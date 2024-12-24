package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/alfanzain/custom-agent-allocation/responses"
)

type QiscusService struct {
	BaseURL   string
	AppID     string
	SecretKey string
}

func NewQiscusService(baseURL, appID, secretKey string) *QiscusService {
	return &QiscusService{
		BaseURL:   baseURL,
		AppID:     appID,
		SecretKey: secretKey,
	}
}

func (qs *QiscusService) AllocateAgent() (*responses.QiscusAllocateAgentResponse, error) {
	endpoint := fmt.Sprintf("%s/api/v1/admin/service/allocate_agent", qs.BaseURL)

	data := url.Values{}
	data.Set("source", "qiscus")

	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Qiscus-App-Id", qs.AppID)
	req.Header.Set("Qiscus-Secret-Key", qs.SecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to allocate agent: received status %v", resp.Status)
	}

	var response *responses.QiscusAllocateAgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	fmt.Println("Agent allocated successfully:", response)
	return response, nil
}

func (qs *QiscusService) AssignAgent(roomID string, agentID uint) (*responses.QiscusAssignAgentResponse, error) {
	endpoint := fmt.Sprintf("%s/api/v1/admin/service/assign_agent", qs.BaseURL)

	data := url.Values{}
	data.Set("room_id", roomID)
	data.Set("agent_id", fmt.Sprintf("%d", agentID))
	data.Set("max_agent", "1")

	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Qiscus-App-Id", qs.AppID)
	req.Header.Set("Qiscus-Secret-Key", qs.SecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to assign agent: received status %v", resp.Status)
	}

	var response *responses.QiscusAssignAgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	fmt.Println("Agent assigned successfully:", response)
	return response, nil
}
