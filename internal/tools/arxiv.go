package tools

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const arxivBaseURL = "https://export.arxiv.org/api/query"

type ArxivPaper struct {
	Title    string   `json:"title"`
	Authors  []string `json:"authors"`
	Year     int      `json:"year"`
	Abstract string   `json:"abstract"`
	URL      string   `json:"url"`
}

type arxivFeed struct {
	Entries []arxivEntry `xml:"entry"`
}

type arxivEntry struct {
	Title    string        `xml:"title"`
	Summary  string        `xml:"summary"`
	Links    []arxivLink   `xml:"link"`
	Authors  []arxivAuthor `xml:"author"`
	Published string       `xml:"published"`
}

type arxivLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

type arxivAuthor struct {
	Name string `xml:"name"`
}

func SearchArxiv(query string, maxResults int) ([]ArxivPaper, error) {
	params := url.Values{}
	params.Set("search_query", "all:"+query)
	params.Set("max_results", fmt.Sprintf("%d", maxResults))
	params.Set("sortBy", "relevance")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(arxivBaseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("arxiv request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("arxiv read: %w", err)
	}

	var feed arxivFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("arxiv parse: %w", err)
	}

	papers := make([]ArxivPaper, 0, len(feed.Entries))
	for _, e := range feed.Entries {
		paper := ArxivPaper{
			Title:    strings.TrimSpace(e.Title),
			Abstract: strings.TrimSpace(e.Summary),
		}

		for _, a := range e.Authors {
			paper.Authors = append(paper.Authors, a.Name)
		}

		if len(e.Published) >= 4 {
			fmt.Sscanf(e.Published[:4], "%d", &paper.Year)
		}

		for _, l := range e.Links {
			if l.Rel == "alternate" {
				paper.URL = l.Href
				break
			}
		}

		papers = append(papers, paper)
	}

	return papers, nil
}
