package main

import (
	"log"

	"github.com/Jarimus/gator/internal/config"
)

func main() {

	// Read config file to a struct
	apiCfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	// Current app state struct
	state := State{
		Config: &apiCfg,
	}

	// Initialize commands
	cmdsMap := make(map[string]func(*State, command) error)
	commands := commands{
		cmds: cmdsMap,
	}
	commands.register("login", handlerLogin)

	// Get command
	cmd, err := getCommand()
	if err != nil {
		log.Fatal(err)
	}

	// Run the command
	err = commands.run(&state, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
