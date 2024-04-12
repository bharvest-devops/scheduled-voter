package grpc

import (
	"context"

	tx "github.com/cosmos/cosmos-sdk/types/tx"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func (c *Client) GetTx(ctx context.Context, hash string) (*tx.GetTxResponse, error) {
	resp, err := c.txServiceClient.GetTx(
		ctx,
		&tx.GetTxRequest{
			Hash: hash,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetProposal(ctx context.Context, proposalId uint64) (*govtypes.QueryProposalResponse, error) {
	resp, err := c.govQueryClient.Proposal(
		ctx,
		&govtypes.QueryProposalRequest{
			ProposalId: proposalId,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
