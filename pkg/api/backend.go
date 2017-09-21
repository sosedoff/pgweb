package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Backend struct {
	Endpoint    string
	Token       string
	PassHeaders string
}

type BackendRequest struct {
	Resource string            `json:"resource"`
	Token    string            `json:"token"`
	Headers  map[string]string `json:"headers"`
}

type BackendCredential struct {
	DatabaseUrl string `json:"database_url"`
}

func (be Backend) FetchCredential(resource string, c *gin.Context) (*BackendCredential, error) {
	request := BackendRequest{
		Resource: resource,
		Token:    be.Token,
		Headers:  map[string]string{},
	}

	for _, name := range strings.Split(be.PassHeaders, ",") {
		request.Headers[strings.ToLower(name)] = c.Request.Header.Get(name)
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(be.Endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		// Any connection-related issues will show up in the server log
		log.Println("Unable to fetch backend credential:", err)

		// We dont want to expose the url of the backend here, so reply with generic error
		return nil, fmt.Errorf("Unable to connect to the auth backend")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Got HTTP error %v from backend", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cred := &BackendCredential{}
	if err := json.Unmarshal(respBody, cred); err != nil {
		return nil, err
	}
	if cred.DatabaseUrl == "" {
		return nil, fmt.Errorf("Database url was not provided")
	}

	return cred, nil
}
