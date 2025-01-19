package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/serux/blogagregator/internal/config"
	"github.com/serux/blogagregator/internal/database"
)

type State struct {
	Config *config.Config
	Db     *database.Queries
}

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	CommandsHandlers map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandsHandlers[name] = f
}
func (c *Commands) Run(s *State, cmd Command) error {
	fun, ok := c.CommandsHandlers[cmd.Name]
	if !ok {
		return fmt.Errorf("command not handled")
	}
	return fun(s, cmd)
}

func HandlerReset(s *State, cmd Command) error {

	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		fmt.Println("Cannot reset users:", err)
		os.Exit(1)
	}

	fmt.Println("Reset successful")

	return nil
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("username required")
	}
	_, err := s.Db.GetOneUserByName(context.Background(), cmd.Arguments[0])
	if err != nil {
		fmt.Println("Cannot find user:", err)
		os.Exit(1)
	}
	s.Config.SetUser(cmd.Arguments[0])
	fmt.Println("Username Set")

	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("name required")
	}

	us, err := s.Db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Arguments[0],
		},
	)
	if err != nil {
		fmt.Println("User exists: ", err)
		os.Exit(1)
	}

	s.Config.SetUser(cmd.Arguments[0])
	fmt.Println("User created:", us)

	return nil
}
