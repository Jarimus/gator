package main

import (
	"github.com/Jarimus/gator/internal/config"
	"github.com/Jarimus/gator/internal/database"
)

type State struct {
	dbQueries *database.Queries
	config    *config.Config
}
