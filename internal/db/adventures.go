package db

import (
	"encoding/json"
	"time"
)

// Adventure represents a row in the adventures table.
type Adventure struct {
	ID          int64
	UserID      int64
	Title       string
	Mode        string
	CurrentStep string
	StateJSON   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Message represents a row in the messages table.
type Message struct {
	ID          int64
	AdventureID int64
	Role        string
	Step        string
	Content     string
	CreatedAt   time.Time
}

// Location represents a single map location with connections to other locations.
type Location struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Connects    []string `json:"connects"`
}

// AdventureState is parsed from StateJSON and holds the structured state of an adventure.
type AdventureState struct {
	Scenario              string   `json:"scenario,omitempty"`
	Setting               string   `json:"setting,omitempty"`
	Transgression         string   `json:"transgression,omitempty"`
	TransgressionRoll     int      `json:"transgression_roll,omitempty"`
	Omens                 string   `json:"omens,omitempty"`
	OmensRoll             int      `json:"omens_roll,omitempty"`
	Manifestation         string   `json:"manifestation,omitempty"`
	ManifestationRoll     int      `json:"manifestation_roll,omitempty"`
	Banishment            string   `json:"banishment,omitempty"`
	BanishmentRoll        int      `json:"banishment_roll,omitempty"`
	Slumber               string   `json:"slumber,omitempty"`
	SlumberRoll           int      `json:"slumber_roll,omitempty"`
	HorrorSummary         string   `json:"horror_summary,omitempty"`
	Survive               string   `json:"survive,omitempty"`
	Solve                 string   `json:"solve,omitempty"`
	Save                  string   `json:"save,omitempty"`
	MapLocations          []Location        `json:"map_locations,omitempty"`
	Theme                 string            `json:"theme,omitempty"`
	FinalDoc              string            `json:"final_doc,omitempty"`
	StepSummaries         map[string]string `json:"step_summaries,omitempty"`
}

// ParseState parses the StateJSON field into an AdventureState.
func (a *Adventure) ParseState() (*AdventureState, error) {
	s := &AdventureState{}
	if a.StateJSON == "" || a.StateJSON == "{}" {
		return s, nil
	}
	if err := json.Unmarshal([]byte(a.StateJSON), s); err != nil {
		return s, err
	}
	return s, nil
}

// CreateAdventure inserts a new adventure and returns its ID.
func (db *DB) CreateAdventure(userID int64, title, mode string) (int64, error) {
	var id int64
	err := db.QueryRow(`
		INSERT INTO adventures (user_id, title, mode, current_step, state_json)
		VALUES (?, ?, ?, 'scenario', '{}')
		RETURNING id
	`, userID, title, mode).Scan(&id)
	return id, err
}

// GetAdventure returns an adventure by ID, verifying ownership.
func (db *DB) GetAdventure(id, userID int64) (*Adventure, error) {
	a := &Adventure{}
	err := db.QueryRow(`
		SELECT id, user_id, title, mode, current_step, state_json, created_at, updated_at
		FROM adventures WHERE id = ? AND user_id = ?
	`, id, userID).Scan(
		&a.ID, &a.UserID, &a.Title, &a.Mode, &a.CurrentStep,
		&a.StateJSON, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// ListAdventures returns all adventures for a user, newest first.
func (db *DB) ListAdventures(userID int64) ([]*Adventure, error) {
	rows, err := db.Query(`
		SELECT id, user_id, title, mode, current_step, state_json, created_at, updated_at
		FROM adventures WHERE user_id = ?
		ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var adventures []*Adventure
	for rows.Next() {
		a := &Adventure{}
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.Title, &a.Mode, &a.CurrentStep,
			&a.StateJSON, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		adventures = append(adventures, a)
	}
	return adventures, rows.Err()
}

// UpdateAdventureStep updates the current_step and state_json of an adventure.
func (db *DB) UpdateAdventureStep(id int64, step, stateJSON string) error {
	_, err := db.Exec(`
		UPDATE adventures SET current_step = ?, state_json = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, step, stateJSON, id)
	return err
}

// UpdateAdventureState updates only the state_json of an adventure.
func (db *DB) UpdateAdventureState(id int64, stateJSON string) error {
	_, err := db.Exec(`
		UPDATE adventures SET state_json = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, stateJSON, id)
	return err
}

// AddMessage inserts a new message for an adventure step.
func (db *DB) AddMessage(adventureID int64, role, step, content string) error {
	_, err := db.Exec(`
		INSERT INTO messages (adventure_id, role, step, content)
		VALUES (?, ?, ?, ?)
	`, adventureID, role, step, content)
	return err
}

// GetMessages returns all messages for an adventure step, oldest first.
func (db *DB) GetMessages(adventureID int64, step string) ([]*Message, error) {
	rows, err := db.Query(`
		SELECT id, adventure_id, role, step, content, created_at
		FROM messages WHERE adventure_id = ? AND step = ?
		ORDER BY created_at ASC
	`, adventureID, step)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		m := &Message{}
		if err := rows.Scan(&m.ID, &m.AdventureID, &m.Role, &m.Step, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

// DeleteStepMessages removes all messages for a specific step of an adventure.
func (db *DB) DeleteStepMessages(adventureID int64, step string) error {
	_, err := db.Exec(`DELETE FROM messages WHERE adventure_id = ? AND step = ?`, adventureID, step)
	return err
}

// DeleteAdventure removes an adventure and its messages, verifying ownership.
func (db *DB) DeleteAdventure(id, userID int64) error {
	// Delete messages first (no cascade in SQLite without pragma)
	if _, err := db.Exec(`DELETE FROM messages WHERE adventure_id = ?`, id); err != nil {
		return err
	}
	_, err := db.Exec(`DELETE FROM adventures WHERE id = ? AND user_id = ?`, id, userID)
	return err
}
