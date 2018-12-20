package bookmark

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/sosedoff/pgweb/pkg/bookmarks"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/discovery"
)

// Provider is a bookmarks discovery provider
type Provider struct {
	basePath string
}

// New returns a new bookmarks provider
func New(opts command.Options) (*Provider, error) {
	return &Provider{
		basePath: bookmarksDir(),
	}, nil
}

// ID returns the provider identificator
func (p Provider) ID() string {
	return "bookmarks"
}

// Name returns the provider name
func (p Provider) Name() string {
	return "Bookmarks"
}

// List returns list of all bookmarks
func (p Provider) List() ([]discovery.Resource, error) {
	files, err := ioutil.ReadDir(p.basePath)
	if err != nil {
		return nil, err
	}

	resources := []discovery.Resource{}

	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".toml" {
			continue
		}

		name := f.Name()

		resources = append(resources, discovery.Resource{
			ID:   name,
			Name: strings.Replace(name, filepath.Ext(name), "", 1),
		})
	}

	return resources, nil
}

// Get returns a database credential for the given resource ID
func (p Provider) Get(id string) (*bookmarks.Bookmark, error) {
	b, err := readBookmark(filepath.Join(p.basePath, id))
	if err != nil {
		return nil, err
	}
	return b, nil
}
