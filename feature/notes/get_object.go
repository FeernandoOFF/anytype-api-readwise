package notes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AnytypeGetObjectsResponse struct {
	Data []AnytypeGetObjectResponseItem `json:"data"`
}

type AnytypeGetObjectResponseItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
	Properties []struct {
		Key   string `json:"key"`
		Value string `json:"text"`
	}
}

func (c *AnytypeClient) GetObjects(spaceID string, typeKey string) ([]AnytypeGetObjectResponseItem, error) {
	// Get all the objects in the space (there's no way to filter)
	endpoint := fmt.Sprintf("/v1/spaces/%s/objects", spaceID)
	resp, err := c.makeRequest("GET", endpoint, nil)

	if err != nil {

		return nil, fmt.Errorf("failed to get objects: %w", err)

	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return nil, fmt.Errorf("anytype API returned status %d", resp.StatusCode)

	}

	var objsResp AnytypeGetObjectsResponse

	if err := json.NewDecoder(resp.Body).Decode(&objsResp); err != nil {

		return nil, fmt.Errorf("failed to decode objects response: %w", err)

	}

	// Filter objects based on type.name matching typeKey
	var filteredData []AnytypeGetObjectResponseItem
	for _, item := range objsResp.Data {
		if item.Type.Name == typeKey {
			filteredData = append(filteredData, item)
		}
	}

	if len(filteredData) == 0 {
		return nil, fmt.Errorf("no objects found with type %s", typeKey)
	}

	return filteredData, nil

}
