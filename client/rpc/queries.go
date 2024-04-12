package rpc

import (
	"context"

	ctypes "github.com/cometbft/cometbft/rpc/core/types"
)

func (c *Client) Subscribe(ctx context.Context, query string) (<-chan ctypes.ResultEvent, error) {
	resultEvent, err := c.RPCClient.Subscribe(ctx, "subscribe", query)
	if err != nil {
		return nil, err
	}

	return resultEvent, nil
}
