package bus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	SessionURL string
	HTTPClient *http.Client
}

func (c *Client) Send(method, endpoint string, body interface{}, result interface{}) error {
	requestBody, err := bodyToJSON(body)
	if err != nil {
		return err
	}

	requestURL := strings.TrimSuffix(c.SessionURL+"/"+endpoint, "/")
	responseBody, err := c.makeRequest(requestURL, method, requestBody)
	if err != nil {
		return err
	}

	if result != nil {
		bodyValue := struct{ Value interface{} }{result}
		if err := json.Unmarshal(responseBody, &bodyValue); err != nil {
			return fmt.Errorf("unexpected response: %s", responseBody)
		}
	}

	return nil
}

func bodyToJSON(body interface{}) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("invalid request body: %s", err)
	}
	return bodyJSON, nil
}

func (c *Client) makeRequest(url, method string, body []byte) ([]byte, error) {
	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("invalid request: %s", err)
	}

	if body != nil {
		request.Header.Add("Content-Type", "application/json")
	}

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, parseResponseError(responseBody)
	}

	return responseBody, nil
}

func parseResponseError(body []byte) error {
	var errBody struct{ Value struct{ Message string } }
	if err := json.Unmarshal(body, &errBody); err != nil {
		return fmt.Errorf("request unsuccessful: %s", body)
	}

	var errMessage struct{ ErrorMessage string }
	if err := json.Unmarshal([]byte(errBody.Value.Message), &errMessage); err != nil {
		return fmt.Errorf("request unsuccessful: %s", errBody.Value.Message)
	}

	return fmt.Errorf("request unsuccessful: %s", errMessage.ErrorMessage)
}
