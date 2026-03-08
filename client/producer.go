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

type CreatedProducer struct {
	UUID         string `json:"uuid"`
	ProducerName string `json:"producerName"`
}

func (client ConsumerClient) CreateProducer(producer NewProducer) (*CreatedProducer, error) {
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

	if resp.StatusCode == http.StatusConflict {
		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Error != "" {
			log.Infof("producer already exists, skipping creation: %s", errResp.Error)
			return nil, fmt.Errorf("producer already exists")
		} else {
			log.Infof("producer already exists, skipping creation")
			return nil, fmt.Errorf("producer already exists")
		}
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("api returned non-201 status: %d", resp.StatusCode)
	}

	var created CreatedProducer
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &created, nil
}
