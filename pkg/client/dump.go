package client

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

type Dump struct {
	Table string
}

func (d *Dump) Export(url string, writer io.Writer) error {
	errOutput := bytes.NewBuffer(nil)

	opts := []string{
		"--no-owner",      // skip restoration of object ownership in plain-text format
		"--clean",         // clean (drop) database objects before recreating
		"--compress", "6", // compression level for compressed formats
	}

	if d.Table != "" {
		opts = append(opts, []string{"--table", d.Table}...)
	}

	opts = append(opts, url)

	cmd := exec.Command("pg_dump", opts...)
	cmd.Stdout = writer
	cmd.Stderr = errOutput

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error: %s. output: %s", err.Error(), errOutput.Bytes())
	}
	return nil
}
