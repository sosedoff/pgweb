package discovery

import (
	"github.com/sosedoff/pgweb/pkg/bookmarks"
)

type (
	// Provider interface represents a third-party provider integration
	Provider interface {
		ID() string
		Name() string
		List() ([]Resource, error)
		Get(string) (*bookmarks.Bookmark, error)
	}

	// Resource represents a third-party database resource information.
	Resource struct {
		ID   string                 `json:"id"`
		Name string                 `json:"name"`
		Meta map[string]interface{} `json:"meta,omitempty"`
	}
)
