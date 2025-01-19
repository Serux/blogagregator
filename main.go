package main

import (
	"fmt"

	"github.com/serux/blogagregator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("Serux")
	cfg = config.Read()
	fmt.Println(cfg)

}
