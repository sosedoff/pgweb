package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Backend represents a third party configuration source
type Backend struct {
	Endpoint    string
	Token       string
	PassHeaders []string
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
func (be Backend) FetchCredential(ctx context.Context, resource string, c *gin.Context) (*BackendCredential, error) {
	logger.WithField("resource", resource).Debug("fetching database credential")

	request := BackendRequest{
		Resource: resource,
		Token:    be.Token,
		Headers:  map[string]string{},
	}

	// Pass white-listed client headers to the backend request
	for _, name := range be.PassHeaders {
		request.Headers[strings.ToLower(name)] = c.Request.Header.Get(name)
	}

	body, err := json.Marshal(request)
	if err != nil {
		logger.WithField("resource", resource).Error("backend request serialization error:", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, be.Endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.WithField("resource", resource).Error("backend credential fetch failed:", err)
		return nil, errBackendConnectError
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("backend credential fetch received HTTP status code %v", resp.StatusCode)

		logger.
			WithField("resource", resource).
			WithField("status", resp.StatusCode).
			Error(err)

		return nil, err
	}

	cred := &BackendCredential{}
	if err := json.NewDecoder(resp.Body).Decode(cred); err != nil {
		return nil, err
	}
	if cred.DatabaseURL == "" {
		return nil, errConnStringRequired
	}

	return cred, nil
}
