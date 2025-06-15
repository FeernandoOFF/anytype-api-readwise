package notes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AnytypeSpace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AnytypeSpacesResponse struct {
	Data []AnytypeSpace `json:"data"`
}

// GetSpaceID returns the space ID from config or the first one in the list
func (c *AnytypeClient) GetSpaceID() (string, error) {
	if c.config.SpaceID != "" {
		return c.config.SpaceID, nil
	}

	spaces, err := c.GetSpaces()
	if err != nil {
		return "", err
	}

	if len(spaces) == 0 {
		return "", fmt.Errorf("no spaces found")
	}

	// Use the first space if none specified
	defaultSpace := spaces[0]
	fmt.Printf("Using space: %s (ID: %s)\n", defaultSpace.Name, defaultSpace.ID)
	return defaultSpace.ID, nil
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
