package cli

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sosedoff/pgweb/pkg/api"
	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/command"

	"github.com/stretchr/testify/assert"
)

var (
	testCommands   map[string]string
	serverHost     string
	serverPort     string
	serverUser     string
	serverPassword string
	serverDatabase string
)

func mapKeys(data map[string]*client.Objects) []string {
	result := []string{}
	for k, _ := range data {
		result = append(result, k)
	}
	return result
}

func pgVersion() (int, int) {
	var major, minor int
	fmt.Sscanf(os.Getenv("PGVERSION"), "%d.%d", &major, &minor)
	return major, minor
}

func getVar(name, def string) string {
	val := os.Getenv(name)
	if val == "" {
		return def
	}
	return val
}

func initVars() {
	// We need to load default options to make sure all stuff works
	if err := command.SetDefaultOptions(); err != nil {
		log.Fatal(err)
	}
	command.Opts.HttpPort = 8081

	serverHost = getVar("PGHOST", "localhost")
	serverPort = getVar("PGPORT", "5432")
	serverUser = getVar("PGUSER", "postgres")
	serverPassword = getVar("PGPASSWORD", "postgres")
	serverDatabase = getVar("PGDATABASE", "booktown")
}

func setupCommands() {
	testCommands = map[string]string{
		"createdb": "createdb",
		"psql":     "psql",
		"dropdb":   "dropdb",
	}

	if onWindows() {
		for k, v := range testCommands {
			testCommands[k] = v + ".exe"
		}
	}
}

func onWindows() bool {
	return runtime.GOOS == "windows"
}

func setup() {
	// No pretty JSON for testsm
	options = command.Opts
	options.DisablePrettyJson = true

	out, err := exec.Command(
		testCommands["createdb"],
		"-U", serverUser,
		"-h", serverHost,
		"-p", serverPort,
		serverDatabase,
	).CombinedOutput()

	if err != nil {
		fmt.Println("Database creation failed:", string(out))
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	out, err = exec.Command(
		testCommands["psql"],
		"-U", serverUser,
		"-h", serverHost,
		"-p", serverPort,
		"-f", "../../data/booktown.sql",
		serverDatabase,
	).CombinedOutput()

	if err != nil {
		fmt.Println("Database import failed:", string(out))
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func setupClient() {
	url := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)
	api.DbClient, _ = client.NewFromUrl(url, nil)
}

func teardownClient() {
	if api.DbClient != nil {
		api.DbClient.Close()
	}
}

func teardown() {
	_, err := exec.Command(
		testCommands["dropdb"],
		"-U", serverUser,
		"-h", serverHost,
		"-p", serverPort,
		serverDatabase,
	).CombinedOutput()

	if err != nil {
		fmt.Println("Teardown error:", err)
	}
}

func testDataImportCsv(t *testing.T) {
	var client *http.Client
	var remoteURL string
	client = &http.Client{Timeout: time.Second * 10}
	remoteURL = "http://localhost:8081/api/import/csv"

	err := Upload(client, remoteURL)
	if err != nil {
		panic(err)
	}
}

func Upload(client *http.Client, url string) (err error) {
	// Prepare a form that you will submit to that URL.
	//prepare the reader instances to encode
	values := map[string]io.Reader{
		"importCSVFile":           mustOpen("test.csv"), // lets assume its this file
		"importCSVFieldDelimiter": strings.NewReader(","),
		"importCSVTableName":      strings.NewReader("from_csv")}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

// this test is coordinated with test.csv
func testResultCsv(t *testing.T) {
	res, _ := api.DbClient.Query("SELECT * FROM from_csv ORDER BY id")
	csv := res.CSV()

	expected := "id,line\n1,line 1\n1,line 1\n2,line-2\n2,line-2\n"

	assert.Equal(t, expected, string(csv))
}

func TestAll(t *testing.T) {
	if onWindows() {
		t.Log("Unit testing on Windows platform is not supported.")
		return
	}

	initVars()
	setupCommands()
	teardown()
	setup()
	setupClient()
	StartServer()

	testDataImportCsv(t)
	testDataImportCsv(t)
	testResultCsv(t)

	teardownClient()
	teardown()
}
