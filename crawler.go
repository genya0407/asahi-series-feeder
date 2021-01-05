package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/itchyny/timefmt-go"
	"github.com/p1ass/feeder"
	"github.com/pkg/errors"
)

type Index struct {
	AuthorImage   string `json:"author_image"`
	AuthorName    string `json:"author_name"`
	AuthorProfile string `json:"author_profile"`
	Description   string `json:"description"`
	ID            string `json:"id"`
	Image         string `json:"image"`
	IsRemoved     string `json:"is_removed"`
	Label         string `json:"label"`
	Number        string `json:"number"`
	Type          string `json:"type"`
}

type Item struct {
	Count       int    `json:"count"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Image       string `json:"image"`
	Limited     int    `json:"limited"`
	ReleaseDate string `json:"release_date"`
	Title       string `json:"title"`
}

func (it *Item) ReleasedAt() *time.Time {
	t, err := timefmt.Parse(it.ReleaseDate, "%Y%m%d%H%M%S")
	if err != nil {
		panic(fmt.Sprintf(`Unrecognized time value: %s`, it.ReleaseDate))
	}
	return &t
}

type AsahiSeriesResponse struct {
	Index Index  `json:"index"`
	Items []Item `json:"items"`
}

type AsahiSeriesCrawler struct {
	SeriesID string
}

func (crawler *AsahiSeriesCrawler) Crawl() (*feeder.Feed, error) {
	url := fmt.Sprintf("https://www.asahicom.jp/rensai/json/as%s.json", crawler.SeriesID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get response.")
	}

	var asahiSeriesResponse AsahiSeriesResponse
	err = json.NewDecoder(resp.Body).Decode(&asahiSeriesResponse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response body.")
	}

	items := []*feeder.Item{}
	for _, article := range asahiSeriesResponse.Items {
		items = append(items, convertToItem(article))
	}

	return &feeder.Feed{
		Title:       asahiSeriesResponse.Index.Label,
		Link:        &feeder.Link{Href: fmt.Sprintf("https://www.asahi.com/rensai/list.html?id=%s", asahiSeriesResponse.Index.ID)},
		Description: asahiSeriesResponse.Index.Description,
		Author: &feeder.Author{
			Name:  "",
			Email: ""},
		Created: time.Time{},
		Items:   items,
	}, nil
}

func convertToItem(a Item) *feeder.Item {
	return &feeder.Item{
		Title:       a.Title,
		Link:        &feeder.Link{Href: fmt.Sprintf("https://digital.asahi.com/articles/%s.html", a.ID)},
		Created:     a.ReleasedAt(),
		ID:          a.ID,
		Description: a.Description,
	}
}
