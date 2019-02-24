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
	"errors"

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
	auxCloser      chan int
	serviceUrl     string
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
	serviceUrl = "http://localhost:8081"
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

func setupDatabase() error {
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
		return errors.New(fmt.Sprintf("Create db failed. Error message: «%s», OS command output: «%s»", err.Error(), string(out)))
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
		return errors.New(fmt.Sprintf("Db import failed. Error message: «%s», OS command output: «%s»", err.Error(), string(out)))
	}
	return nil
}

func setupServer() {
	auxCloser = make(chan int)
	go Run(auxCloser)
}

func teardownServer() {
	closer := func() { 
		auxCloser <- 1
	}
	go closer()
}

type formDataType map[string]io.Reader

func setupClient() (err error) {
	// url := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)
	// Generate session id
	// Login with this url
	// Assert success

	var client *http.Client = &http.Client{Timeout: time.Second * 10}
	apiUrl := serviceUrl + "/api/connect"
	postgresUrlString := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)

	formData := formDataType{
		"url":           strings.NewReader(postgresUrlString), // lets assume its this file
	}
 var req *http.Request
	req, err = preparePostRequest(apiUrl, formData)
	if err != nil {
		return
	}
	
	// Submit the request
	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		return
	}
	
	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}

func teardownClient() (err error) {
	// disconnect here
	var client *http.Client = &http.Client{Timeout: time.Second * 1000}
 var req *http.Request
	req, err = preparePostRequest(serviceUrl + "/api/disconnect", formDataType{})
	if err != nil {
		return
	}
	
	// Submit the request
	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		return
	}
	
	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}




func teardownDatabase() error {
	out, err := exec.Command(
		testCommands["dropdb"],
		"-U", serverUser,
		"-h", serverHost,
		"-p", serverPort,
		serverDatabase,
	).CombinedOutput()

	if err != nil {
		return errors.New(fmt.Sprintf("Dropdb failed. Error message: «%s», drop db command output: «%s»", err.Error(), string(out)))
	}
	return nil
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
	values := formDataType{
		"importCSVFile":           mustOpen("test.csv"), // lets assume its this file
		"importCSVFieldDelimiter": strings.NewReader(","),
		"importCSVTableName":      strings.NewReader("from_csv")}

 req, err := preparePostRequest("POST", values)
	// Now that you have a form, you can submit it to your handler.
	if err != nil {
		return
	}

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


func preparePostRequest(apiUrl string, formData map[string]io.Reader) (req *http.Request, err error) {
	req = nil
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range formData {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
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
			return 
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()
	req, err = http.NewRequest("POST", apiUrl, &b)

	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	sessionId := "test-sess-ion-id"
	req.Header.Add("x-session-id", sessionId)

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

// returns true if there is an error
func reportIfErr(stepName string,err error) bool {
	if err != nil {
		fmt.Println("Step ",stepName," error: "+err.Error())
		return true
	}
	return false
}

func TestAll(t *testing.T) {
	if onWindows() {
		t.Log("Unit testing on Windows platform is not supported.")
		return
	}
	initVars()
	setupCommands()
	// We expect that database does not exist, so we ignore errors here
	teardownDatabase()
	if reportIfErr("setupDatabase",setupDatabase()) {
		return
	}
	defer func(){
		reportIfErr("teardownDatabase",teardownDatabase())
	}()
	
	
	setupServer()
	defer teardownServer()


	// FIXME there must be a better way to wait for server to start, e.g. 
	time.Sleep(5 * time.Second)
	if reportIfErr("setupClient",setupClient()) {
		return
	}
	defer func(){
		reportIfErr("teardownClient",teardownClient())
	}()

	//testDataImportCsv(t)
	//testDataImportCsv(t)
	//testResultCsv(t)

}
