package client

import (
	"context"

	"account/internal/conf"
	customerpb "customer/api/customer/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CustomerClient struct {
	client customerpb.CustomerServiceClient
}

func NewCustomerClient(
	cfg *conf.CustomerClient,
) (*CustomerClient, error) {

	conn, err := grpc.NewClient(
		cfg.Endpoint,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		return nil, err
	}

	return &CustomerClient{
		client: customerpb.NewCustomerServiceClient(conn),
	}, nil
}

func (c *CustomerClient) GetCustomer(ctx context.Context, customerID int64) (*customerpb.GetCustomerReply, error) {

	return c.client.GetCustomer(ctx, &customerpb.GetCustomerRequest{
		CustomerId: customerID,
	},
	)
}
