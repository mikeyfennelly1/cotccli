package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type NewProducer struct {
	Name       string `json:"name"`
	StreamName string `json:"streamName"`
}

type ProducerMetadata struct {
	UUID         string `json:"uuid"`
	ProducerName string `json:"producerName"`
	StreamName   string `json:"streamName"`
}

func NewProducerClient(baseUrl string) {

}

func (client StreamControllerClient) CreateProducer(producer NewProducer) (*ProducerMetadata, error) {
	log.Infof("creating producer: %v", producer)

	jsonData, err := json.Marshal(producer)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/consumer/producers", client.BaseUrl), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, errFromResponseBody(*resp)
	}

	var created ProducerMetadata
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &created, nil
}

type Producer struct {
	UUID         string `json:"uuid"`
	ProducerName string `json:"producerName"`
	StreamId     string `json:"streamId"`
}

func (client ReportingClient) GetProducers() ([]Producer, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/reporting/producers", client.BaseUrl))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errFromResponseBody(*resp)
	}

	var producers []Producer
	if err := json.NewDecoder(resp.Body).Decode(&producers); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return producers, nil
}
