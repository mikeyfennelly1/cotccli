package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type CollectorClient struct {
	BaseUrl string
}

func (client CollectorClient) Health() error {
	resp, err := http.Get(fmt.Sprintf("%s/api/collector/health", client.BaseUrl))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("collector: unhealthy (status %d)", resp.StatusCode)
	}
	return nil
}

func (client CollectorClient) SendMessage(message Message) error {
	topic := "TIMESERIES"

	log.Infof("sending message to topic %s", topic)

	jsonData, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		return err
	}
	log.Debugf("message json:\n%s", string(jsonData))

	endpoint := fmt.Sprintf("%s/api/collector/%s", client.BaseUrl, topic)
	log.Debugf("sending POST request to endpoint: %s", endpoint)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Debugf("response status: %d, body: %s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusCreated {
		log.Warnf("api returned non-201 status code: %d", resp.StatusCode)
	}

	return nil
}

type Message struct {
	ProducerId   string             `json:"producer_id"`
	ProducerName string             `json:"producer_name"`
	ReadTime     int64              `json:"read_time"`
	Values       map[string]float64 `json:"values"`
}
