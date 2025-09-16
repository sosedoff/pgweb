package connect

// Request holds the resource request details
type Request struct {
	Resource string            `json:"resource"`
	Token    string            `json:"token"`
	Headers  map[string]string `json:"headers,omitempty"`
}

// Credential holds the database connection string
type Credential struct {
	DatabaseURL string `json:"database_url"`
}
