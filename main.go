package main

import (
	"context"
	"fmt"
	"os"

	"bharvest.io/scheduled-voter/app"
	"bharvest.io/scheduled-voter/utils"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	ctx := context.Background()

	f, err := os.ReadFile("config.toml")
	if err != nil {
		utils.Error(err)
		panic(err)
	}
	cfg := app.Config{}
	err = toml.Unmarshal(f, &cfg)
	if err != nil {
		utils.Error(err)
		panic(err)
	}

	tgTitle := fmt.Sprintf("🤖 Scheduled Voter 🤖")
	utils.SetTg(cfg.Tg.Enable, tgTitle, cfg.Tg.Token, cfg.Tg.ChatID)

	baseapp := app.NewBaseApp(&cfg)
	baseapp.Run(ctx)
}
