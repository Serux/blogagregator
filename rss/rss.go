package rss

import (
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

	r, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "gator")
	cl := http.Client{}
	resp, err := cl.Do(r)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(data))
	RSSFeed := RSSFeed{}
	err = xml.Unmarshal(data, &RSSFeed)
	if err != nil {
		return nil, err
	}
	RSSFeed.Channel.Description = html.UnescapeString(RSSFeed.Channel.Description)
	RSSFeed.Channel.Title = html.UnescapeString(RSSFeed.Channel.Title)

	for i := range RSSFeed.Channel.Item {
		RSSFeed.Channel.Item[i].Description = html.UnescapeString(RSSFeed.Channel.Item[i].Description)
		RSSFeed.Channel.Item[i].Title = html.UnescapeString(RSSFeed.Channel.Item[i].Title)

	}
	return &RSSFeed, nil

}
