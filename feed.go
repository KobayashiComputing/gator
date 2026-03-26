package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"log"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	f := RSSFeed{}

	// create a new http request and make sure it doesn't throw an error...
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil) // (*Request, error)
	if err != nil {
		return &f, err
	}

	// set the "user-agent" in the header to our app, which is "gator"
	req.Header.Set("user-agent", "gator")

	// create http client
	client := &http.Client{}

	// execute the request
	resp, err := client.Do(req)
	if err != nil {
		return &f, err
	}
	defer resp.Body.Close()

	// get the body data...
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &f, err
	}

	if err := xml.Unmarshal([]byte(body), &f); err != nil {
		log.Fatalf("Error unmarshaling XML: %v", err)
	}

	// unescape the channel's title and description...
	f.Channel.Description = html.UnescapeString(f.Channel.Description)
	f.Channel.Title = html.UnescapeString(f.Channel.Title)

	// fmt.Printf("Title: %s\n", f.Channel.Title)
	// fmt.Printf("Desctiption: %s\n\n", f.Channel.Description)

	// unescape the title and description for each item in the channel
	for _, item := range(f.Channel.Item) {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &f, nil
}