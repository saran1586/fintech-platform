package client

import (
	"context"

	transactionpb "transaction/api/transaction/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TransactionClient struct {
	client transactionpb.TransactionServiceClient
	conn   *grpc.ClientConn
}

func NewTransactionClient(endpoint string) (*TransactionClient, error) {
	conn, err := grpc.NewClient(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := transactionpb.NewTransactionServiceClient(conn)

	return &TransactionClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *TransactionClient) CreateTransaction(
	ctx context.Context,
	req *transactionpb.CreateTransactionRequest,
) (*transactionpb.CreateTransactionReply, error) {
	return c.client.CreateTransaction(ctx, req)
}