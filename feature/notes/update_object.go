package notes

import (
	"anytype-readwise/feature/bookmarks"
	"encoding/json"
	"fmt"
)

type AnytypeUpdateObjectRequest struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

func (c *AnytypeClient) CreateBookUpdateRequest(book bookmarks.ReadwiseBook, anytypeObject AnytypeGetObjectResponseItem) AnytypeUpdateObjectRequest {

	return AnytypeUpdateObjectRequest{
		Name: anytypeObject.Name + "UPDATED!",
		Body: book.Author, // FIXME There's no way to update the body of a note currently
	}

}

func (c *AnytypeClient) UpdateObject(spaceID string, objectID string, req AnytypeUpdateObjectRequest) (*AnytypeObject, error) {
	endpoint := fmt.Sprintf("/v1/spaces/%s/objects/%s", spaceID, objectID)
	resp, err := c.makeRequest("POST", endpoint, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create object: %w", err)
	}
	defer resp.Body.Close()

	var objResp AnytypeCreateObjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&objResp); err != nil {
		return nil, fmt.Errorf("failed to decode object response: %w", err)
	}

	object := objResp.Object.toAnytypeObject()

	return &object, nil
}
