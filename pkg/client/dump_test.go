package client

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func test_DumpExport(t *testing.T) {
	url := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)

	savePath := "/tmp/dump.sql.gz"
	os.Remove(savePath)

	saveFile, err := os.Create(savePath)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer func() {
		saveFile.Close()
		os.Remove(savePath)
	}()

	// Test full db dump
	dump := Dump{}
	err = dump.Export(url, saveFile)
	assert.NoError(t, err)

	// Test nonexistent database
	invalidUrl := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, "foobar")
	err = dump.Export(invalidUrl, saveFile)
	assert.Contains(t, err.Error(), `database "foobar" does not exist`)

	// Test dump of non existent db
	dump = Dump{Table: "foobar"}
	err = dump.Export(url, saveFile)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "pg_dump: no matching tables were found")
}
