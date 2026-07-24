// Package client is a minimal client for the ServiceNow Table API
// (/api/now/table/{table}), authenticating with HTTP basic auth.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client talks to a single ServiceNow instance's Table API over HTTP basic auth.
type Client struct {
	BaseURL    string // e.g. "https://dev12345.service-now.com"
	Username   string
	Password   string
	HTTPClient *http.Client
}

// New builds a Client for the given instance host (e.g. "dev12345.service-now.com")
// or full base URL, authenticating with basic auth.
func New(instance, username, password string) *Client {
	base := instance
	if !strings.Contains(base, "://") {
		base = "https://" + base
	}
	return &Client{
		BaseURL:    strings.TrimRight(base, "/"),
		Username:   username,
		Password:   password,
		HTTPClient: &http.Client{},
	}
}

type envelope struct {
	Result json.RawMessage `json:"result"`
	Error  *apiError       `json:"error"`
}

type apiError struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (c *Client) tableURL(table, sysID string, params url.Values) string {
	u := fmt.Sprintf("%s/api/now/table/%s", c.BaseURL, table)
	if sysID != "" {
		u += "/" + sysID
	}
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	return u
}

func (c *Client) request(method, requestURL string, body []byte) (json.RawMessage, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, requestURL, reader)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("authentication failed (401)")
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("servicenow API error (%d): %s", resp.StatusCode, string(respBody))
	}
	if len(respBody) == 0 {
		return nil, nil
	}

	var env envelope
	if err := json.Unmarshal(respBody, &env); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	if env.Error != nil {
		return nil, fmt.Errorf("servicenow error: %s: %s", env.Error.Message, env.Error.Detail)
	}
	return env.Result, nil
}

// List retrieves records from a table, optionally filtered/sorted via sysparm_* query params.
func (c *Client) List(table string, params url.Values) ([]map[string]any, error) {
	raw, err := c.request(http.MethodGet, c.tableURL(table, "", params), nil)
	if err != nil {
		return nil, err
	}
	var records []map[string]any
	if err := json.Unmarshal(raw, &records); err != nil {
		return nil, fmt.Errorf("decoding records: %w", err)
	}
	return records, nil
}

// Get retrieves a single record by sys_id.
func (c *Client) Get(table, sysID string, params url.Values) (map[string]any, error) {
	raw, err := c.request(http.MethodGet, c.tableURL(table, sysID, params), nil)
	if err != nil {
		return nil, err
	}
	var record map[string]any
	if err := json.Unmarshal(raw, &record); err != nil {
		return nil, fmt.Errorf("decoding record: %w", err)
	}
	return record, nil
}

// Create inserts a new record into a table.
func (c *Client) Create(table string, fields map[string]any) (map[string]any, error) {
	body, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}
	raw, err := c.request(http.MethodPost, c.tableURL(table, "", nil), body)
	if err != nil {
		return nil, err
	}
	var record map[string]any
	if err := json.Unmarshal(raw, &record); err != nil {
		return nil, fmt.Errorf("decoding record: %w", err)
	}
	return record, nil
}

// Update patches an existing record identified by sys_id.
func (c *Client) Update(table, sysID string, fields map[string]any) (map[string]any, error) {
	body, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}
	raw, err := c.request(http.MethodPatch, c.tableURL(table, sysID, nil), body)
	if err != nil {
		return nil, err
	}
	var record map[string]any
	if err := json.Unmarshal(raw, &record); err != nil {
		return nil, fmt.Errorf("decoding record: %w", err)
	}
	return record, nil
}

// Delete removes a record identified by sys_id.
func (c *Client) Delete(table, sysID string) error {
	_, err := c.request(http.MethodDelete, c.tableURL(table, sysID, nil), nil)
	return err
}

// Ping performs a lightweight request to verify the credentials/instance are valid.
func (c *Client) Ping() error {
	params := url.Values{"sysparm_limit": []string{"1"}}
	_, err := c.List("sys_user", params)
	return err
}
