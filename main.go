package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/MarDoA/gator/internal/config"
	"github.com/MarDoA/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("couldn't read config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbQueries := database.New(db)
	s := &state{cfg: &cfg, db: dbQueries}

	cmds := commands{cmds: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAggregate)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerListFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnFollow))
	cmds.register("browse", middlewareLoggedIn(handleBrowse))
	args := os.Args

	if len(args) < 2 {
		log.Fatal("no commands passed")
	}
	err = cmds.run(s, command{name: args[1], args: args[2:]})
	if err != nil {
		log.Fatal(err)
	}

}
