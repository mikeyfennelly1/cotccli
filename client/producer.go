package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type NewProducer struct {
	Name  string `json:"name"`
	Group string `json:"group"`
}

type ProducerMetadata struct {
	ProducerName string `json:"producerName"`
	UUID         string
}

type ProducerClient struct {
	BaseUrl string
}

func NewProducerClient(baseUrl string) ProducerClient {
	return ProducerClient{baseUrl}
}

func (client ProducerClient) Health() error {
	resp, err := http.Get(fmt.Sprintf("%s/api/producer/health", client.BaseUrl))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("collector: unhealthy (status %d)", resp.StatusCode)
	}
	return nil
}

func (client ProducerClient) CreateProducer(producer NewProducer) (*ProducerMetadata, error) {
	log.Infof("creating producer: %v", producer)

	jsonData, err := json.Marshal(producer)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/producer", client.BaseUrl), bytes.NewBuffer(jsonData))
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
	GroupId      string `json:"groupName"`
}

func (client ProducerClient) GetProducersForGroup(groupName string) ([]Producer, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/producer/%s", client.BaseUrl, groupName))
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

func (client ProducerClient) GetProducers() ([]Producer, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/producer", client.BaseUrl))
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
