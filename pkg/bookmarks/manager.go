package bookmarks

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Manager struct {
	dir string
}

func NewManager(dir string) Manager {
	return Manager{
		dir: dir,
	}
}

func (m Manager) Get(id string) (*Bookmark, error) {
	bookmarks, err := m.list()
	if err != nil {
		return nil, err
	}

	for _, b := range bookmarks {
		if b.ID == id {
			return &b, nil
		}
	}

	return nil, fmt.Errorf("bookmark %v not found", id)
}

func (m Manager) List() ([]Bookmark, error) {
	return m.list()
}

func (m Manager) ListIDs() ([]string, error) {
	bookmarks, err := m.list()
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(bookmarks))
	for i, bookmark := range bookmarks {
		ids[i] = bookmark.ID
	}

	return ids, nil
}

func (m Manager) list() ([]Bookmark, error) {
	result := []Bookmark{}

	if m.dir == "" {
		return result, nil
	}

	info, err := os.Stat(m.dir)
	if err != nil {
		// Do not fail if base dir does not exists: it's not created by default
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "[WARN] bookmarks dir %s does not exist\n", m.dir)
			return result, nil
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path %s is not a directory", m.dir)
	}

	dirEntries, err := os.ReadDir(m.dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		name := entry.Name()
		if filepath.Ext(name) != ".toml" {
			continue
		}

		bookmark, err := readBookmark(filepath.Join(m.dir, name))
		if err != nil {
			// Do not fail if one of the bookmarks is invalid
			fmt.Fprintf(os.Stderr, "[WARN] bookmark file %s is invalid: %s\n", name, err)
			continue
		}

		result = append(result, bookmark)
	}

	return result, nil
}

func readBookmark(path string) (Bookmark, error) {
	bookmark := Bookmark{
		ID: fileBasename(path),
	}

	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = fmt.Errorf("bookmark file %s does not exist", path)
		}
		return bookmark, err
	}

	buff, err := os.ReadFile(path)
	if err != nil {
		return bookmark, err
	}

	_, err = toml.Decode(string(buff), &bookmark)

	if bookmark.Port == 0 {
		bookmark.Port = 5432
	}

	// List of all supported postgres modes
	modes := []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}
	valid := false

	for _, mode := range modes {
		if bookmark.SSLMode == mode {
			valid = true
			break
		}
	}

	// Fall back to a default mode if mode is not set or invalid
	// Typical typo: ssl mode set to "disabled"
	if bookmark.SSLMode == "" || !valid {
		bookmark.SSLMode = "disable"
	}

	// Set default SSH port if it's not provided by user
	if bookmark.SSH != nil && bookmark.SSH.Port == "" {
		bookmark.SSH.Port = "22"
	}

	return bookmark, err
}

func fileBasename(path string) string {
	filename := filepath.Base(path)
	return strings.Replace(filename, filepath.Ext(path), "", 1)
}
