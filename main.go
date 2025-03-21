package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Jarimus/gator/internal/config"
	"github.com/Jarimus/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	// Read config file to a struct
	apiCfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", apiCfg.DbUrl)
	if err != nil {
		log.Fatalf("error opening connection to database: %s", err)
	}
	dbQueries := database.New(db)

	// Current app state struct
	state := State{
		config:    &apiCfg,
		dbQueries: dbQueries,
	}

	// Initialize commands
	cmdsMap := make(map[string]func(*State, command) error)
	commands := commands{
		cmds: cmdsMap,
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerListUsers)
	commands.register("agg", middlewareLoggedIn(handlerAggregateRSS))
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerListFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollowFeed))
	commands.register("following", middlewareLoggedIn(handlerListFeedFollows))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))
	commands.register("browse", middlewareLoggedIn(handlerBrowsePosts))

	// If program started with no args, display help
	if len(os.Args) == 1 || (len(os.Args) >= 2 && os.Args[1] == "help") {
		help(commands)
		return
	}

	// Get command from arguments
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

func help(commands commands) {
	fmt.Print("gator is used with commands:\n")
	for cmdKey := range commands.cmds {
		fmt.Println(cmdKey)
	}
}
