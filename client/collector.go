package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func NewCollectorClient(host string, port int) *CollectorClient {
	baseUrl := fmt.Sprintf("http://%s:%d", host, port)
	return &CollectorClient{
		BaseUrl: baseUrl,
	}
}

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

func (client CollectorClient) SendMessage(message Message, topic string) error {
	log.Infof("sending message to api: %v", message)

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", client.BaseUrl, topic), bytes.NewBuffer(jsonData))
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

	if resp.StatusCode != http.StatusCreated {
		log.Warnf("api returned non-201 status code: %d", resp.StatusCode)
	}

	return nil
}

type Message struct {
	ProducerName string
	ReadTime     int64              `json:"read_time"`
	Values       map[string]float64 `json:"values"`
}
