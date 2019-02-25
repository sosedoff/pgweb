package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	// "github.com/sosedoff/pgweb/pkg/cli"
	"github.com/stretchr/testify/assert"
)

type formDataType map[string]io.Reader

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

func getVar(name, def string) string {
	val := os.Getenv(name)
	if val == "" {
		return def
	}
	return val
}

func initVars() {
	// We need to load default options to make sure all stuff works
	//if err := command.SetDefaultOptions(); err != nil {
	//	log.Fatal(err)
	//}
	//command.Opts.HttpPort = 8081

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
	//options = command.Opts
	//options.DisablePrettyJson = true

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
	go Run([]string{})
}

func teardownServer() {
	// do nothing
}

func setupClient() (err error) {
	// url := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)
	// Generate session id
	// Login with this url
	// Assert success

	var client *http.Client = &http.Client{Timeout: time.Second * 10}
	apiUrl := serviceUrl + "/api/connect"
	postgresUrlString := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)

	formData := formDataType{
		"url": strings.NewReader(postgresUrlString), // lets assume its this file
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
		err = decodeErrorFromHTTPResponsesBody(res)
	}
	return
}

func decodeErrorFromHTTPResponsesBody(res *http.Response) (err error) {
	var body map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&body)
	if err == nil {
		err = errors.New(fmt.Sprintf("Http response status %v, response body is %#v", res.StatusCode, body))
	}
	return
}

func teardownClient() (err error) {
	// disconnect here
	var client *http.Client = &http.Client{Timeout: time.Second * 1000}
	var req *http.Request
	req, err = preparePostRequest(serviceUrl+"/api/disconnect", formDataType{})
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
		err = decodeErrorFromHTTPResponsesBody(res)
	}
	return
}

func teardownDatabase() error {
	out, err := exec.Command(
		testCommands["dropdb"],
		"--if-exists",
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

func dataImportCSV(tableName, fieldDelimiter, fileName string) (err error) {
	var client *http.Client
	client = &http.Client{Timeout: time.Second * 10}
	apiUrl := "http://localhost:8081/api/import/csv"

	fd := formDataType{
		"importCSVTableName":      strings.NewReader(tableName),
		"importCSVFieldDelimiter": strings.NewReader(fieldDelimiter),
		"importCSVFile":           mustOpen(fileName),
	}

	var req *http.Request
	req, err = preparePostRequest(apiUrl, fd)
	// Now that you have a form, you can submit it to your handler.
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
		err = decodeErrorFromHTTPResponsesBody(res)
	}
	return
}

// Produces error if the query fails or result does not match the expectation
func errIfQueryResultMismatch(t *testing.T, query, expectedResult string) (err error) {
	var client *http.Client
	client = &http.Client{Timeout: time.Second * 10}
	apiUrl := "http://localhost:8081/api/query"

	fd := formDataType{
		"query": strings.NewReader(query),
	}

	var req *http.Request
	req, err = preparePostRequest(apiUrl, fd)
	// Now that you have a form, you can submit it to your handler.
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
		err = decodeErrorFromHTTPResponsesBody(res)
		return
	}

	defer res.Body.Close()
	var htmlData []byte
	htmlData, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	actualResult := string(htmlData)
	if expectedResult != actualResult {
		err = errors.New(fmt.Sprintf("Query «%s», expected «%s», actual «%s»", query, expectedResult, actualResult))
		return
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

// returns true if there is an error
func reportIfErr(err error) bool {
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func testDataImportCSVSimple(t *testing.T) {
	// we import the same file twice to check that
	// existing records are not deleted prior to import
	for i := 0; i < 2; i++ {
		if !assert.True(t, reportIfErr(dataImportCSV("from_csv", ",", "../../data/import_csv/test.csv"))) {
			return
		}
	}
	if !assert.True(t,
		reportIfErr(
			errIfQueryResultMismatch(t,
				"select id, line from from_csv order by id",
				`{"columns":["id","line"],"rows":[["1","line 1"],["1","line 1"],["2","line-2"],["2","line-2"]]}`))) {
		return
	}
}

func testDataImportCSVIncorrectData(t *testing.T) {
	err := dataImportCSV("from_csv_bad_data", ",", "../../data/import_csv/incorrect-data.csv")
	assert.NotNilf(t, err, "testDataImportCSVIncorrectData must have returned error")
	if err != nil {
		msg := err.Error()
		assert.Contains(t, msg, "record on line 2: wrong number of fields")
	}
}

func testDataImportCSVAlternativeDelimiter(t *testing.T) {
	if !assert.True(t,
		reportIfErr(dataImportCSV("from_csv_alternative_delimiter",
			";",
			"../../data/import_csv/alternative-delimiter.csv"))) {
		return
	}
	if !assert.True(t,
		reportIfErr(
			errIfQueryResultMismatch(t,
				"select id, line from from_csv_alternative_delimiter order by id",
				`{"columns":["id","line"],"rows":[["1","line"]]}`))) {
		return
	}
}

func TestAll(t *testing.T) {
	if onWindows() {
		t.Log("Unit testing on Windows platform is not supported.")
		return
	}
	initVars()
	setupCommands()
	if !assert.Truef(t,
		reportIfErr(teardownDatabase()),
		"first teardownDatabase failed") {
		return
	}
	time.Sleep(1 * time.Second)
	if !assert.Truef(t,
		reportIfErr(setupDatabase()),
		"setupDatabase failed") {
		return
	}
	defer func() {
		time.Sleep(1 * time.Second)
		assert.Truef(t,
			reportIfErr(teardownDatabase()),
			"second teardownDatabase failed")
	}()

	time.Sleep(1 * time.Second)
	setupServer()
	defer teardownServer()

	// FIXME there must be a better way to wait for server to start
	time.Sleep(1 * time.Second)
	if !assert.Truef(t,
		reportIfErr(setupClient()),
		"setupClient failed") {
		return
	}

	defer func() {
		time.Sleep(1 * time.Second)
		assert.Truef(t,
			reportIfErr(teardownClient()),
			"teardownClient failed")
	}()

	t.Run("testDataImportCSVSimple", testDataImportCSVSimple)
	t.Run("testDataImportCSVIncorrectData", testDataImportCSVIncorrectData)
	t.Run("testDataImportCSVAlternativeDelimiter", testDataImportCSVAlternativeDelimiter)
}
