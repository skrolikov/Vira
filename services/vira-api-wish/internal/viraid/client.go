package viraid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"vira-api-wish/internal/types"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL, HTTPClient: &http.Client{}}
}

// Register
func (c *Client) Register(ctx context.Context, req types.RegisterRequest) (*types.AuthResponse, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(req); err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Post(c.BaseURL+"/register", "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("vira-id register: status %d", resp.StatusCode)
	}
	var out types.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Login
func (c *Client) Login(ctx context.Context, req types.LoginRequest) (*types.AuthResponse, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(req); err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Post(c.BaseURL+"/login", "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vira-id login: status %d", resp.StatusCode)
	}
	var out types.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
