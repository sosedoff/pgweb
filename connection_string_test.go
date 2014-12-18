package main

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Invalid_Url(t *testing.T) {
	opts := Options{}
	examples := []string{
		"postgresql://foobar",
		"foobar",
	}

	for _, val := range examples {
		opts.Url = val
		str, err := buildConnectionString(opts)

		assert.Equal(t, "", str)
		assert.Error(t, err)
		assert.Equal(t, "Invalid URL. Valid format: postgres://user:password@host:port/db?sslmode=mode", err.Error())
	}
}

func Test_Valid_Url(t *testing.T) {
	url := "postgres://myhost/database"
	str, err := buildConnectionString(Options{Url: url})

	assert.Equal(t, nil, err)
	assert.Equal(t, url, str)
}

func Test_Url_And_Ssl_Flag(t *testing.T) {
	str, err := buildConnectionString(Options{
		Url: "postgres://myhost/database",
		Ssl: "disable",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://myhost/database?sslmode=disable", str)
}

func Test_Localhost_Url_And_No_Ssl_Flag(t *testing.T) {
	str, err := buildConnectionString(Options{
		Url: "postgres://localhost/database",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://localhost/database?sslmode=disable", str)

	str, err = buildConnectionString(Options{
		Url: "postgres://127.0.0.1/database",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://127.0.0.1/database?sslmode=disable", str)
}

func Test_Localhost_Url_And_Ssl_Flag(t *testing.T) {
	str, err := buildConnectionString(Options{
		Url: "postgres://localhost/database",
		Ssl: "require",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://localhost/database?sslmode=require", str)

	str, err = buildConnectionString(Options{
		Url: "postgres://127.0.0.1/database",
		Ssl: "require",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://127.0.0.1/database?sslmode=require", str)
}

func Test_Localhost_Url_And_Ssl_Arg(t *testing.T) {
	str, err := buildConnectionString(Options{
		Url: "postgres://localhost/database?sslmode=require",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://localhost/database?sslmode=require", str)

	str, err = buildConnectionString(Options{
		Url: "postgres://127.0.0.1/database?sslmode=require",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://127.0.0.1/database?sslmode=require", str)
}

func Test_Flag_Args(t *testing.T) {
	str, err := buildConnectionString(Options{
		Host:   "host",
		Port:   5432,
		User:   "user",
		Pass:   "password",
		DbName: "db",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://user:password@host:5432/db", str)
}

func Test_Localhost(t *testing.T) {
	opts := Options{
		Host:   "localhost",
		Port:   5432,
		User:   "user",
		Pass:   "password",
		DbName: "db",
	}

	str, err := buildConnectionString(opts)
	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://user:password@localhost:5432/db?sslmode=disable", str)

	opts.Host = "127.0.0.1"
	str, err = buildConnectionString(opts)
	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://user:password@127.0.0.1:5432/db?sslmode=disable", str)
}

func Test_Localhost_And_Ssl(t *testing.T) {
	opts := Options{
		Host:   "localhost",
		Port:   5432,
		User:   "user",
		Pass:   "password",
		DbName: "db",
		Ssl:    "require",
	}

	str, err := buildConnectionString(opts)
	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://user:password@localhost:5432/db?sslmode=require", str)
}

func Test_No_User(t *testing.T) {
	opts := Options{Host: "host", Port: 5432, DbName: "db"}
	u, _ := user.Current()
	str, err := buildConnectionString(opts)

	assert.Equal(t, nil, err)
	assert.Equal(t, fmt.Sprintf("postgres://%s@host:5432/db", u.Username), str)
}

func Test_Port(t *testing.T) {
	opts := Options{Host: "host", User: "user", Port: 5000, DbName: "db"}
	str, err := buildConnectionString(opts)

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://user@host:5000/db", str)
}
