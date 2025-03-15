package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Jarimus/gator/internal/database"
	RSS "github.com/Jarimus/gator/internal/rss"
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

func handlerDisplayUsers(s *State, _ command) error {
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

func handlerAggregateRSS(s *State, _ command) error {
	rss, err := RSS.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%v", rss)

	return nil
}
