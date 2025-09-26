package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/MarDoA/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	xm, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	rss := RSSFeed{}
	err = xml.Unmarshal(xm, &rss)
	if err != nil {
		return nil, err
	}
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	for i, it := range rss.Channel.Item {
		it.Title = html.UnescapeString(it.Title)
		it.Description = html.UnescapeString(it.Description)
		rss.Channel.Item[i] = it
	}
	return &rss, nil
}

func scrapFeeds(s *state) error {
	feeddb, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		UpdatedAt: time.Now().UTC(), ID: feeddb.ID})
	if err != nil {
		return err
	}
	feed, err := fetchFeed(context.Background(), feeddb.Url)
	if err != nil {
		return err
	}
	for _, it := range feed.Channel.Item {
		if it.Title == "" && it.Description == "" {
			continue
		}
		date, err := pasreTime(it.PubDate)
		if err != nil {
			return err
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       sql.NullString{String: it.Title, Valid: it.Title != ""},
			Url:         it.Link,
			Description: sql.NullString{String: it.Description, Valid: it.Description != ""},
			PublishedAt: sql.NullTime{Time: date, Valid: date != time.Time{}},
			FeedID:      feeddb.ID,
		})
		if err != nil {
			if pgerror, ok := err.(*pq.Error); ok {
				if pgerror.Code == "23505" {
					continue
				}
			}
			return err
		}

	}
	return nil
}

func pasreTime(s string) (time.Time, error) {
	formats := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"02 Jan 2006 15:04:05 MST",
	}

	for _, layout := range formats {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, nil
}
