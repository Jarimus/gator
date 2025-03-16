package main

import (
	"context"
	"fmt"

	"github.com/Jarimus/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *State, cmd command, user database.User) error) func(*State, command) error {
	return func(s *State, cmd command) error {
		user, err := s.dbQueries.GetUserByName(context.Background(), s.config.CurrentUser)
		if err != nil {
			return fmt.Errorf("error getting current user: %s", err)
		}

		return handler(s, cmd, user)
	}
}
