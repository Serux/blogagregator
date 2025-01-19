package main

import (
	"fmt"
	"os"

	"github.com/serux/blogagregator/commands"
	"github.com/serux/blogagregator/internal/config"
)

func main() {
	cfg := config.Read()
	stt := commands.State{Config: &cfg}
	cmds := commands.Commands{CommandsHandlers: map[string]func(*commands.State, commands.Command) error{}}
	cmds.Register("login", commands.HandlerLogin)
	if len(os.Args) < 2 {
		fmt.Println("Not enough args.")
		os.Exit(1)
	}
	c := commands.Command{Name: os.Args[1], Arguments: os.Args[2:]}

	err := cmds.Run(&stt, c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
