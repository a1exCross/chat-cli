package main

import (
	"flag"
	"fmt"
	"github.com/a1exCross/chat-cli/cmd/root"
	"github.com/a1exCross/chat-cli/internal/app"
	"github.com/a1exCross/chat-cli/internal/config"
	"github.com/a1exCross/chat-cli/internal/utils"
	"github.com/fatih/color"
	"log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("falied to load config: %v", err)
	}

}

func main() {
	root.Execute()

	cfg := config.LoadExecConfig()

	if utils.TokenHasExpired(cfg.AccessToken) {
		err := utils.TryReauthorize(app.NewServiceProvider(), cfg)
		if err != nil {
			fmt.Printf(color.RedString("\nNeed log in\n"))

			return
		} else {
			fmt.Printf(color.GreenString("\nHello, %s!\n\n", cfg.Username))
		}
	}
}
