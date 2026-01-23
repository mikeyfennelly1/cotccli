package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/sysinfo"
	log "github.com/sirupsen/logrus"
)

// PushToAggregator sends a sysinfo message to the API
func PushToAggregator(message sysinfo.Message, host string, port int) error {
	log.Infof("sending message to api: %v", message)

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	apiEndpoint := formatAPIEndpoint(host, port)

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

func formatAPIEndpoint(host string, port int) string {
	return fmt.Sprintf("http://%s:%d/sysinfo", host, port)
}
