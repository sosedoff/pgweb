package queries

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrQueryDirNotExist  = errors.New("queries directory does not exist")
	ErrQueryFileNotExist = errors.New("query file does not exist")
)

type Store struct {
	dir string
}

func NewStore(dir string) *Store {
	return &Store{
		dir: dir,
	}
}

func (s Store) Read(id string) (*Query, error) {
	path := filepath.Join(s.dir, fmt.Sprintf("%s.sql", id))
	return readQuery(path)
}

func (s Store) ReadAll() ([]Query, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = ErrQueryDirNotExist
		}
		return nil, err
	}

	queries := []Query{}

	for _, entry := range entries {
		name := entry.Name()
		if filepath.Ext(name) != ".sql" {
			continue
		}

		path := filepath.Join(s.dir, name)
		query, err := readQuery(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[WARN] skipping %q query file due to error: %v\n", name, err)
			continue
		}
		if query == nil {
			continue
		}

		queries = append(queries, *query)
	}

	return queries, nil
}

func readQuery(path string) (*Query, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrQueryFileNotExist
		}
		return nil, err
	}
	dataStr := string(data)

	meta, err := parseMetadata(dataStr)
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, nil
	}

	return &Query{
		ID:   strings.Replace(filepath.Base(path), ".sql", "", 1),
		Path: path,
		Meta: meta,
		Data: sanitizeMetadata(dataStr),
	}, nil
}
