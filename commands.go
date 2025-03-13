package main

import (
	"errors"
	"fmt"
	"os"
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
	currentUser := cmd.args[0]
	err := s.Config.SetUser(currentUser)
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to: %s\n", currentUser)

	return nil
}
