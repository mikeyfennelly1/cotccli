package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ServerErr struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

func (ce ServerErr) toGolangErr() error {
	return fmt.Errorf("message: %s", ce.Message)
}

func errFromResponseBody(resp http.Response) error {
	var clientErr ServerErr
	if err := json.NewDecoder(resp.Body).Decode(&clientErr); err == nil && clientErr.Timestamp != "" {
		return clientErr.toGolangErr()
	}
	return fmt.Errorf("unexpected response from server, unable to parse to ServerErr")
}
