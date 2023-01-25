package queries

type Query struct {
	ID   string
	Path string
	Meta *Metadata
	Data string
}

// IsPermitted returns true if a query is allowed to execute for a given db context
func (q Query) IsPermitted(host, user, database, mode string) bool {
	// All fields must be provided for matching
	if q.Meta == nil || host == "" || user == "" || database == "" || mode == "" {
		return false
	}

	meta := q.Meta

	return meta.Host.matches(host) &&
		meta.User.matches(user) &&
		meta.Database.matches(database) &&
		meta.Mode.matches(mode)
}
