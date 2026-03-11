package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Group struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Children  []Group    `json:"children"`
	Producers []Producer `json:"producers"`
}

func (client GroupControllerClient) GetGroupHierarchy() error {
	resp, err := http.Get(fmt.Sprintf("%s/api/group", client.BaseUrl))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api returned non-200 status: %d", resp.StatusCode)
	}

	var streams []Group
	if err := json.NewDecoder(resp.Body).Decode(&streams); err != nil {
		return err
	}

	for _, stream := range streams {
		printGroupTree(stream, "", true)
	}

	return nil
}

type groupLookup struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func (client GroupControllerClient) GetGroupUUIDByName(name string) (string, error) {
	url := fmt.Sprintf("%s/api/group?name=%s", client.BaseUrl, name)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("api returned non-200 status: %d", resp.StatusCode)
	}

	var stream groupLookup
	if err := json.NewDecoder(resp.Body).Decode(&stream); err != nil {
		return "", err
	}

	return stream.UUID, nil
}

func (client GroupControllerClient) SubscribeToGroupEvents(groupId string) error {
	url := fmt.Sprintf("%s/api/group/events?group=%s", client.BaseUrl, groupId)
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

func (client ProducerClient) GetProducersByGroupId(streamId string) ([]Producer, error) {
	url := fmt.Sprintf("%s/api/group/streams/%s/producers", client.BaseUrl, streamId)
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

func (client ProducerClient) GetProducerByName(name string) (*ProducerMetadata, error) {
	url := fmt.Sprintf("%s/api/producer?name=%s", client.BaseUrl, name)
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

func printGroupTree(node Group, prefix string, isRoot bool) {
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
		printGroupTree(child, childPrefix, false)
	}
}

type NewGroup struct {
	Name   string `json:"name"`
	Parent string `json:"parent"`
}

func NewGroupControllerClient(baseUrl string) *GroupControllerClient {
	return &GroupControllerClient{BaseUrl: baseUrl}
}

type GroupControllerClient struct {
	BaseUrl string
}

func (client GroupControllerClient) Health() error {
	resp, err := http.Get(fmt.Sprintf("%s/api/group/health", client.BaseUrl))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("group: unhealthy (status %d)", resp.StatusCode)
	}
	return nil
}

func (client GroupControllerClient) CreateGroup(group string) error {
	log.Infof("creating group: %s", group)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/group?name=%s", client.BaseUrl, group), nil)
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
		return errFromResponseBody(*resp)
	}

	return nil
}

func (client GroupControllerClient) DeleteGroup(name string) error {
	url := fmt.Sprintf("%s/api/group/streams?name=%s", client.BaseUrl, name)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("api returned non-200 status: %d", resp.StatusCode)
	}

	return nil
}
