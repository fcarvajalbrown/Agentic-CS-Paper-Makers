package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const zenodoBaseURL = "https://zenodo.org/api/records"

type ZenadoPaper struct {
	Title    string   `json:"title"`
	Authors  []string `json:"authors"`
	Year     int      `json:"year"`
	Abstract string   `json:"abstract"`
	URL      string   `json:"url"`
}

type zenodoResponse struct {
	Hits struct {
		Hits []zenodoEntry `json:"hits"`
	} `json:"hits"`
}

type zenodoEntry struct {
	ID       int    `json:"id"`
	DOI      string `json:"doi"`
	Metadata struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		PublicationDate string `json:"publication_date"`
		Creators []struct {
			Name string `json:"name"`
		} `json:"creators"`
	} `json:"metadata"`
	Links struct {
		HTML string `json:"html"`
	} `json:"links"`
}

func SearchZenodo(query string, maxResults int) ([]ZenadoPaper, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("size", fmt.Sprintf("%d", maxResults))
	params.Set("type", "publication")
	params.Set("sort", "bestmatch")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(zenodoBaseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("zenodo request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("zenodo status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("zenodo read: %w", err)
	}

	var result zenodoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("zenodo parse: %w", err)
	}

	papers := make([]ZenadoPaper, 0, len(result.Hits.Hits))
	for _, e := range result.Hits.Hits {
		paper := ZenadoPaper{
			Title:    e.Metadata.Title,
			Abstract: strings.TrimSpace(e.Metadata.Description),
			URL:      e.Links.HTML,
		}

		for _, c := range e.Metadata.Creators {
			paper.Authors = append(paper.Authors, c.Name)
		}

		if len(e.Metadata.PublicationDate) >= 4 {
			fmt.Sscanf(e.Metadata.PublicationDate[:4], "%d", &paper.Year)
		}

		papers = append(papers, paper)
	}

	return papers, nil
}
