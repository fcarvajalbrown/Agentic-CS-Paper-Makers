package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const semanticScholarBaseURL = "https://api.semanticscholar.org/graph/v1/paper/search"

type SemanticPaper struct {
	Title    string   `json:"title"`
	Authors  []string `json:"authors"`
	Year     int      `json:"year"`
	Abstract string   `json:"abstract"`
	URL      string   `json:"url"`
}

type semanticResponse struct {
	Data []semanticEntry `json:"data"`
}

type semanticEntry struct {
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Abstract string `json:"abstract"`
	PaperID  string `json:"paperId"`
	Authors  []struct {
		Name string `json:"name"`
	} `json:"authors"`
	ExternalIDs struct {
		ArXiv string `json:"ArXiv"`
	} `json:"externalIds"`
}

func SearchSemanticScholar(query string, maxResults int) ([]SemanticPaper, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("limit", fmt.Sprintf("%d", maxResults))
	params.Set("fields", "title,authors,year,abstract,externalIds")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(semanticScholarBaseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("semantic scholar request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("semantic scholar status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("semantic scholar read: %w", err)
	}

	var result semanticResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("semantic scholar parse: %w", err)
	}

	papers := make([]SemanticPaper, 0, len(result.Data))
	for _, e := range result.Data {
		paper := SemanticPaper{
			Title:    e.Title,
			Year:     e.Year,
			Abstract: e.Abstract,
		}

		for _, a := range e.Authors {
			paper.Authors = append(paper.Authors, a.Name)
		}

		if e.ExternalIDs.ArXiv != "" {
			paper.URL = "https://arxiv.org/abs/" + e.ExternalIDs.ArXiv
		} else {
			paper.URL = "https://www.semanticscholar.org/paper/" + e.PaperID
		}

		papers = append(papers, paper)
	}

	return papers, nil
}
