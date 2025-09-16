package connect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type Backend struct {
	Endpoint    string
	Token       string
	PassHeaders []string

	logger *logrus.Logger
}

func NewBackend(endpoint string, token string) Backend {
	return Backend{
		Endpoint: endpoint,
		Token:    token,
		logger:   logrus.StandardLogger(),
	}
}

func (be *Backend) SetLogger(logger *logrus.Logger) {
	be.logger = logger
}

func (be *Backend) SetPassHeaders(headers []string) {
	be.PassHeaders = headers
}

func (be *Backend) FetchCredential(ctx context.Context, resource string, headers http.Header) (*Credential, error) {
	be.logger.WithField("resource", resource).Debug("fetching database credential")

	request := Request{
		Resource: resource,
		Token:    be.Token,
		Headers:  map[string]string{},
	}

	// Pass allow-listed client headers to the backend request
	for _, name := range be.PassHeaders {
		request.Headers[strings.ToLower(name)] = headers.Get(name)
	}

	body, err := json.Marshal(request)
	if err != nil {
		be.logger.WithField("resource", resource).Error("backend request serialization error:", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, be.Endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		be.logger.WithField("resource", resource).Error("backend credential fetch failed:", err)
		return nil, errBackendConnectError
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("backend credential fetch received HTTP status code %v", resp.StatusCode)

		be.logger.
			WithField("resource", request.Resource).
			WithField("status", resp.StatusCode).
			Error(err)

		return nil, err
	}

	cred := &Credential{}
	if err := json.NewDecoder(resp.Body).Decode(cred); err != nil {
		return nil, err
	}

	if cred.DatabaseURL == "" {
		return nil, errConnStringRequired
	}

	return cred, nil
}
