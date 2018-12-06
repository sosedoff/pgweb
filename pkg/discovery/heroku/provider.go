package heroku

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bgentry/go-netrc/netrc"
	"github.com/heroku/heroku-go/v3"

	"github.com/sosedoff/pgweb/pkg/bookmarks"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/discovery"
)

var (
	errTokenMissing = errors.New("Heroku API token is missing")
)

// Provider represents Heroku discovery provider
type Provider struct {
	service *heroku.Service
}

// New returns a new Heroku provider instance
func New(opts command.Options) (*Provider, error) {
	// Try to locate Heroku API token using environment
	if opts.HerokuToken == "" {
		opts.HerokuToken = os.Getenv("HEROKU_TOKEN")
	}
	// Try to read token from local .netrc file
	if opts.HerokuToken == "" {
		netrcPath := filepath.Join(os.Getenv("HOME"), ".netrc")
		machine, err := netrc.FindMachine(netrcPath, "api.heroku.com")
		if machine != nil && err == nil {
			opts.HerokuToken = machine.Password
		}
	}
	// Final validation
	if opts.HerokuToken == "" {
		return nil, errTokenMissing
	}

	provider := Provider{
		service: heroku.NewService(&http.Client{
			Transport: &heroku.Transport{
				Password: opts.HerokuToken,
			},
			Timeout: time.Second * 5,
		}),
	}

	return &provider, nil
}

// ID returns the provider identificator
func (p Provider) ID() string {
	return "heroku"
}

// Name returns the provider name
func (p Provider) Name() string {
	return "Heroku"
}

// List returns list of all Heroku postgres addons
func (p Provider) List() ([]discovery.Resource, error) {
	addons, err := p.service.AddOnList(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	resources := []discovery.Resource{}
	for _, addon := range addons {
		// Only support Heroku Postgres addons for now
		if addon.AddonService.Name != "heroku-postgresql" {
			continue
		}

		resources = append(resources, discovery.Resource{
			ID:   addon.ID,
			Name: addon.App.Name,
			Meta: map[string]interface{}{
				"name": addon.Name,
				"plan": addon.Plan.Name,
			},
		})
	}

	return resources, nil
}

// Get returns a database credential for the given resource ID
func (p Provider) Get(id string) (*bookmarks.Bookmark, error) {
	// Fetch Heroku addon configuration list
	configs, err := p.service.AddOnConfigList(context.Background(), id, nil)
	if err != nil {
		return nil, err
	}

	// Filter out the database url config
	for _, c := range configs {
		if c.Name == "url" && c.Value != nil {
			return &bookmarks.Bookmark{Url: *c.Value}, nil
		}
	}

	return nil, errors.New("cant find config")
}
