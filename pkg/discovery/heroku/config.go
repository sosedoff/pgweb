package heroku

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bgentry/go-netrc/netrc"
	"github.com/sosedoff/pgweb/pkg/command"
)

func netrcPath() string {
	return filepath.Join(os.Getenv("HOME"), ".netrc")
}

func readConfig(opts *command.Options) {
	log.Println("[heroku] reading configuration from ~/.netrc file")

	machine, err := netrc.FindMachine(netrcPath(), "api.heroku.com")
	if err != nil {
		log.Println("[heroku] cant read netrc file:", err)
		return
	}
	if machine == nil {
		log.Println("[heroku] api.heroku.com section is not found in .netrc file")
		return
	}

	opts.HerokuToken = machine.Password
}
