package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MarDoA/gator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}
type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.cmds[cmd.name]
	if !ok {
		return errors.New("command doesn't exist")
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't list users: %w", err)
	}
	for _, user := range users {
		if user.Name == s.cfg.User {
			fmt.Println("* " + user.Name + " (current)")
		} else {
			fmt.Println("* " + user.Name)
		}
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing username")
	}
	//check if name in database
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set user: %w", err)
	}
	fmt.Printf("The user have been set to: %s\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %v (name)", cmd.name)
	}
	uName := cmd.args[0]
	//check if name already in database
	_, err := s.db.GetUser(context.Background(), uName)
	if err == nil {
		return fmt.Errorf("user with name %q already exists", uName)
	}
	if err != sql.ErrNoRows {
		return err
	}

	u, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}
	err = s.cfg.SetUser(uName)
	if err != nil {
		return fmt.Errorf("couldn't set user: %w", err)
	}
	fmt.Println("user was created succefully")
	printUser(u)
	return nil
}
func printUser(user database.User) {
	fmt.Printf(" * Name:    %v\n", user.Name)
}
