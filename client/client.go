package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/sysinfo"
	log "github.com/sirupsen/logrus"
)

const apiEndpoint = "http://localhost:8080/sysinfo"

// PushToAggregator sends a sysinfo message to the API
func PushToAggregator(message sysinfo.Message) error {
	log.Infof("sending message to api: %v", message)

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warnf("api returned non-200 status code: %d", resp.StatusCode)
	}

	return nil
}
