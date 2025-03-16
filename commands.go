package main

import (
	"context"
	"errors"
	"fmt"
	"html"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/Jarimus/gator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*State, command) error
}

// registers a new handler function to the commands struct.
func (c *commands) register(name string, f func(*State, command) error) {
	c.cmds[name] = f
}

// Gets the command line arguments and return a command struct
func getCommand() (command, error) {
	if len(os.Args) < 2 {
		return command{}, errors.New("not enough arguments: need <command> [optionals]")
	}
	if len(os.Args) == 2 {
		return command{name: os.Args[1]}, nil
	}
	return command{
		name: os.Args[1],
		args: os.Args[2:],
	}, nil
}

// Runs the command of the given name
func (c *commands) run(s *State, cmd command) error {
	err := c.cmds[cmd.name](s, cmd)
	return err
}

// Logs in the user in the first argument by setting a
func handlerLogin(s *State, cmd command) error {
	// Check for args
	if len(cmd.args) == 0 {
		return errors.New("login command expects a single optional argument: <login> <username>")
	}

	// Use the first arg to set the current user
	userName := cmd.args[0]

	// Check if the user exists
	_, err := s.dbQueries.GetUserByName(context.Background(), userName)
	if err != nil {
		return errors.New("username not found")
	}

	err = s.config.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to: %s\n", userName)

	return nil
}

// Registers a new user
func handlerRegister(s *State, cmd command) error {
	// Check for args
	if len(cmd.args) == 0 {
		return errors.New("register command expects a single optional argument: <register> <username>")
	}

	// Use the first arg to get the username
	newUser := cmd.args[0]

	// First check if the current user exists
	dbUser, err := s.dbQueries.GetUserByName(context.Background(), newUser)
	if err != nil {
	} else {
		if dbUser.Name == newUser {
			return errors.New("user already exists")
		}
	}

	// parameters for the database quer
	params := database.CreateUserParams{
		ID:   uuid.New(),
		Name: newUser,
	}

	// Register the user
	_, err = s.dbQueries.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	// Login the registered user
	s.config.SetUser(newUser)

	fmt.Printf("New user registered: %s\nCurrent user set to registered user.\n", newUser)

	return nil
}

func handlerReset(s *State, cmd command) error {
	err := s.dbQueries.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Print("All users deleted.\n")
	return nil
}

func handlerListUsers(s *State, _ command) error {
	users, err := s.dbQueries.GetAllUsers(context.Background())
	if err != nil {
		return err
	}
	if users == nil {
		fmt.Print("No users.\n")
	}

	for i, user := range users {
		if user.Name == s.config.CurrentUser {
			fmt.Printf("%d: %s (current)\n", i+1, user.Name)
		} else {
			fmt.Printf("%d: %s\n", i+1, user.Name)
		}

	}
	return nil
}

func handlerAggregateRSS(s *State, cmd command) error {

	if len(cmd.args) < 1 {
		return errors.New("not enought arguments. usage: agg <time_between_requests> (valid time units: \"ns\", \"us\", \"ms\", \"s\", \"m\", \"h\")")
	}

	// parse time arguments
	interval, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}

}

// Stores a given feed (title+url) to the database, connected to the current user
func handlerAddFeed(s *State, cmd command, dbUser database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("invalid arguments. usage: addFeed <\"feed title\"> <url>")
	}

	// Unpack args
	feedTitle, url := cmd.args[0], cmd.args[1]

	// Store the feed in the database
	feedParams := database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   feedTitle,
		Url:    url,
		UserID: dbUser.ID,
	}

	dbFeed, err := s.dbQueries.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:     uuid.New(),
		FeedID: feedParams.ID,
		UserID: dbUser.ID,
	}

	s.dbQueries.CreateFeedFollow(context.Background(), feedFollowParams)

	fmt.Printf("New feed added and followed!\nID: %s\nName: %s\nurl: %s\nfor user: %s\n", dbFeed.ID, dbFeed.Name, dbFeed.Url, s.config.CurrentUser)

	return nil
}

func handlerListFeeds(s *State, _ command) error {
	feeds, err := s.dbQueries.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}
	if feeds == nil {
		fmt.Print("No feeds found.\n")
		return nil
	}

	fmt.Print("*******************\n")
	for _, feed := range feeds {
		user, err := s.dbQueries.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Name: %s\nurl: %s\nuser: %s\n*******************\n", feed.Name, feed.Url, user.Name)

	}

	return nil
}

func handlerFollowFeed(s *State, cmd command, dbCurrentUser database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("not enought arguments: follow <\"url\">")
	}

	url := cmd.args[0]

	dbFeed, err := s.dbQueries.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:     uuid.New(),
		FeedID: dbFeed.ID,
		UserID: dbCurrentUser.ID,
	}

	dbFeedFollow, err := s.dbQueries.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("New feed follow added!\nFeed name: %s\nCurrent user: %s\n", dbFeedFollow.FeedName, dbFeedFollow.UserName)

	return nil
}

func handlerUnfollowFeed(s *State, cmd command, dbUser database.User) error {

	if len(cmd.args) < 1 {
		fmt.Print("Invalid arguments. usage: unfollow <url>\n")
		return nil
	}

	url := cmd.args[0]

	dbFeed, err := s.dbQueries.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	unfollowParams := database.UnfollowFeedParams{
		FeedID: dbFeed.ID,
		UserID: dbUser.ID,
	}

	err = s.dbQueries.UnfollowFeed(context.Background(), unfollowParams)
	if err != nil {
		return err
	}

	fmt.Printf("Unfollowed %s for user %s\n", dbFeed.Name, dbUser.Name)

	return nil
}

func handlerListFeedFollows(s *State, cmd command, dbUser database.User) error {

	dbFeedFollows, err := s.dbQueries.GetFeedFollowsForUserByID(context.Background(), dbUser.ID)
	if err != nil {
		return err
	}

	if dbFeedFollows == nil {
		fmt.Print("No feeds being followed.\n")
		return nil
	}

	fmt.Printf("Feeds %s is following:\n", dbFeedFollows[0].UserName)
	for i, feedFollow := range dbFeedFollows {
		fmt.Printf("%d: %s (%s)\n", i+1, feedFollow.FeedName, feedFollow.Url)
	}

	return nil
}

func handlerBrowsePosts(s *State, cmd command, dbUser database.User) error {
	var limit int
	var err error
	if len(cmd.args) >= 1 {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return errors.New("optional argument needs to be an integer. defaults to 2. usage: browse [limit]")
		}
	} else {
		limit = 2
	}

	params := database.GetPostsForUserParams{
		UserID: dbUser.ID,
		Limit:  int32(limit),
	}

	posts, err := s.dbQueries.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}

	for _, post := range posts {
		description := html.UnescapeString(post.Description)
		pTagRemover, err := regexp.Compile(`<.*?>`)
		if err != nil {
			return err
		}
		description = pTagRemover.ReplaceAllString(description, "")
		fmt.Printf("%s (%v)\n\n%s\n\n%s\n", post.Title, post.PublishedAt.Format("02-01-2006"), post.Url, description)
		print("********************************************************\n\n")
	}

	return nil
}
