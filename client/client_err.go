package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ServerErr struct {
	Error string `json:"error"`
}

func (ce ServerErr) toGolangErr() error {
	return fmt.Errorf("ServerErr: %s", ce.Error)
}

func errFromResponseBody(resp http.Response) error {
	var clientErr ServerErr
	if err := json.NewDecoder(resp.Body).Decode(&clientErr); err == nil && clientErr.Error != "" {
		return clientErr.toGolangErr()
	}
	return fmt.Errorf("unexpected response from server, unable to parse to ServerErr")
}
