package feedback

type Entry struct {
	Type        string `json:"type"`
	Timestamp   string `json:"timestamp"`
	Doc         string `json:"doc"`
	Section     string `json:"section"`
	Scenario    string `json:"scenario"`
	Impact      string `json:"impact"`
	Description string `json:"description"`

	Goal             string `json:"goal,omitempty"`
	DocsUsed         string `json:"docs_used,omitempty"`
	DocsMissing      string `json:"docs_missing,omitempty"`
	GapsHit          int    `json:"gaps_hit,omitempty"`
	FrictionPoints   int    `json:"friction_points,omitempty"`
	MethodologyGrade string `json:"methodology_grade,omitempty"`
	Notes            string `json:"notes,omitempty"`

	ProjectTier   string `json:"project_tier"`
	VirgilVersion string `json:"virgil_version"`
}

type ReportSection struct {
	Title   string
	Entries []ReportEntry
}

type ReportEntry struct {
	Doc      string
	Scenario string
	Count    int
	Impact   string
}
