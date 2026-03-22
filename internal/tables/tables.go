package tables

import "math/rand"

// TableEntry represents a range-based D100 table entry.
type TableEntry struct {
	Min  int
	Max  int
	Text string
}

// Lookup finds the entry for a given roll (0-99).
func Lookup(table []TableEntry, roll int) string {
	for _, e := range table {
		if roll >= e.Min && roll <= e.Max {
			return e.Text
		}
	}
	return ""
}

// RollD100 returns a random 0-99.
func RollD100() int { return rand.Intn(100) }

// RollD10 returns a random 0-9.
func RollD10() int { return rand.Intn(10) }

// Scenarios is a list of 10 starting scenarios (index 0-9, displayed as 1-10).
var Scenarios = []string{
	"Explore the Unknown",             // 1
	"Investigate a Strange Rumor",     // 2
	"Salvage a Derelict Ship",         // 3
	"Exterminate an Otherworldly Threat", // 4
	"Visit an Offworld Colony",        // 5
	"Undertake a Dangerous Mission",   // 6
	"Survive a Colossal Disaster",     // 7
	"Respond to a Distress Signal",    // 8
	"Transport Precious Cargo",        // 9
	"Make Contact with the Beyond",    // 10
}

// ScenarioDescriptions provides the flavour text for each scenario (parallel to Scenarios).
var ScenarioDescriptions = []string{
	"The crew is hired to survey an uncharted planet, or to explore the interior of a strange vessel which has recently appeared at the edge of Rim space. They are on their own in unfamiliar territory with no one to call on for help.",
	"Something is alive in the vents. Colonists are disappearing. Someone is leaving messages on the comms terminal, from two years in the future. Separating fact from fiction was the easy part. The hard part was learning to live with the truth.",
	"Distress signal on repeat. Scans show no signs of life. Finding a derelict ship can be a chance to strike it rich off scrap and loot, but not every risk has a reward, and every abandoned ship is abandoned for a reason.",
	"No one goes outside anymore. They've been scratching at the walls for days and the Company has made their decision. Wipe them out, by any means necessary. Just bring back a sample for testing when you're done — a living one.",
	"The Company hasn't heard from the miners on PK-294, and the shareholders are getting restless. Take the next Jumpliner out to the Rim and get production back on track. Offworld colonies for many are a new start, but for others a new end.",
	"A C-level's child has been 'kidnapped' by a fringe religious group. Android liberation activists want to sabotage a corporate synthetic production facility. There's always work for people with no scruples in need of quick credits.",
	"Abandon ship! Radiation leaks and warp anomalies. Make it to the escape pods before the whole station collapses. Unstable environments and dangerous weather. Escaping disaster is never a bad place to start.",
	"Help is never nearby out on the Rim, and responding to a distress signal is a spacer's duty. There's only one problem though: you can never tell the legitimate cries for help from the traps laid by wolves in sheep's clothing.",
	"You needed a job, and your contact came through. They won't tell you what's in the container, only that it needs to be at Outpost-683 in six weeks. Don't open it, don't scan it, and whatever you do, don't listen to anything it says.",
	"They found it at the edge of the system, on a frozen and forgotten moon. Eons old, intricate stonework, when our probe arrived it started humming a tune. They've been looking for us for a long, long time. Now they've found us.",
}

// Settings is a D10 table of adventure settings.
var Settings = []TableEntry{
	{0, 0, "Space Station"},
	{1, 1, "Aboard Your Own Ship"},
	{2, 2, "Military Outpost"},
	{3, 3, "Prison Complex"},
	{4, 4, "Derelict Spacecraft"},
	{5, 5, "Religious Compound"},
	{6, 6, "Mining Colony"},
	{7, 7, "Research Facility"},
	{8, 8, "Underwater Base"},
	{9, 9, "Mothership"},
}

// Transgressions is a D100 table — what initiated the Horror.
var Transgressions = []TableEntry{
	{0, 4, "Making first contact"},
	{5, 9, "Studying arcane text"},
	{10, 14, "Boarding a ship"},
	{15, 19, "Opening a grave"},
	{20, 24, "Mining strange ore"},
	{25, 29, "Trespassing"},
	{30, 34, "Gross negligence"},
	{35, 39, "Tampering with biology"},
	{40, 44, "Reneging on a deal"},
	{45, 49, "Disturbing holy site"},
	{50, 54, "Leaving someone behind"},
	{55, 59, "Study of strange relic"},
	{60, 64, "Forgotten atrocity"},
	{65, 69, "Interfacing with forbidden technology"},
	{70, 74, "Landing on uncharted planet"},
	{75, 79, "Altering its natural habitat"},
	{80, 84, "Breaking a cultural taboo"},
	{85, 89, "Failing to stop a previous Transgression"},
	{90, 94, "Ingesting an unknown substance"},
	{95, 99, "Allowing harm to come to an innocent"},
}

// Omens is a D100 table — the warning signs before the Horror manifests.
var Omens = []TableEntry{
	{0, 4, "Dead animals"},
	{5, 9, "Visions of future victims"},
	{10, 14, "Writing on the wall"},
	{15, 19, "Stigmata"},
	{20, 24, "Unexplained suicides"},
	{25, 29, "Distress signal"},
	{30, 34, "Stranger appears"},
	{35, 39, "Abnormal birth"},
	{40, 44, "Unlucky numbers"},
	{45, 49, "Ancient distress beacon"},
	{50, 54, "Android having visions"},
	{55, 59, "Ancient recorded warning"},
	{60, 64, "Researcher's incoherent notes and findings"},
	{65, 69, "Irrational computer behavior"},
	{70, 74, "Significant astrological alignment"},
	{75, 79, "Speaking in tongues"},
	{80, 84, "Mysterious disappearances"},
	{85, 89, "Strange weather phenomena"},
	{90, 94, "Ancient calendar foretells of its arrival"},
	{95, 99, "Gruesomely displayed corpse(s)"},
}

// Manifestations is a D100 table — the form the Horror takes.
var Manifestations = []TableEntry{
	{0, 4, "Alien creature"},
	{5, 9, "Deranged killer"},
	{10, 14, "Elder evil returned"},
	{15, 19, "Cult"},
	{20, 24, "Tainted technology"},
	{25, 29, "Colossal space being"},
	{30, 34, "Ruthless apex predator"},
	{35, 39, "Ghost or spirit"},
	{40, 44, "Tangled mass of flesh"},
	{45, 49, "Mutated being"},
	{50, 54, "Child"},
	{55, 59, "Biological experiment"},
	{60, 64, "Sentient environment"},
	{65, 69, "Gateway or portal"},
	{70, 74, "Dream"},
	{75, 79, "Cybernetic organism"},
	{80, 84, "Haunted location"},
	{85, 89, "Doppelganger"},
	{90, 94, "Invisible being"},
	{95, 99, "Mothership"},
}

// Banishments is a D100 table — how to stop the Horror.
var Banishments = []TableEntry{
	{0, 4, "Righting a wrong"},
	{5, 9, "Human sacrifice"},
	{10, 14, "Vaccine"},
	{15, 19, "Only harmed by fire"},
	{20, 24, "Nuking from orbit"},
	{25, 29, "Obscure occult ritual"},
	{30, 34, "Returning it to its home"},
	{35, 39, "Tough, but killable"},
	{40, 44, "Giving it what it wants"},
	{45, 49, "Special weapon"},
	{50, 54, "Making a pact with it"},
	{55, 59, "Serving it"},
	{60, 64, "Learning its true identity"},
	{65, 69, "Certain kinds of light"},
	{70, 74, "It can't be killed, only avoided"},
	{75, 79, "Inter remains in their rightful resting place"},
	{80, 84, "Closing portal/gate to another realm"},
	{85, 89, "Requires a certain time/location"},
	{90, 94, "Sending it to another dimension"},
	{95, 99, "Trapping it inside a powerful container"},
}

// Slumbers is a D100 table — what happens after the Horror is defeated/banished.
var Slumbers = []TableEntry{
	{0, 4, "Returns in 100 years"},
	{5, 9, "Recurring hallucinations"},
	{10, 14, "Victims forever scarred"},
	{15, 19, "Slumbers until next Jump"},
	{20, 24, "Retreats into hiding"},
	{25, 29, "Feigns death, escapes"},
	{30, 34, "Awaits next Transgression"},
	{35, 39, "Recurring nightmares"},
	{40, 44, "Possesses closest victim"},
	{45, 49, "Awakens if Transgression is repeated"},
	{50, 54, "Hibernates deep underground"},
	{55, 59, "Whispers from the shadows"},
	{60, 64, "Evolves into its more powerful form"},
	{65, 69, "Hidden in the background of screens and images"},
	{70, 74, "Slumbers in its killer's mind"},
	{75, 79, "Herald of a greater Horror to come"},
	{80, 84, "Uploads into nearest computer"},
	{85, 89, "Never stay in one place for too long or it finds you"},
	{90, 94, "Parental entity comes looking for answers"},
	{95, 99, "Apocalyptic events set in motion"},
}

// Themes is a D100 table — thematic keywords for the adventure.
var Themes = []TableEntry{
	{0, 3, "Death, ancient, arise"},
	{4, 9, "Underwater, sunken, drowning"},
	{10, 12, "Politics, government, nationalism"},
	{13, 16, "Humanity, love, memory"},
	{17, 19, "Resistance, struggle, suffering"},
	{20, 22, "Travel, road-weariness, rural"},
	{23, 25, "Darkness, absence, void"},
	{26, 29, "Medicine, hospitals, surgery"},
	{30, 32, "Rust, the Machine, noise"},
	{33, 35, "Transformation, rebirth, loss"},
	{36, 38, "Childhood, innocence, time"},
	{39, 41, "Underground, crime, buried"},
	{42, 43, "Fading beauty, age, fame"},
	{44, 46, "Technology, excess, decay"},
	{47, 49, "Abduction, identity, silence"},
	{50, 52, "The City, rain, flood"},
	{53, 55, "Fear, the afterlife, prophecy"},
	{56, 58, "Factories, work, oppression"},
	{59, 61, "Belief, god, hell"},
	{62, 64, "Cold, deep, snow"},
	{65, 67, "Fire, ashes, war"},
	{68, 71, "Hunger, famine, food"},
	{72, 74, "Pleasure, touch, passion"},
	{75, 77, "Artifice, dolls, toys"},
	{78, 81, "Meat, slaughter, animal"},
	{82, 84, "Truth, solitude, loneliness"},
	{85, 87, "Wilderness, nature, growth"},
	{88, 91, "Capitalism, greed, fortune"},
	{92, 94, "Chaos, change, laughter"},
	{95, 99, "Abandoned, empty, forgotten"},
}

// PuzzleComponents is a D100 table of puzzle types.
var PuzzleComponents = []TableEntry{
	{0, 4, "Alarm — If solved incorrectly sets off an alarm, alerting nearby enemies"},
	{5, 9, "Connect the Dots — Connect a series of 'dots' using some kind of item"},
	{10, 13, "Construction — Construct an item using pieces that are given or must be found"},
	{14, 17, "Dilemma — Choose between the lesser of two evils or the greater of two goods"},
	{18, 21, "Egg Find — Safely deliver a delicate or vulnerable item to a certain location"},
	{22, 25, "Find the Clue — Search for one or more objects which lead to Questions, Puzzles, or Answers"},
	{26, 29, "Guardian — Defeat, appease, or bypass a guardian that denies access"},
	{30, 34, "Hazardous Path — Route blocked by danger which must be avoided or circumnavigated"},
	{35, 38, "Illusion — Contains an element that appears to be one thing but is actually another"},
	{39, 42, "Labyrinth — Forces players to navigate a maze or complicated path"},
	{43, 47, "Lock & Key — Collect an item and use it to gain access"},
	{48, 51, "Missing Part — Locate a missing item in order for the puzzle to work properly"},
	{52, 56, "Mundane Obstacle — A real-world problem (broken elevator, collapsed pylon, etc.)"},
	{57, 60, "Outside-the-Box — No obvious solution; players must bring outside resources"},
	{61, 64, "Pattern Recognition — Players notice repeated symbols or repeated information"},
	{65, 69, "Remote Switch — Activate a switch in another location to bypass the obstacle"},
	{70, 73, "Riddle — Requires recitation of a coded phrase or password"},
	{74, 78, "Rising Tide — Danger which escalates naturally, forcing players to solve it quickly"},
	{79, 82, "Sacrifice — Requires players to sacrifice something of great value"},
	{83, 86, "Sequence — Complete a certain number of steps in a specific order"},
	{87, 90, "Teamwork — Multiple players must do something at the same time"},
	{91, 93, "Timelock — Solve the puzzle in a certain amount of time or fail"},
	{94, 96, "Trap — The puzzle punishes players for failed attempts"},
	{97, 99, "Trial and Error — Players must experiment with several ingredients"},
}

// NPCType defines the role matrix for NPCs.
type NPCType struct {
	Name        string
	Power       string // "powerful", "neither", "powerless"
	Helpfulness string // "helpful", "neither", "unhelpful"
	Description string
}

// NPCTypes lists the 10 NPC archetypes from the Warden's Operations Manual.
var NPCTypes = []NPCType{
	{"Gatekeeper", "powerful", "unhelpful", "They have the power to keep you from what you want. As unhelpful to you as they can possibly be."},
	{"Employer", "powerful", "neither", "Wields some of the greatest power over you of anyone. Neither helpful nor unhelpful. Your job is to help them."},
	{"Benefactor", "powerful", "helpful", "Patrons, sponsors, and contacts who can help you and have the resources to do so. The rarest kind of ally."},
	{"Traitor", "neither", "unhelpful", "Double agents, spies, backstabbers. Could be anyone. The more helpful they are, the worse the betrayal is."},
	{"Survivor", "neither", "neither", "Only cares about what's best for them and will do whatever it takes. Help them or get out of their way."},
	{"Expert", "neither", "helpful", "The cream of the crop. When you need them, they are the best kind of help. Their power is their expertise."},
	{"Coward", "powerless", "unhelpful", "Panicked civilians, frightened bystanders, or worthless sycophants. Completely useless."},
	{"Victim", "powerless", "neither", "Powerless to help you, though not forever. They need you the most, but can you spare the effort?"},
	{"Drinking Buddy", "powerless", "helpful", "Friends in low places. They'd do anything for you. There's just not that much they actually can do."},
	{"Wildcard", "any", "any", "Extremely unpredictable agents of chaos, particularly when their backs are against the wall. Use sparingly."},
}

// StepOrder defines the canonical order of adventure generation steps.
var StepOrder = []string{
	"scenario",
	"setting",
	"tombs",
	"survive",
	"solve",
	"save",
	"map",
	"final",
	"complete",
}

// NextStep returns the step that follows the given step.
func NextStep(current string) string {
	for i, s := range StepOrder {
		if s == current && i+1 < len(StepOrder) {
			return StepOrder[i+1]
		}
	}
	return "complete"
}

// StepOption is a selectable option for the primary select on a step's form.
type StepOption struct {
	Value string
	Label string
}

// StepOptions returns the primary select options for a given step.
func StepOptions(step string) []StepOption {
	switch step {
	case "scenario":
		opts := make([]StepOption, len(Scenarios))
		for i, s := range Scenarios {
			opts[i] = StepOption{Value: s, Label: s}
		}
		return opts
	case "setting":
		opts := make([]StepOption, len(Settings))
		for i, s := range Settings {
			opts[i] = StepOption{Value: s.Text, Label: s.Text}
		}
		return opts
	case "tombs":
		return []StepOption{
			{Value: "roll-all", Label: "Roll all five (random)"},
			{Value: "transgression", Label: "Roll Transgression only"},
			{Value: "omens", Label: "Roll Omens only"},
			{Value: "manifestation", Label: "Roll Manifestation only"},
			{Value: "banishment", Label: "Roll Banishment only"},
			{Value: "slumber", Label: "Roll Slumber only"},
		}
	case "survive":
		return []StepOption{
			{Value: "psychological", Label: "Psychological trauma"},
			{Value: "violence", Label: "Violent encounters"},
			{Value: "environmental", Label: "Environmental hazards"},
			{Value: "resources", Label: "Resource scarcity"},
			{Value: "social", Label: "Social pressure"},
		}
	case "solve":
		return nil
	case "save":
		return []StepOption{
			{Value: "Gatekeeper", Label: "Gatekeeper — powerful, unhelpful"},
			{Value: "Employer", Label: "Employer — powerful, neither"},
			{Value: "Benefactor", Label: "Benefactor — powerful, helpful"},
			{Value: "Traitor", Label: "Traitor — neither, unhelpful"},
			{Value: "Survivor", Label: "Survivor — neither, neither"},
			{Value: "Expert", Label: "Expert — neither, helpful"},
			{Value: "Coward", Label: "Coward — powerless, unhelpful"},
			{Value: "Victim", Label: "Victim — powerless, neither"},
			{Value: "Drinking Buddy", Label: "Drinking Buddy — powerless, helpful"},
			{Value: "Wildcard", Label: "Wildcard — unpredictable"},
		}
	case "map":
		return []StepOption{
			{Value: "full-map", Label: "Generate all 10 locations"},
			{Value: "entry", Label: "Focus: entry point"},
			{Value: "horror-lair", Label: "Focus: the Horror's lair"},
			{Value: "safe-zone", Label: "Focus: a place of safety"},
			{Value: "locked-area", Label: "Focus: locked or secret area"},
		}
	case "final":
		return []StepOption{
			{Value: "full-doc", Label: "Generate complete adventure document"},
			{Value: "summary", Label: "Generate a one-page summary"},
		}
	}
	return nil
}

// StepLabel returns the human-readable label for a step name.
func StepLabel(step string) string {
	labels := map[string]string{
		"scenario": "Scenario",
		"setting":  "Setting",
		"tombs":    "Horror (TOMBS)",
		"survive":  "Survive",
		"solve":    "Solve",
		"save":     "Save",
		"map":      "Map",
		"final":    "Final Doc",
		"complete": "Complete",
	}
	if l, ok := labels[step]; ok {
		return l
	}
	return step
}

// StepOptionPrompt converts a step selection value to a prompt string.
func StepOptionPrompt(step, value string) string {
	switch step {
	case "scenario":
		return "I want to use this scenario: " + value
	case "setting":
		return "The setting is: " + value
	case "tombs":
		switch value {
		case "roll-all":
			return "Roll all five TOMBS elements for me."
		default:
			return "Roll the " + value + " for me."
		}
	case "survive":
		return "The main survival pressure for this adventure is: " + value + "."
	case "solve":
		if value == "full-mystery" {
			return "Generate the full mystery structure."
		}
		return "Focus on " + value + "."
	case "save":
		return "Generate NPCs for these roles: " + value + "."
	case "map":
		if value == "full-map" {
			return "Generate all 10 locations for the map."
		}
		return "Focus on " + value + "."
	case "final":
		if value == "full-doc" {
			return "Generate the complete adventure document."
		}
		return "Generate a one-page summary of the adventure."
	}
	return value
}

// SystemPrompt returns the LLM system prompt for a given step.
func SystemPrompt(step string) string {
	base := `You are a collaborative Warden's assistant for the Mothership RPG (sci-fi horror tabletop RPG). The user is the Warden (game master) preparing an adventure framework for a future session — they are not a player. Your job is to help define the situation, atmosphere, NPCs, hooks, and hidden truths that the Warden will present to players at the table. Never ask questions that only players can answer (e.g., what the crew decides to do, who volunteers, what choices they make). Focus on what the Warden needs to decide ahead of time. Be concise and evocative — use the gritty, terse tone of the Mothership universe. When asking follow-up questions, ask only ONE at a time.`

	contextNote := ` If context from prior steps is provided, treat those decisions as established — do not re-ask what was already decided. Open by briefly acknowledging what's been locked in, confirm any choice that follows naturally from prior decisions, then focus the conversation on what this step uniquely contributes.`

	prompts := map[string]string{
		"scenario": base + ` The user is choosing a starting scenario. Help them flesh it out: who hired the crew, what's the immediate hook, what feels wrong from the start. This step sets the premise — keep it focused on situation and motivation, not location or the Horror.`,

		"setting": base + contextNote + ` This step is about the physical space — not what's happening, but where. Make it specific and atmospheric: current state of the location, layout, who's present, what's visibly wrong. If the scenario already implies a setting type, confirm it and move straight to developing the details rather than re-asking the obvious.`,

		"tombs": base + contextNote + ` The user is defining the Horror using the TOMBS cycle (Transgression, Omens, Manifestation, Banishment, Slumber). When rolls are provided, synthesize them into a coherent, terrifying Horror with a name and nature. Tie it to the established scenario and setting — the Horror should feel inevitable given what's already decided.`,

		"survive": base + contextNote + ` This step defines the physical or psychological pressure that forces hard choices — not the mystery, not the NPCs, just the relentless survival threat. Players can usually only do one of: Survive, Solve, or Save. Ground the threat in the established Horror and setting. What is actively trying to kill or break the crew?`,

		"solve": base + contextNote + ` This step builds the mystery using three parts — work through all three before the step is complete. The three questions form a chain, each unlocking the next: (1) "What happened here?" establishes the event or truth the players uncover. (2) "Who did it?" — once they know what happened, who caused it? (3) "Where are they?" — once they know who, where is that person, creature, or thing now? Work through all three in order, making each answer specific to the established Horror, setting, and scenario. Then define the PUZZLES — 2-3 concrete obstacles between the players and those answers (a locked area, an uncooperative NPC, a missing piece of evidence). Finally define the CLUES — specific things in the world (objects, people, places) that point toward each answer. When a clue is found it should raise a new question. The Horror is known to the Warden but hidden from players — focus on what evidence exists and how players might encounter it.`,

		"save": base + contextNote + ` This step creates the NPCs. Use the character matrix: Gatekeeper (powerful/unhelpful), Employer, Benefactor, Traitor, Survivor, Expert, Coward, Victim, Drinking Buddy, Wildcard. Each NPC needs a name, role in the matrix, one defining physical or behavioral trait, and one concrete want that puts them in conflict or alignment with the crew. Root them in the established setting and scenario.`,

		"map": base + contextNote + ` This step maps the adventure location as 10 numbered boxes in a rough flowchart. Use the established setting, Horror, survival threat, mystery, and NPCs as your raw material — every significant element from prior steps should have a physical home on the map. Box 1 is the crew's entry point. The Horror's lair must appear. Place the survival threat in specific locations (e.g., if the threat is environmental, which areas are affected?). Give NPCs starting locations. Include at least one locked or secret area tied to the mystery's clues. Include a defensible safe zone. Each box gets a short name and 1-2 sentences grounded in the specific details already established — not generic sci-fi rooms. Do not invent new lore; express what has already been decided spatially.`,

		"final": base + contextNote + ` Generate the complete adventure document using everything established across all steps. Sections: Session Title, Scenario, Setting, The Horror (TOMBS elements woven into narrative), Something to Survive, Something to Solve, Someone to Save (full NPC details), The Map (all 10 locations). Terse and evocative throughout — this is a Warden's reference, not a novel.`,
	}
	if p, ok := prompts[step]; ok {
		return p
	}
	return "You are a collaborative Warden's assistant for the Mothership RPG. Help the user develop their adventure."
}
