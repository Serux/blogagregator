package commands

import (
	"fmt"

	"github.com/serux/blogagregator/internal/config"
)

type State struct {
	Config *config.Config
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

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("username required")
	}
	s.Config.SetUser(cmd.Arguments[0])
	fmt.Println("Username Set")

	return nil
}
