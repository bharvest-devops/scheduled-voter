package app

import (
	"context"
	"fmt"
	"time"

	"bharvest.io/scheduled-voter/client/grpc"
	"bharvest.io/scheduled-voter/client/rpc"
	"bharvest.io/scheduled-voter/script"
	"bharvest.io/scheduled-voter/utils"
	"github.com/cometbft/cometbft/types"
)

func (app *BaseApp) SubProposal(ctx context.Context) {
	rpcClient, err := rpc.New(app.cfg.General.RPC)
	if err != nil {
		app.chErr <- err
		return
	}

	rpcClient.Connect(ctx)
	defer rpcClient.Terminate(ctx)

	// For keep connect webscoket.
	// Webscoket connection will disconnected if no data is sent for an extended period of time.
	go func() {
		chNewBlock, err := rpcClient.Subscribe(ctx, "tm.event = 'NewBlockHeader'")
		if err != nil {
			panic(err)
		}
		for newBlock := range chNewBlock {
			height := newBlock.Data.(types.EventDataNewBlockHeader).Header.Height
			utils.Info(fmt.Sprintf("New block height: %d", height))
		}
	}()

	chEvent, err := rpcClient.Subscribe(ctx, "tm.event = 'Tx' AND message.action = '/cosmos.gov.v1beta1.MsgSubmitProposal'")
	if err != nil {
		app.chErr <- err
		return
	}
	utils.Info("Subscribed to new proposal event")
	for e := range chEvent {
		events := e.Events

		propId := events["submit_proposal.proposal_id"][0]
		app.chSubProposal <- propId

		msg := fmt.Sprintf("New proposal #%s", propId)
		utils.Info(msg)
		utils.SendTg(msg)
	}

	return
}

func (app *BaseApp) ProcVote(ctx context.Context, propId string) error {
	propInfo := fmt.Sprintf(
		"%s prop #%s : %s",
		app.cfg.Vote.Chain,
		propId,
		app.cfg.Vote.Option,
	)

	client := grpc.New(app.cfg.General.GRPC)
	err := client.Connect(ctx, app.cfg.General.GRPCSecureConnection)
	defer client.Terminate(ctx)
	if err != nil {
		utils.Error(err)
		utils.SendTg(fmt.Sprintf("Tx failed(%s): %s", propInfo, err.Error()))
		return err
	}

	cliResp, err := script.Vote(
		app.cfg.Vote.VoteScriptPath,
		app.cfg.Vote.Chain,
		propId,
		app.cfg.Vote.Option,
	)
	if err != nil {
		utils.Error(err)
		utils.SendTg(fmt.Sprintf("Tx failed(%s): %s", propInfo, err.Error()))
		return err
	}

	time.Sleep(10 * time.Second)

	txResp, err := client.GetTx(ctx, cliResp.Txhash)
	if err != nil {
		utils.Error(err)
		utils.SendTg(fmt.Sprintf("Tx failed(%s): %s", propInfo, err.Error()))
		return err
	}

	if txResp.TxResponse.Code == 0 {
		utils.Info("Tx success")
		utils.SendTg(fmt.Sprintf("Tx success(%s): %s", propInfo, txResp.TxResponse.TxHash))
	} else {
		utils.Error(fmt.Errorf("Tx failed(%s): %s", propInfo, txResp.TxResponse.RawLog))
		utils.SendTg(fmt.Sprintf("Tx failed(%s): %s", propInfo, txResp.TxResponse.TxHash))
	}
	return nil
}
