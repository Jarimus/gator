package rss

import (
	"bytes"
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, bytes.NewBuffer([]byte("stuff")))
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "gator")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var respData RSSFeed

	err = xml.Unmarshal(body, &respData)

	// Unescape html to decode escaped HTML entities
	respData.Channel.Title = html.UnescapeString(respData.Channel.Title)
	respData.Channel.Description = html.UnescapeString(respData.Channel.Description)
	for i, item := range respData.Channel.Item {
		respData.Channel.Item[i].Title = html.UnescapeString(item.Title)
		respData.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &respData, err
}
