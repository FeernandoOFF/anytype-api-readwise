// anytype.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AnytypeClient struct {
	apiKey     string
	baseURL    string
	version    string
	httpClient *http.Client
}

type AnytypeSpace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AnytypeSpacesResponse struct {
	Data []AnytypeSpace `json:"data"`
}

type AnytypeObject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AnytypeObjectResponse struct {
	Object AnytypeObject `json:"object"`
}

type PropertyValue interface {
	Type() string
}

type TextValue struct {
	Value string `json:"text"`
}

func (v TextValue) Type() string { return "text" }

type URLValue struct {
	Value string `json:"url"`
}

func (v URLValue) Type() string { return "url" }

type NumberValue struct {
	Value int `json:"number"`
}

func (v NumberValue) Type() string { return "number" }

type SelectValue struct {
	Value string `json:"select"`
}

func (v SelectValue) Type() string { return "select" }

type Property struct {
	Key   string `json:"key"`
	Value string `json:"text"`
}

type CreateObjectRequest struct {
	Name       string      `json:"name"`
	TypeKey    string      `json:"type_key"`
	Body       string      `json:"body"`
	Icon       *ObjectIcon `json:"icon,omitempty"`
	Properties []Property  `json:"properties,omitempty"`
}

type ObjectIcon struct {
	Emoji  string `json:"emoji"`
	Format string `json:"format"`
}

func NewAnytypeClient(apiKey, baseURL, version string) *AnytypeClient {
	return &AnytypeClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		version: version,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *AnytypeClient) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&reqBody).Encode(body); err != nil {
			return nil, err
		}
	}

	fmt.Println(reqBody.String())

	req, err := http.NewRequest(method, c.baseURL+endpoint, &reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Anytype-Version", c.version)
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

func (c *AnytypeClient) GetSpaces() ([]AnytypeSpace, error) {
	resp, err := c.makeRequest("GET", "/v1/spaces", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch spaces: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("anytype API returned status %d", resp.StatusCode)
	}

	var spacesResp AnytypeSpacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&spacesResp); err != nil {
		return nil, fmt.Errorf("failed to decode spaces response: %w", err)
	}

	return spacesResp.Data, nil
}

func (c *AnytypeClient) CreateObject(spaceID string, req CreateObjectRequest) (*AnytypeObject, error) {
	endpoint := fmt.Sprintf("/v1/spaces/%s/objects", spaceID)
	resp, err := c.makeRequest("POST", endpoint, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create object: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("anytype API returned status %d", resp.StatusCode)
	}

	var objResp AnytypeObjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&objResp); err != nil {
		return nil, fmt.Errorf("failed to decode object response: %w", err)
	}

	return &objResp.Object, nil
}
