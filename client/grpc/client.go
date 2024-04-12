package grpc

import (
	"context"
	"crypto/tls"

	"bharvest.io/scheduled-voter/utils"
	tx "github.com/cosmos/cosmos-sdk/types/tx"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func New(host string) *Client {
	return &Client{
		host: host,
	}
}

func (c *Client) Connect(ctx context.Context, secureConnection bool) error {
	options := []grpc.DialOption{grpc.WithBlock()}
	if secureConnection {
		options = append(
			options,
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		)
	} else {
		options = append(
			options,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	}

	conn, err := grpc.DialContext(
		ctx,
		c.host,
		options...,
	)
	if err != nil {
		return err
	}

	c.conn = conn
	c.txServiceClient = tx.NewServiceClient(conn)
	c.govQueryClient = govtypes.NewQueryClient(conn)

	utils.Info("GRPC connected")

	return nil
}

func (c *Client) Terminate(_ context.Context) error {
	err := c.conn.Close()
	utils.Info("GRPC connection terminated")

	return err
}
