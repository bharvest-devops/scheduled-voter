package grpc

import (
	tx "github.com/cosmos/cosmos-sdk/types/tx"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"google.golang.org/grpc"
)

type Client struct {
	host string
	conn *grpc.ClientConn
	txServiceClient tx.ServiceClient
	govQueryClient govtypes.QueryClient
}
