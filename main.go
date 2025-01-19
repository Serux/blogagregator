package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/serux/blogagregator/commands"
	"github.com/serux/blogagregator/internal/config"
	"github.com/serux/blogagregator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Read()
	stt := commands.State{Config: &cfg}
	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	stt.Db = database.New(db)

	cmds := commands.Commands{CommandsHandlers: map[string]func(*commands.State, commands.Command) error{}}
	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerReset)
	cmds.Register("users", commands.HandlerGetUsers)

	if len(os.Args) < 2 {
		fmt.Println("Not enough args.")
		os.Exit(1)
	}
	c := commands.Command{Name: os.Args[1], Arguments: os.Args[2:]}

	err = cmds.Run(&stt, c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
