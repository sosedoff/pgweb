package client

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"strings"
)

var (
	unsupportedDumpOptions = []string{
		"search_path",
	}
)

// Dump represents a database dump
type Dump struct {
	Table string
}

// CanExport returns true if database dump tool could be used without an error
func (d *Dump) CanExport() bool {
	err := exec.Command("pg_dump", "--version").Run()
	return err == nil
}

// Export streams the database dump to the specified writer
func (d *Dump) Export(connstr string, writer io.Writer) error {
	if str, err := removeUnsupportedOptions(connstr); err != nil {
		return err
	} else {
		connstr = str
	}

	errOutput := bytes.NewBuffer(nil)

	opts := []string{
		"--no-owner",      // skip restoration of object ownership in plain-text format
		"--clean",         // clean (drop) database objects before recreating
		"--compress", "6", // compression level for compressed formats
	}

	if d.Table != "" {
		opts = append(opts, []string{"--table", d.Table}...)
	}

	opts = append(opts, connstr)

	cmd := exec.Command("pg_dump", opts...)
	cmd.Stdout = writer
	cmd.Stderr = errOutput

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error: %s. output: %s", err.Error(), errOutput.Bytes())
	}
	return nil
}

// removeUnsupportedOptions removes any options unsupported for making a db dump
func removeUnsupportedOptions(input string) (string, error) {
	uri, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	q := uri.Query()
	for _, opt := range unsupportedDumpOptions {
		q.Del(opt)
		q.Del(strings.ToUpper(opt))
	}
	uri.RawQuery = q.Encode()

	return uri.String(), nil
}
