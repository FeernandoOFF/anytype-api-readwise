// readwise.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ReadwiseClient implements the BookmarksProvider interface
type ReadwiseClient struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

type ReadwiseBook struct {
	ID            int         `json:"id"`
	Title         string      `json:"title"`
	Author        string      `json:"author"`
	Category      string      `json:"category"`
	Source        string      `json:"source"`
	NumHighlights int         `json:"num_highlights"`
	LastHighlight time.Time   `json:"last_highlight_at"`
	Updated       time.Time   `json:"updated"`
	CoverImageURL string      `json:"cover_image_url"`
	Highlights    []Highlight `json:"highlights,omitempty"`
}

type Highlight struct {
	ID            int       `json:"id"`
	Text          string    `json:"text"`
	Note          string    `json:"note"`
	Location      int       `json:"location"`
	LocationType  string    `json:"location_type"`
	HighlightedAt time.Time `json:"highlighted_at"`
	URL           string    `json:"url"`
	Color         string    `json:"color"`
	Updated       time.Time `json:"updated"`
}

type ReadwiseBooksResponse struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []ReadwiseBook `json:"results"`
}

type ReadwiseHighlightsResponse struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous string      `json:"previous"`
	Results  []Highlight `json:"results"`
}

func NewReadwiseClient(token string) *ReadwiseClient {
	return &ReadwiseClient{
		token:   token,
		baseURL: "https://readwise.io/api/v2",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *ReadwiseClient) makeRequest(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Token "+c.token)
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

func (c *ReadwiseClient) GetBooks() ([]ReadwiseBook, error) {
	var allBooks []ReadwiseBook
	url := "/books/"

	for url != "" {
		resp, err := c.makeRequest(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch books: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("readwise API returned status %d", resp.StatusCode)
		}

		var booksResp ReadwiseBooksResponse
		if err := json.NewDecoder(resp.Body).Decode(&booksResp); err != nil {
			return nil, fmt.Errorf("failed to decode books response: %w", err)
		}

		allBooks = append(allBooks, booksResp.Results...)

		if booksResp.Next != "" {
			url = booksResp.Next[len(c.baseURL):]
		} else {
			url = ""
		}
	}

	return allBooks, nil
}

func (c *ReadwiseClient) GetHighlights(bookID int) ([]Highlight, error) {
	var allHighlights []Highlight
	url := fmt.Sprintf("/highlights/?book_id=%d", bookID)

	for url != "" {
		resp, err := c.makeRequest(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch highlights: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("readwise API returned status %d", resp.StatusCode)
		}

		var highlightsResp ReadwiseHighlightsResponse
		if err := json.NewDecoder(resp.Body).Decode(&highlightsResp); err != nil {
			return nil, fmt.Errorf("failed to decode highlights response: %w", err)
		}

		allHighlights = append(allHighlights, highlightsResp.Results...)

		if highlightsResp.Next != "" {
			url = highlightsResp.Next[len(c.baseURL):]
		} else {
			url = ""
		}
	}

	return allHighlights, nil
}
