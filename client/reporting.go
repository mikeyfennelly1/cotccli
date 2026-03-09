package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func NewReportingClient(baseUrl string) *ReportingClient {
	return &ReportingClient{BaseUrl: baseUrl}
}

type ReportingClient struct {
	BaseUrl string
}

type StreamNode struct {
	ID       int          `json:"id"`
	Name     string       `json:"name"`
	Children []StreamNode `json:"children"`
	Sources  []string     `json:"sources"`
}

func (client ReportingClient) Health() error {
	resp, err := http.Get(fmt.Sprintf("%s/api/reporting/health", client.BaseUrl))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("reporting: unhealthy (status %d)", resp.StatusCode)
	}
	return nil
}

func (client ReportingClient) GetStreamHierarchy() error {
	resp, err := http.Get(fmt.Sprintf("%s/api/reporting/streams", client.BaseUrl))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api returned non-200 status: %d", resp.StatusCode)
	}

	var streams []StreamNode
	if err := json.NewDecoder(resp.Body).Decode(&streams); err != nil {
		return err
	}

	for _, stream := range streams {
		printStreamTree(stream, "", true)
	}

	return nil
}

type streamLookup struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func (client ReportingClient) GetStreamUUIDByName(name string) (string, error) {
	url := fmt.Sprintf("%s/api/reporting/streams?name=%s", client.BaseUrl, name)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("api returned non-200 status: %d", resp.StatusCode)
	}

	var stream streamLookup
	if err := json.NewDecoder(resp.Body).Decode(&stream); err != nil {
		return "", err
	}

	return stream.UUID, nil
}

func (client ReportingClient) SubscribeToStream(streamID string) error {
	url := fmt.Sprintf("%s/api/reporting/streams/%s", client.BaseUrl, streamID)
	log.Debugf("subscribing to stream at url: %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/event-stream")

	resp, err := (&http.Client{Timeout: 0}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api returned non-200 status: %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			fmt.Println(strings.TrimPrefix(line, "data: "))
		}
	}

	return scanner.Err()
}

func (client ReportingClient) GetProducersByStreamId(streamId string) ([]Producer, error) {
	url := fmt.Sprintf("%s/api/reporting/streams/%s/producers", client.BaseUrl, streamId)
	resp, err := http.Get(url)
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

func (client ReportingClient) GetProducerByName(name string) (*ProducerMetadata, error) {
	url := fmt.Sprintf("%s/api/reporting/producers?name=%s", client.BaseUrl, name)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errFromResponseBody(*resp)
	}

	var producer ProducerMetadata
	if err := json.NewDecoder(resp.Body).Decode(&producer); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &producer, nil
}

func printStreamTree(node StreamNode, prefix string, isRoot bool) {
	if isRoot {
		fmt.Println(node.Name)
	}
	for i, child := range node.Children {
		isLast := i == len(node.Children)-1
		connector := "├── "
		childPrefix := prefix + "│   "
		if isLast {
			connector = "└── "
			childPrefix = prefix + "    "
		}
		fmt.Printf("%s%s%s\n", prefix, connector, child.Name)
		printStreamTree(child, childPrefix, false)
	}
}
