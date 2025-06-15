package notes

import (
	"anytype-readwise/core"
	"anytype-readwise/feature/bookmarks"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type AnytypeClient struct {
	apiKey     string
	baseURL    string
	version    string
	httpClient *http.Client
	config     *core.Config
}

type AnytypeObject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewAnytypeClient(apiKey, baseURL, version string, config *core.Config) *AnytypeClient {
	return &AnytypeClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		version: version,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		config: config,
	}
}

func (c *AnytypeClient) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&reqBody).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, &reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Anytype-Version", c.version)
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

func (c *AnytypeClient) CreateOrUpdateNoteFromBook(spaceID string, book bookmarks.ReadwiseBook, content string) (*AnytypeObject, error) {
	objects, err := c.GetObjects(spaceID, c.config.ObjectType)
	if err != nil {
		return nil, fmt.Errorf("failed to get objects for space %s: %w", spaceID, err)
	}

	for _, obj := range objects {
		// Check properties (description, to see if it matches BookmarkID)
		for _, prop := range obj.Properties {
			if prop.Key == "description" {
				if prop.Value == strconv.Itoa(book.ID) {
					fmt.Println("Found a matching note!:", book.Title, prop.Value)
					req := c.CreateBookUpdateRequest(book, obj)
					updatedObject, err := c.UpdateObject(spaceID, obj.ID, req)
					if err != nil {
						fmt.Println("Failed to update object:", err)
						return nil, nil // Currently fails to update
					}
					return updatedObject, nil
				}

			}

		}
	}

	req := c.CreateBookObjectRequest(book, content)
	createdObject, err := c.CreateObject(spaceID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create object: %w", err)
	}
	object := createdObject.toAnytypeObject()
	return &object, err
}
