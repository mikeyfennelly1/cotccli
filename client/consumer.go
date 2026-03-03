package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type NewStream struct {
	Name   string `json:"name"`
	Parent string `json:"parent"`
}

func NewConsumerClient(baseUrl string) *ConsumerClient {
	return &ConsumerClient{BaseUrl: baseUrl}
}

type ConsumerClient struct {
	BaseUrl string
}

func (client ConsumerClient) CreateStream(stream NewStream) error {
	log.Infof("creating stream: %v", stream)

	jsonData, err := json.Marshal(stream)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/consumer/streams", client.BaseUrl), bytes.NewBuffer(jsonData))
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
		return fmt.Errorf("api returned non-201 status: %d", resp.StatusCode)
	}

	return nil
}
