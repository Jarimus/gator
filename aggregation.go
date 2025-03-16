package main

import (
	"context"
	"fmt"

	"github.com/Jarimus/gator/internal/database"
	rss "github.com/Jarimus/gator/internal/rss"
	"github.com/araddon/dateparse"
	"github.com/google/uuid"
)

func scrapeFeeds(s *State) error {
	nextFeed, err := s.dbQueries.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = s.dbQueries.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return err
	}

	rssFeed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	fmt.Print("*******************************************\n")
	fmt.Printf("Fetching RSS feed from %s\n", rssFeed.Channel.Title)
	fmt.Print("*******************************************\n")
	for _, item := range rssFeed.Channel.Item {

		pubDate, err := dateparse.ParseAny(item.PubDate)
		if err != nil {
			return err
		}

		params := database.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pubDate,
			FeedID:      nextFeed.ID,
		}

		_, err = s.dbQueries.CreatePost(context.Background(), params)
		if err != nil {
			return err
		}
	}

	return nil
}
