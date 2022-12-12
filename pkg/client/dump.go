package client

import (
	"bytes"
	"context"
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

// Validate checks availability and version of pg_dump CLI
func (d *Dump) Validate(serverVersion string) error {
	out := bytes.NewBuffer(nil)

	cmd := exec.Command("pg_dump", "--version")
	cmd.Stdout = out
	cmd.Stderr = out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_dump command failed: %s", out.Bytes())
	}

	detected, dumpVersion := detectDumpVersion(out.String())
	if detected && serverVersion != "" {
		satisfied := checkVersionRequirement(dumpVersion, serverVersion)
		if !satisfied {
			return fmt.Errorf("pg_dump version %v not compatible with server version %v", dumpVersion, serverVersion)
		}
	}

	return nil
}

// Export streams the database dump to the specified writer
func (d *Dump) Export(ctx context.Context, connstr string, writer io.Writer) error {
	if str, err := removeUnsupportedOptions(connstr); err != nil {
		return err
	} else {
		connstr = str
	}

	opts := []string{
		"--no-owner",      // skip restoration of object ownership in plain-text format
		"--clean",         // clean (drop) database objects before recreating
		"--compress", "6", // compression level for compressed formats
	}

	if d.Table != "" {
		opts = append(opts, []string{"--table", d.Table}...)
	}

	opts = append(opts, connstr)
	errOutput := bytes.NewBuffer(nil)

	cmd := exec.CommandContext(ctx, "pg_dump", opts...)
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
