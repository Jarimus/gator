package main

import (
	"context"
	"fmt"

	rss "github.com/Jarimus/gator/internal/rss"
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
	fmt.Printf("RSS feed from %s\n", rssFeed.Channel.Title)
	fmt.Print("*******************************************\n")
	for i, item := range rssFeed.Channel.Item {
		fmt.Printf("%d: %s (%s)\n", i+1, item.Title, item.Link)
	}

	return nil
}
