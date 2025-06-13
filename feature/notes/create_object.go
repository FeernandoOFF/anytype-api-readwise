package notes

import (
	"anytype-readwise/feature/bookmarks"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type AnytypeCreateObjectResponseItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AnytypeCreateObjectResponse struct {
	Object AnytypeCreateObjectResponseItem `json:"object"`
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

type CreateObjectProperty struct {
	Key   string `json:"key"`
	Value string `json:"text"`
}

type CreateObjectRequest struct {
	Name       string                 `json:"name"`
	TypeKey    string                 `json:"type_key"`
	Body       string                 `json:"body"`
	Icon       *ObjectIcon            `json:"icon,omitempty"`
	Properties []CreateObjectProperty `json:"properties,omitempty"`
}

type ObjectIcon struct {
	Emoji  string `json:"emoji"`
	Format string `json:"format"`
}

// CreateBookObjectRequest creates a CreateObjectRequest for a book
func (c *AnytypeClient) CreateBookObjectRequest(book bookmarks.ReadwiseBook, content string) CreateObjectRequest {
	return CreateObjectRequest{
		Name:    fmt.Sprintf("%s - %s [SYNC]", book.Title, book.Author),
		TypeKey: strings.ToLower(c.config.ObjectType), // FIXME does this only happens for the bookmarks??
		Body:    content,
		Icon: &ObjectIcon{
			Emoji:  "📚",
			Format: "emoji",
		},
		Properties: []CreateObjectProperty{
			{
				Key:   "description",
				Value: strconv.Itoa(book.ID),
			},
		},
	}
}

func (c *AnytypeClient) CreateObject(spaceID string, req CreateObjectRequest) (*AnytypeCreateObjectResponseItem, error) {
	endpoint := fmt.Sprintf("/v1/spaces/%s/objects", spaceID)
	resp, err := c.makeRequest("POST", endpoint, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create object: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("anytype API returned status %d", resp.StatusCode)
	}

	var objResp AnytypeCreateObjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&objResp); err != nil {
		return nil, fmt.Errorf("failed to decode object response: %w", err)
	}

	return &objResp.Object, nil
}

func (respObj *AnytypeCreateObjectResponseItem) toAnytypeObject() AnytypeObject {
	return AnytypeObject{
		ID:   respObj.ID,
		Name: respObj.Name,
	}
}
