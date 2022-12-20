package connection

import (
	"fmt"
	"net/url"
	"os/user"
	"testing"

	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/stretchr/testify/assert"
)

func TestBuildStringFromOptions(t *testing.T) {
	t.Run("valid url", func(t *testing.T) {
		url := "postgres://myhost/database"
		str, err := BuildStringFromOptions(command.Options{URL: url})

		assert.NoError(t, err)
		assert.Equal(t, url, str)
	})

	t.Run("with sslmode param", func(t *testing.T) {
		str, err := BuildStringFromOptions(command.Options{
			URL:     "postgres://myhost/database",
			SSLMode: "disable",
		})

		assert.NoError(t, err)
		assert.Equal(t, "postgres://myhost/database?sslmode=disable", str)
	})

	t.Run("sets sslmode param if not set", func(t *testing.T) {
		str, err := BuildStringFromOptions(command.Options{
			URL: "postgres://localhost/database",
		})
		assert.NoError(t, err)
		assert.Equal(t, "postgres://localhost/database?sslmode=disable", str)

		str, err = BuildStringFromOptions(command.Options{
			URL: "postgres://127.0.0.1/database",
		})
		assert.NoError(t, err)
		assert.Equal(t, "postgres://127.0.0.1/database?sslmode=disable", str)
	})

	t.Run("sslmode as an option", func(t *testing.T) {
		str, err := BuildStringFromOptions(command.Options{
			URL:     "postgres://localhost/database",
			SSLMode: "require",
		})
		assert.NoError(t, err)
		assert.Equal(t, "postgres://localhost/database?sslmode=require", str)

		str, err = BuildStringFromOptions(command.Options{
			URL:     "postgres://127.0.0.1/database",
			SSLMode: "require",
		})
		assert.NoError(t, err)
		assert.Equal(t, "postgres://127.0.0.1/database?sslmode=require", str)
	})

	t.Run("localhost and sslmode flag", func(t *testing.T) {
		str, err := BuildStringFromOptions(command.Options{
			URL: "postgres://localhost/database?sslmode=require",
		})
		assert.NoError(t, err)
		assert.Equal(t, "postgres://localhost/database?sslmode=require", str)

		str, err = BuildStringFromOptions(command.Options{
			URL: "postgres://127.0.0.1/database?sslmode=require",
		})
		assert.NoError(t, err)
		assert.Equal(t, "postgres://127.0.0.1/database?sslmode=require", str)
	})

	t.Run("extended options", func(t *testing.T) {
		str, err := BuildStringFromOptions(command.Options{
			URL: "postgres://localhost/database?sslmode=require&sslcert=cert&sslkey=key&sslrootcert=ca",
		})
		assert.NoError(t, err)
		assert.Equal(t, "postgres://localhost/database?sslcert=cert&sslkey=key&sslmode=require&sslrootcert=ca", str)
	})

	t.Run("from flags", func(t *testing.T) {
		str, err := BuildStringFromOptions(command.Options{
			Host:   "host",
			Port:   5432,
			User:   "user",
			Pass:   "password",
			DbName: "db",
		})

		assert.NoError(t, err)
		assert.Equal(t, "postgres://user:password@host:5432/db", str)
	})

	t.Run("localhost", func(t *testing.T) {
		opts := command.Options{
			Host:   "localhost",
			Port:   5432,
			User:   "user",
			Pass:   "password",
			DbName: "db",
		}

		str, err := BuildStringFromOptions(opts)
		assert.NoError(t, err)
		assert.Equal(t, "postgres://user:password@localhost:5432/db?sslmode=disable", str)

		opts.Host = "127.0.0.1"
		str, err = BuildStringFromOptions(opts)
		assert.NoError(t, err)
		assert.Equal(t, "postgres://user:password@127.0.0.1:5432/db?sslmode=disable", str)
	})

	t.Run("localhost and ssl", func(t *testing.T) {
		opts := command.Options{
			Host:        "localhost",
			Port:        5432,
			User:        "user",
			Pass:        "password",
			DbName:      "db",
			SSLMode:     "require",
			SSLKey:      "keyPath",
			SSLCert:     "certPath",
			SSLRootCert: "caPath",
		}

		str, err := BuildStringFromOptions(opts)
		assert.NoError(t, err)
		assert.Equal(t, "postgres://user:password@localhost:5432/db?sslcert=certPath&sslkey=keyPath&sslmode=require&sslrootcert=caPath", str)
	})

	t.Run("no user", func(t *testing.T) {
		opts := command.Options{Host: "host", Port: 5432, DbName: "db"}
		u, _ := user.Current()
		str, err := BuildStringFromOptions(opts)
		userAndPass := url.UserPassword(u.Username, "").String()

		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("postgres://%s@host:5432/db", userAndPass), str)
	})

	t.Run("port", func(t *testing.T) {
		opts := command.Options{Host: "host", User: "user", Port: 5000, DbName: "db"}
		str, err := BuildStringFromOptions(opts)
		assert.NoError(t, err)
		assert.Equal(t, "postgres://user:@host:5000/db", str)
	})

	t.Run("with pgpass", func(t *testing.T) {
		opts := command.Options{
			Host:     "localhost",
			Port:     5432,
			User:     "username",
			DbName:   "dbname",
			Passfile: "../../data/passfile",
		}

		str, err := BuildStringFromOptions(opts)
		assert.NoError(t, err)
		assert.Equal(t, "postgres://username:password@localhost:5432/dbname?sslmode=disable", str)

		opts.User = "foobar"
		str, err = BuildStringFromOptions(opts)
		assert.NoError(t, err)
		assert.Equal(t, "postgres://foobar:@localhost:5432/dbname?sslmode=disable", str)

		opts.Host = "127.0.0.1"
		opts.DbName = "foobar2"
		str, err = BuildStringFromOptions(opts)
		assert.NoError(t, err)
		assert.Equal(t, "postgres://foobar:password2@127.0.0.1:5432/foobar2?sslmode=disable", str)
	})

	t.Run("with connection timeout", func(t *testing.T) {
		opts := command.Options{
			URL:         "postgres://user:pass@localhost:5432/dbname",
			OpenTimeout: 30,
		}
		str, err := BuildStringFromOptions(opts)
		assert.NoError(t, err)
		assert.Equal(t, "postgres://user:pass@localhost:5432/dbname?connect_timeout=30&sslmode=disable", str)

		opts = command.Options{
			Host:        "localhost",
			Port:        5432,
			User:        "username",
			DbName:      "dbname",
			OpenTimeout: 30,
		}

		str, err = BuildStringFromOptions(opts)
		assert.NoError(t, err)
		assert.Equal(t, "postgres://username:@localhost:5432/dbname?connect_timeout=30&sslmode=disable", str)
	})

	t.Run("invalid url", func(t *testing.T) {
		opts := command.Options{}
		examples := []string{
			"postgre://foobar",
			"tcp://blah",
			"foobar",
		}

		for _, val := range examples {
			opts.URL = val
			str, err := BuildStringFromOptions(opts)

			assert.Equal(t, "", str)
			assert.Error(t, err)
			assert.Equal(t, "Invalid URL. Valid format: postgres://user:password@host:port/db?sslmode=mode", err.Error())
		}
	})
}

func TestFormatURL(t *testing.T) {
	examples := []struct {
		name   string
		input  command.Options
		result string
		err    string
	}{
		{
			name:  "empty opts",
			input: command.Options{},
		},
		{
			name:  "invalid url",
			input: command.Options{URL: "barurl"},
			err:   "Invalid URL",
		},
		{
			name: "good",
			input: command.Options{
				URL: "postgres://user:pass@localhost:5432/dbname",
			},
			result: "postgres://user:pass@localhost:5432/dbname?sslmode=disable",
		},
		{
			name: "password lookup, password set",
			input: command.Options{
				URL:      "postgres://username:@localhost:5432/dbname",
				Passfile: "../../data/passfile",
			},
			result: "postgres://username:password@localhost:5432/dbname?sslmode=disable",
		},
		{
			name: "password lookup, password not set",
			input: command.Options{
				URL:      "postgres://username@localhost:5432/dbname",
				Passfile: "../../data/passfile",
			},
			result: "postgres://username:password@localhost:5432/dbname?sslmode=disable",
		},
		{
			name: "with timeout setting",
			input: command.Options{
				URL:         "postgres://username@localhost:5432/dbname",
				OpenTimeout: 30,
			},
			result: "postgres://username@localhost:5432/dbname?connect_timeout=30&sslmode=disable",
		},
	}

	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			str, err := FormatURL(ex.input)

			if ex.err != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), ex.err)
			}
			assert.Equal(t, ex.result, str)
		})
	}
}

func TestIsBlank(t *testing.T) {
	assert.Equal(t, true, IsBlank(command.Options{}))
	assert.Equal(t, false, IsBlank(command.Options{Host: "host", User: "user"}))
	assert.Equal(t, false, IsBlank(command.Options{Host: "host", User: "user", DbName: "db"}))
	assert.Equal(t, false, IsBlank(command.Options{URL: "url"}))
}
