package queries

import (
	"fmt"
	"os"
	"path/filepath"
)

type Store struct {
	dir string
}

func NewStore(dir string) Store {
	return Store{
		dir: dir,
	}
}

func (s Store) ReadAll() ([]Query, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	queries := []Query{}

	for _, entry := range entries {
		name := entry.Name()
		path := filepath.Join(s.dir, name)

		fmt.Println("==>", name)

		if filepath.Ext(name) != ".sql" {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		dataStr := string(data)

		meta, err := parseMetadata(dataStr)
		if err != nil {
			return nil, err
		}
		if meta == nil {
			continue
		}

		queries = append(queries, Query{
			ID:   entry.Name(),
			Path: path,
			Meta: meta,
			Data: dataStr,
		})
	}

	return queries, nil
}
