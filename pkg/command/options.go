package command

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Version  bool   `short:"v" long:"version" description:"Print version"`
	Debug    bool   `short:"d" long:"debug" description:"Enable debugging mode" default:"false"`
	Url      string `long:"url" description:"Database connection string"`
	Host     string `long:"host" description:"Server hostname or IP"`
	Port     int    `long:"port" description:"Server port" default:"5432"`
	User     string `long:"user" description:"Database user"`
	Pass     string `long:"pass" description:"Password for user"`
	DbName   string `long:"db" description:"Database name"`
	Ssl      string `long:"ssl" description:"SSL option"`
	HttpHost string `long:"bind" description:"HTTP server host" default:"localhost"`
	HttpPort uint   `long:"listen" description:"HTTP server listen port" default:"8081"`
	AuthUser string `long:"auth-user" description:"HTTP basic auth user"`
	AuthPass string `long:"auth-pass" description:"HTTP basic auth password"`
	SkipOpen bool   `short:"s" long:"skip-open" description:"Skip browser open on start"`
	Sessions bool   `long:"sessions" description:"Enable multiple database sessions" default:"false"`
}

var Opts Options

func ParseOptions() error {
	_, err := flags.ParseArgs(&Opts, os.Args)
	if err != nil {
		return err
	}

	if Opts.Url == "" {
		Opts.Url = os.Getenv("DATABASE_URL")
	}

	if os.Getenv("SESSIONS") != "" {
		Opts.Sessions = true
	}

	return nil
}
