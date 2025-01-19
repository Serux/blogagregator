package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/serux/blogagregator/internal/config"
	"github.com/serux/blogagregator/internal/database"
	"github.com/serux/blogagregator/rss"
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

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("format: addfeed name url")
	}
	user, err := s.Db.GetOneUserByName(context.Background(), s.Config.Current_user_name)
	if err != nil {
		return err
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
		Url:       cmd.Arguments[1],
		UserID:    user.ID,
	}
	ret, err := s.Db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Println(ret)
	return nil
}

func HandlerAgg(s *State, cmd Command) error {

	ret, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")

	if err != nil {
		return err
	}
	fmt.Println(ret)

	return nil
}

func HandlerReset(s *State, cmd Command) error {

	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		fmt.Println("Cannot reset users:")
		return err
	}

	fmt.Println("Reset successful")

	return nil
}

func HandlerGetUsers(s *State, cmd Command) error {

	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		fmt.Println("Cannot get users:")
		return err
	}

	for _, user := range users {
		current := ""
		if user.Name == s.Config.Current_user_name {
			current = "(current)"
		}
		fmt.Printf("* %v %v\n", user.Name, current)
	}

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
