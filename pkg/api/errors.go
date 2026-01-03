package api

import (
	"errors"
)

var (
	errNotConnected         = errors.New("Not connected")
	errNotPermitted         = errors.New("Not permitted")
	errInvalidConnString    = errors.New("Invalid connection string")
	errSessionRequired      = errors.New("Session ID is required")
	errSessionLocked        = errors.New("Session is locked")
	errURLRequired          = errors.New("URL parameter is required")
	errQueryRequired        = errors.New("Query parameter is required")
	errDatabaseNameRequired = errors.New("Database name is required")
)
