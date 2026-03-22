package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lehmann314159/moshag/internal/auth"
	"github.com/lehmann314159/moshag/internal/db"
	"github.com/lehmann314159/moshag/internal/ollama"
	"github.com/lehmann314159/moshag/internal/tables"
)

// adventurePageData is passed to the adventure workspace template.
type adventurePageData struct {
	PageData
	Adventure   *db.Adventure
	State       *db.AdventureState
	Messages    []*db.Message
	StepOrder   []string
	StepOptions []tables.StepOption
	Phase       string // "gather" or "collaborate"
}

// stepAutoStart returns true for steps that skip the gather phase and go straight to AI conversation.
func stepAutoStart(step string) bool {
	return len(tables.StepOptions(step)) == 0
}

// autoStartChat fires the opening AI message for a step in a goroutine.
func (h *Handlers) autoStartChat(adventureID int64, step string, state *db.AdventureState) {
	userMsg := "Let's begin the " + tables.StepLabel(step) + " step."
	if err := h.db.AddMessage(adventureID, "user", step, userMsg); err != nil {
		log.Printf("auto-start save user message: %v", err)
		return
	}
	dbMessages, _ := h.db.GetMessages(adventureID, step)
	chatMessages := buildChatMessages(step, state, dbMessages)
	response, err := h.ollama.Chat(context.Background(), chatMessages)
	if err != nil {
		log.Printf("auto-start chat %s/%d: %v", step, adventureID, err)
		return
	}
	if err := h.db.AddMessage(adventureID, "assistant", step, response); err != nil {
		log.Printf("auto-start save assistant message: %v", err)
	}
}

// chatPartialData is passed to the chat-messages partial.
type chatPartialData struct {
	AdventureID int64
	Step        string
	Messages    []*db.Message
}

// collaboratePanelData is passed to the collaborate-panel partial.
type collaboratePanelData struct {
	Adventure *db.Adventure
	State     *db.AdventureState
	Messages  []*db.Message
	Streaming bool
}

// stepFormData is passed to the step-form partial.
type stepFormData struct {
	Adventure   *db.Adventure
	State       *db.AdventureState
	StepOptions []tables.StepOption
}

// currentUserID returns the logged-in user's DB ID, or the guest user ID if not logged in.
func currentUserID(r *http.Request) int64 {
	if sess := auth.GetSession(r.Context()); sess != nil {
		return sess.DBUserID
	}
	return db.GuestUserID
}

// parseAdventureID parses the {id} URL param and returns the int64, or -1 on error.
func parseAdventureID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}

// NewAdventure handles POST /adventures/new — creates a new adventure.
func (h *Handlers) NewAdventure(w http.ResponseWriter, r *http.Request) {
	userID := currentUserID(r)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		title = "Untitled Adventure"
	}
	mode := r.FormValue("mode")
	if mode == "" {
		mode = "manual"
	}

	id, err := h.db.CreateAdventure(userID, title, mode)
	if err != nil {
		log.Printf("create adventure: %v", err)
		http.Error(w, "Failed to create adventure", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/adventures/"+strconv.FormatInt(id, 10), http.StatusFound)
}

// ShowAdventure handles GET /adventures/{id}.
func (h *Handlers) ShowAdventure(w http.ResponseWriter, r *http.Request) {
	id, err := parseAdventureID(r)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	adventure, err := h.db.GetAdventure(id, currentUserID(r))
	if err != nil {
		log.Printf("get adventure %d: %v", id, err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	state, err := adventure.ParseState()
	if err != nil {
		log.Printf("parse state for adventure %d: %v", id, err)
		state = &db.AdventureState{}
	}

	messages, err := h.db.GetMessages(id, adventure.CurrentStep)
	if err != nil {
		log.Printf("get messages for adventure %d: %v", id, err)
		messages = nil
	}

	phase := "gather"
	if len(messages) > 0 {
		phase = "collaborate"
	} else if stepAutoStart(adventure.CurrentStep) {
		phase = "collaborate"
		go h.autoStartChat(id, adventure.CurrentStep, state)
	}

	data := adventurePageData{
		PageData:    pageData(r, "MOSHAG — "+adventure.Title, "adventure"),
		Adventure:   adventure,
		State:       state,
		Messages:    messages,
		StepOrder:   tables.StepOrder,
		StepOptions: tables.StepOptions(adventure.CurrentStep),
		Phase:       phase,
	}
	h.render(w, "base", data)
}

// StartChat handles POST /adventures/{id}/start — first message from the gather form.
// Sends the initial message to the AI and returns the collaborate panel.
func (h *Handlers) StartChat(w http.ResponseWriter, r *http.Request) {
	id, err := parseAdventureID(r)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	adventure, err := h.db.GetAdventure(id, currentUserID(r))
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	selection := strings.Join(r.Form["selection"], ", ")
	guidance := r.FormValue("guidance")
	context := r.FormValue("message")

	userMsg := buildUserMessage(adventure.CurrentStep, selection, guidance, context)
	if userMsg == "" {
		userMsg = "Let's begin the " + tables.StepLabel(adventure.CurrentStep) + " step."
	}

	// Save selection to state.
	state, err := adventure.ParseState()
	if err != nil {
		state = &db.AdventureState{}
	}
	if selection != "" {
		saveSelectionToState(adventure.CurrentStep, selection, state)
		if b, err := json.Marshal(state); err == nil {
			_ = h.db.UpdateAdventureState(id, string(b))
		}
	}

	// Save user message.
	if err := h.db.AddMessage(id, "user", adventure.CurrentStep, userMsg); err != nil {
		log.Printf("add user message: %v", err)
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	// Load messages saved so far (user message only — AI will stream).
	dbMessages, _ := h.db.GetMessages(id, adventure.CurrentStep)

	h.renderPartial(w, "collaborate-panel", collaboratePanelData{
		Adventure: adventure,
		State:     state,
		Messages:  dbMessages,
		Streaming: true,
	})
}

// Chat handles POST /adventures/{id}/chat — continue the conversation.
func (h *Handlers) Chat(w http.ResponseWriter, r *http.Request) {
	id, err := parseAdventureID(r)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	adventure, err := h.db.GetAdventure(id, currentUserID(r))
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	userMsg := r.FormValue("message")
	if userMsg == "" {
		// Return current messages unchanged rather than an error.
		msgs, _ := h.db.GetMessages(id, adventure.CurrentStep)
		h.renderPartial(w, "chat-messages", chatPartialData{
			AdventureID: id,
			Step:        adventure.CurrentStep,
			Messages:    msgs,
		})
		return
	}

	// Save user message.
	if err := h.db.AddMessage(id, "user", adventure.CurrentStep, userMsg); err != nil {
		log.Printf("add user message: %v", err)
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	// Load messages saved so far (user message only — AI will stream).
	dbMessages, _ := h.db.GetMessages(id, adventure.CurrentStep)

	h.renderPartial(w, "chat-messages-stream", chatPartialData{
		AdventureID: id,
		Step:        adventure.CurrentStep,
		Messages:    dbMessages,
	})
}

// Done handles POST /adventures/{id}/done — generate a step summary, save it, and advance.
func (h *Handlers) Done(w http.ResponseWriter, r *http.Request) {
	id, err := parseAdventureID(r)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	adventure, err := h.db.GetAdventure(id, currentUserID(r))
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if adventure.CurrentStep == "complete" {
		http.Error(w, "Adventure already complete", http.StatusBadRequest)
		return
	}

	state, err := adventure.ParseState()
	if err != nil {
		state = &db.AdventureState{}
	}

	// Generate a step summary from the conversation.
	dbMessages, _ := h.db.GetMessages(id, adventure.CurrentStep)
	if len(dbMessages) > 0 {
		chatMessages := buildChatMessages(adventure.CurrentStep, state, dbMessages)
		summaryPrompt := "In 2-3 sentences, summarize the key decisions made for the " +
			tables.StepLabel(adventure.CurrentStep) +
			" step. Be specific about the details chosen — this will inform the rest of the adventure."
		chatMessages = append(chatMessages, ollama.Message{Role: "user", Content: summaryPrompt})

		summary, err := h.ollama.Chat(r.Context(), chatMessages)
		if err != nil {
			log.Printf("generate step summary for %s/%d: %v", adventure.CurrentStep, id, err)
		} else if summary != "" {
			if state.StepSummaries == nil {
				state.StepSummaries = make(map[string]string)
			}
			state.StepSummaries[adventure.CurrentStep] = summary
		}
	}

	// Advance to next step and save updated state.
	nextStep := tables.NextStep(adventure.CurrentStep)
	stateBytes, _ := json.Marshal(state)
	if err := h.db.UpdateAdventureStep(id, nextStep, string(stateBytes)); err != nil {
		log.Printf("update step: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Reload adventure for new step.
	adventure, err = h.db.GetAdventure(id, currentUserID(r))
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	state, _ = adventure.ParseState()
	messages, _ := h.db.GetMessages(id, nextStep)

	phase := "gather"
	if len(messages) > 0 {
		phase = "collaborate"
	} else if stepAutoStart(nextStep) {
		phase = "collaborate"
		go h.autoStartChat(id, nextStep, state)
	}

	data := adventurePageData{
		PageData:    pageData(r, "MOSHAG — "+adventure.Title, "adventure"),
		Adventure:   adventure,
		State:       state,
		Messages:    messages,
		StepOrder:   tables.StepOrder,
		StepOptions: tables.StepOptions(nextStep),
		Phase:       phase,
	}
	h.renderPartial(w, "adventure-workspace", data)
}

// Roll handles POST /adventures/{id}/roll — roll a table for the current step.
func (h *Handlers) Roll(w http.ResponseWriter, r *http.Request) {
	id, err := parseAdventureID(r)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	adventure, err := h.db.GetAdventure(id, currentUserID(r))
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	tableName := r.FormValue("table")

	state, err := adventure.ParseState()
	if err != nil {
		state = &db.AdventureState{}
	}

	switch tableName {
	case "setting":
		roll := tables.RollD10()
		result := tables.Lookup(tables.Settings, roll)
		state.Setting = result
	case "transgression":
		roll := tables.RollD100()
		result := tables.Lookup(tables.Transgressions, roll)
		state.Transgression = result
		state.TransgressionRoll = roll
	case "omens":
		roll := tables.RollD100()
		result := tables.Lookup(tables.Omens, roll)
		state.Omens = result
		state.OmensRoll = roll
	case "manifestation":
		roll := tables.RollD100()
		result := tables.Lookup(tables.Manifestations, roll)
		state.Manifestation = result
		state.ManifestationRoll = roll
	case "banishment":
		roll := tables.RollD100()
		result := tables.Lookup(tables.Banishments, roll)
		state.Banishment = result
		state.BanishmentRoll = roll
	case "slumber":
		roll := tables.RollD100()
		result := tables.Lookup(tables.Slumbers, roll)
		state.Slumber = result
		state.SlumberRoll = roll
	case "all-tombs":
		state.TransgressionRoll = tables.RollD100()
		state.Transgression = tables.Lookup(tables.Transgressions, state.TransgressionRoll)
		state.OmensRoll = tables.RollD100()
		state.Omens = tables.Lookup(tables.Omens, state.OmensRoll)
		state.ManifestationRoll = tables.RollD100()
		state.Manifestation = tables.Lookup(tables.Manifestations, state.ManifestationRoll)
		state.BanishmentRoll = tables.RollD100()
		state.Banishment = tables.Lookup(tables.Banishments, state.BanishmentRoll)
		state.SlumberRoll = tables.RollD100()
		state.Slumber = tables.Lookup(tables.Slumbers, state.SlumberRoll)
	}

	stateBytes, err := json.Marshal(state)
	if err != nil {
		log.Printf("marshal state: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := h.db.UpdateAdventureState(id, string(stateBytes)); err != nil {
		log.Printf("update state: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If all TOMBS rolls are done, trigger LLM synthesis.
	allTombsDone := state.Transgression != "" && state.Omens != "" &&
		state.Manifestation != "" && state.Banishment != "" && state.Slumber != ""
	if allTombsDone && state.HorrorSummary == "" && adventure.CurrentStep == "tombs" {
		go h.synthesizeHorror(id, adventure, state)
	}

	adventure, err = h.db.GetAdventure(id, currentUserID(r))
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	state, _ = adventure.ParseState()
	h.renderPartial(w, "step-form", stepFormData{
		Adventure:   adventure,
		State:       state,
		StepOptions: tables.StepOptions(adventure.CurrentStep),
	})
}

// NextStep handles POST /adventures/{id}/next — advance without generating a summary.
func (h *Handlers) NextStep(w http.ResponseWriter, r *http.Request) {
	id, err := parseAdventureID(r)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	adventure, err := h.db.GetAdventure(id, currentUserID(r))
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if adventure.CurrentStep == "complete" {
		http.Error(w, "Adventure already complete", http.StatusBadRequest)
		return
	}

	nextStep := tables.NextStep(adventure.CurrentStep)

	state, _ := adventure.ParseState()
	stateBytes, _ := json.Marshal(state)

	if err := h.db.UpdateAdventureStep(id, nextStep, string(stateBytes)); err != nil {
		log.Printf("update step: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	adventure, err = h.db.GetAdventure(id, currentUserID(r))
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	state, _ = adventure.ParseState()
	messages, _ := h.db.GetMessages(id, nextStep)

	phase := "gather"
	if len(messages) > 0 {
		phase = "collaborate"
	} else if stepAutoStart(nextStep) {
		phase = "collaborate"
		go h.autoStartChat(id, nextStep, state)
	}

	data := adventurePageData{
		PageData:    pageData(r, "MOSHAG — "+adventure.Title, "adventure"),
		Adventure:   adventure,
		State:       state,
		Messages:    messages,
		StepOrder:   tables.StepOrder,
		StepOptions: tables.StepOptions(nextStep),
		Phase:       phase,
	}
	h.renderPartial(w, "adventure-workspace", data)
}

// DeleteAdventure handles DELETE /adventures/{id}.
func (h *Handlers) DeleteAdventure(w http.ResponseWriter, r *http.Request) {
	userID := currentUserID(r)

	id, err := parseAdventureID(r)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err := h.db.DeleteAdventure(id, userID); err != nil {
		log.Printf("delete adventure %d: %v", id, err)
		http.Error(w, "Failed to delete", http.StatusInternalServerError)
		return
	}

	adventures, err := h.db.ListAdventures(userID)
	if err != nil {
		adventures = nil
	}
	h.renderPartial(w, "adventure-list", struct {
		Adventures []*db.Adventure
	}{Adventures: adventures})
}

// ClearStep handles POST /adventures/{id}/clear — deletes messages for the current step and returns to gather phase.
func (h *Handlers) ClearStep(w http.ResponseWriter, r *http.Request) {
	id, err := parseAdventureID(r)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	adventure, err := h.db.GetAdventure(id, currentUserID(r))
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err := h.db.DeleteStepMessages(id, adventure.CurrentStep); err != nil {
		log.Printf("clear step messages: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	state, _ := adventure.ParseState()
	data := adventurePageData{
		PageData:    pageData(r, "MOSHAG — "+adventure.Title, "adventure"),
		Adventure:   adventure,
		State:       state,
		Messages:    nil,
		StepOrder:   tables.StepOrder,
		StepOptions: tables.StepOptions(adventure.CurrentStep),
		Phase:       "gather",
	}
	h.renderPartial(w, "adventure-workspace", data)
}

// Messages handles GET /adventures/{id}/messages — polls for chat messages.
func (h *Handlers) Messages(w http.ResponseWriter, r *http.Request) {
	id, err := parseAdventureID(r)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	adventure, err := h.db.GetAdventure(id, currentUserID(r))
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	step := r.URL.Query().Get("step")
	if step == "" {
		step = adventure.CurrentStep
	}

	msgs, _ := h.db.GetMessages(id, step)
	h.renderPartial(w, "chat-messages", chatPartialData{
		AdventureID: id,
		Step:        step,
		Messages:    msgs,
	})
}

// buildChatMessages converts DB messages to Ollama chat format, prepending system prompt.
// It injects summaries from previously completed steps into the system prompt.
func buildChatMessages(step string, state *db.AdventureState, dbMessages []*db.Message) []ollama.Message {
	systemPrompt := tables.SystemPrompt(step)

	if state != nil && len(state.StepSummaries) > 0 {
		contextLines := ""
		for _, prevStep := range tables.StepOrder {
			if prevStep == step {
				break
			}
			if summary, ok := state.StepSummaries[prevStep]; ok && summary != "" {
				contextLines += "- " + tables.StepLabel(prevStep) + ": " + summary + "\n"
			}
		}
		if contextLines != "" {
			systemPrompt += "\n\nContext from previous steps:\n" + contextLines
		}
	}

	messages := []ollama.Message{
		{Role: "system", Content: systemPrompt},
	}
	for _, m := range dbMessages {
		messages = append(messages, ollama.Message{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	return messages
}

// synthesizeHorror sends an LLM synthesis after all TOMBS rolls are done.
func (h *Handlers) synthesizeHorror(adventureID int64, adventure *db.Adventure, state *db.AdventureState) {
	prompt := "The Warden has rolled the TOMBS cycle. Here are the results:\n\n" +
		"Transgression (what was done): " + state.Transgression + "\n" +
		"Omens (warning signs): " + state.Omens + "\n" +
		"Manifestation (what the Horror is): " + state.Manifestation + "\n" +
		"Banishment (how to stop it): " + state.Banishment + "\n" +
		"Slumber (what happens after): " + state.Slumber + "\n\n" +
		"Synthesize these into a coherent, terrifying Horror for the adventure. " +
		"Give it a name, describe its nature, and explain how these elements connect."

	messages := []ollama.Message{
		{Role: "system", Content: tables.SystemPrompt("tombs")},
		{Role: "user", Content: prompt},
	}

	response, err := h.ollama.Chat(context.Background(), messages)
	if err != nil {
		log.Printf("synthesize horror for %d: %v", adventureID, err)
		return
	}

	if err := h.db.AddMessage(adventureID, "assistant", "tombs", response); err != nil {
		log.Printf("save horror synthesis: %v", err)
		return
	}

	if adventure != nil {
		freshState, _ := adventure.ParseState()
		freshState.HorrorSummary = response
		stateBytes, _ := json.Marshal(freshState)
		if err := h.db.UpdateAdventureState(adventureID, string(stateBytes)); err != nil {
			log.Printf("update horror summary: %v", err)
		}
	}
}

// Stream handles GET /adventures/{id}/stream — SSE streaming of an AI chat response.
// It builds the chat message history, streams tokens from Ollama, saves the full response, then sends [DONE].
func (h *Handlers) Stream(w http.ResponseWriter, r *http.Request) {
	id, err := parseAdventureID(r)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	adventure, err := h.db.GetAdventure(id, currentUserID(r))
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	step := r.URL.Query().Get("step")
	if step == "" {
		step = adventure.CurrentStep
	}

	state, _ := adventure.ParseState()
	dbMessages, _ := h.db.GetMessages(id, step)
	chatMessages := buildChatMessages(step, state, dbMessages)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	sendEvent := func(data string) {
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}

	var assembled strings.Builder
	err = h.ollama.ChatStream(r.Context(), chatMessages, func(token string) {
		assembled.WriteString(token)
		encoded, _ := json.Marshal(token)
		sendEvent(string(encoded))
	})

	if err != nil {
		log.Printf("stream %s/%d: %v", step, id, err)
		sendEvent(`"[ERROR]"`)
		return
	}

	fullResponse := assembled.String()
	if fullResponse != "" {
		if err := h.db.AddMessage(id, "assistant", step, fullResponse); err != nil {
			log.Printf("save streamed response %s/%d: %v", step, id, err)
		}
	}

	sendEvent(`"[DONE]"`)
}

// saveSelectionToState writes the primary select choice into the structured state.
func saveSelectionToState(step, selection string, state *db.AdventureState) {
	switch step {
	case "scenario":
		state.Scenario = selection
	case "setting":
		state.Setting = selection
	}
}

// buildUserMessage combines the selection, guidance style, and freeform context
// into a single prompt string to send to the LLM.
func buildUserMessage(step, selection, guidance, context string) string {
	var parts []string

	if selection != "" {
		parts = append(parts, tables.StepOptionPrompt(step, selection))
	}
	if context != "" {
		parts = append(parts, context)
	}

	msg := ""
	for i, p := range parts {
		if i == 0 {
			msg = p
		} else {
			msg += " — " + p
		}
	}

	suffix := map[string]string{
		"clarify": "Ask me one clarifying question.",
		"choices": "Give me three broad choices, each in one or two sentences.",
		"idea":    "Give me one fully-fleshed out idea, written as if it's already decided.",
	}[guidance]

	if suffix != "" {
		if msg == "" {
			msg = suffix
		} else {
			msg += " — " + suffix
		}
	}

	return msg
}

// buildStateContext creates a context summary from adventure state for LLM prompts.
func buildStateContext(state *db.AdventureState) string {
	if state == nil {
		return ""
	}
	var parts []string
	if state.Scenario != "" {
		parts = append(parts, "Scenario: "+state.Scenario)
	}
	if state.Setting != "" {
		parts = append(parts, "Setting: "+state.Setting)
	}
	if state.Transgression != "" {
		parts = append(parts, "TOMBS — Transgression: "+state.Transgression)
	}
	if state.Manifestation != "" {
		parts = append(parts, "TOMBS — Manifestation: "+state.Manifestation)
	}
	if state.Banishment != "" {
		parts = append(parts, "TOMBS — Banishment: "+state.Banishment)
	}
	if state.HorrorSummary != "" {
		parts = append(parts, "The Horror (synthesized): "+state.HorrorSummary)
	}
	if state.Survive != "" {
		parts = append(parts, "Survive: "+state.Survive)
	}
	if state.Solve != "" {
		parts = append(parts, "Solve: "+state.Solve)
	}
	if state.Save != "" {
		parts = append(parts, "Save (NPCs): "+state.Save)
	}
	if len(parts) == 0 {
		return ""
	}
	result := "Adventure context established so far:\n"
	for _, p := range parts {
		result += "- " + p + "\n"
	}
	return result
}
