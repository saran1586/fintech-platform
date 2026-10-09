package client

import (
	"context"

	accountpb "account/api/account/v1"
	"payments/internal/conf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AccountClient struct {
	client accountpb.AccountServiceClient
}

func NewAccountClient(
	cfg *conf.Data_AccountClient,
) (*AccountClient, error) {

	conn, err := grpc.NewClient(
		cfg.Endpoint,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		return nil, err
	}

	return &AccountClient{
		client: accountpb.NewAccountServiceClient(conn),
	}, nil
}

func (c *AccountClient) GetAccount(
	ctx context.Context,
	accountID int64,
) (*accountpb.GetAccountReply, error) {

	return c.client.GetAccount(
		ctx,
		&accountpb.GetAccountRequest{
			AccountId: accountID,
		},
	)
}

func (c *AccountClient) DebitAccount(
	ctx context.Context,
	accountID int64,
	amount string,
) (*accountpb.DebitAccountReply, error) {

	return c.client.DebitAccount(
		ctx,
		&accountpb.DebitAccountRequest{
			AccountId: accountID,
			Amount:    amount,
		},
	)
}

func (c *AccountClient) CreditAccount(
	ctx context.Context,
	accountID int64,
	amount string,
) (*accountpb.CreditAccountReply, error) {

	return c.client.CreditAccount(
		ctx,
		&accountpb.CreditAccountRequest{
			AccountId: accountID,
			Amount:    amount,
		},
	)
}