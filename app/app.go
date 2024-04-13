package app

import (
	"context"
	"fmt"
	"time"

	"bharvest.io/scheduled-voter/utils"
)

func NewBaseApp(cfg *Config) *BaseApp {
	return &BaseApp{
		cfg:           cfg,
		chSubProposal: make(chan string, 10),
		chErr:         make(chan error),
	}
}

func (app *BaseApp) Run(ctx context.Context) {
	appCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	utils.Info(fmt.Sprintf("Start scheduled voter for %d minutes", app.cfg.Vote.Duration))

	go app.SubProposal(appCtx)
	for {
		select {
		case propId := <-app.chSubProposal:
			// ProcVote will be thread safe
			// because of sequence of vote tx
			app.ProcVote(appCtx, propId)
		case <-time.After(time.Duration(app.cfg.Vote.Duration) * time.Minute):
			utils.Info("Program ended")
			return
		case err := <-app.chErr:
			utils.Error(err)
			return
		}
	}
}
