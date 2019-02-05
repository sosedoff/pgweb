package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Backend represents a third party configuration source
type Backend struct {
	Endpoint    string
	Token       string
	PassHeaders string
}

// BackendRequest represents a payload sent to the third-party source
type BackendRequest struct {
	Resource string            `json:"resource"`
	Token    string            `json:"token"`
	Headers  map[string]string `json:"headers"`
}

// BackendCredential represents the third-party response
type BackendCredential struct {
	DatabaseURL string `json:"database_url"`
}

// FetchCredential sends an authentication request to a third-party service
func (be Backend) FetchCredential(resource string, c *gin.Context) (*BackendCredential, error) {
	request := BackendRequest{
		Resource: resource,
		Token:    be.Token,
		Headers:  map[string]string{},
	}

	// Pass white-listed client headers to the backend request
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
		return nil, fmt.Errorf("unable to connect to the auth backend")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("got HTTP error %v from backend", resp.StatusCode)
	}

	cred := &BackendCredential{}
	if err := json.NewDecoder(resp.Body).Decode(cred); err != nil {
		return nil, err
	}
	if cred.DatabaseURL == "" {
		return nil, fmt.Errorf("database URL was not provided")
	}

	return cred, nil
}
