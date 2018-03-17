package statements

import (
	"errors"
)

// DatabaseSql is the sql queries interface
type DatabaseDialect struct {
	Databases *string
	Schemas *string
	Info * string
	TableIndexes *string
	TableConstraints *string
	TableInfo *string
	TableSchema *string
	MaterializedView *string
	Objects *string
	Activity *map[string]string
	Provider string
}


var (
	sqlDialects  = make(map[string]*DatabaseDialect)
	errorAlreadyRegistered = errors.New("Dialect already registered")
	errorNotFound = errors.New("Dialect not found")
)

// RegisterDialect register new available sql dialect
func RegisterDialect( name string, dialect* DatabaseDialect ) error {
	if _, ok := sqlDialects[name]; ok {
		return errorAlreadyRegistered
	}
	sqlDialects[name] = dialect
	return nil
}


// FindDialect search dialect in registered dialect list
func FindDialect( name string ) (*DatabaseDialect, error) {
	if dialect, ok := sqlDialects[name]; !ok {
		return nil,errorNotFound
	} else {
		return dialect, nil
	}
}
