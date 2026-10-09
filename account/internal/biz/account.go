package biz

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
    pb "account/api/account/v1"
	"account/internal/mapper"
	"account/internal/model"
    "account/internal/client"
	"github.com/bwmarrin/snowflake"
	"log/slog"
)

type AccountRepo interface {
	Create(ctx context.Context, account *model.Account) error
	GetByID(ctx context.Context, id int64) (*model.Account, error)
	GetAll(ctx context.Context) ([]*model.Account, error)
	Update(ctx context.Context, account *model.Account) error
	Delete(ctx context.Context, id int64) error
}


type AccountUsecase struct {
	repo   AccountRepo
	customerClient *client.CustomerClient
	node   *snowflake.Node
	logger *slog.Logger

}

func NewAccountUsecase(repo AccountRepo,customerClient *client.CustomerClient, logger *slog.Logger, ) (*AccountUsecase, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}

	return &AccountUsecase{
		repo:   repo,
		customerClient: customerClient,
		node:   node,
		logger: logger,
		
	}, nil
}
func (uc *AccountUsecase) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*model.Account, error) {
	if req.CustomerId <= 0 {
		return nil, errors.New("customer id is required")
	}

	customer, err := uc.customerClient.GetCustomer(
		ctx,
		req.CustomerId,
	)
	if err != nil {
		uc.logger.Error(
			"failed to get customer",
			"customer_id",
			req.CustomerId,
			"error",
			err,
		)

		return nil, errors.New("customer not found")
	}

	if customer.Status != "active" {
		return nil, errors.New("customer is not active")
	}

	if customer.KycStatus != "verified" {
		return nil, errors.New("customer kyc is not verified")
	}

	accountType := strings.ToLower(strings.TrimSpace(req.AccountType))

	if accountType != "checking" &&
		accountType != "savings" &&
		accountType != "credit" &&
		accountType != "wallet" {
		return nil, errors.New("invalid account type")
	}

	currency := strings.ToUpper(strings.TrimSpace(req.Currency))

	if len(currency) != 3 {
		return nil, errors.New("currency must be 3 characters")
	}

	account := mapper.ToAccount(req)

	account.AccountID = uc.node.Generate().Int64()
	account.AccountNumber = "ACC" + strconv.FormatInt(account.AccountID, 10)
	account.AccountType = accountType
	account.Currency = currency
	account.Status = "active"
	account.AvailableBalance = 0
	account.PendingBalance = 0
	account.CreditLimit = 0
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	if err := uc.repo.Create(ctx, account); err != nil {
		uc.logger.Error("failed to create account", "error", err)
		return nil, err
	}

	uc.logger.Info("account created", "account_id", account.AccountID)

	return account, nil
}


func (uc *AccountUsecase) GetAccount(ctx context.Context, id int64) (*model.Account, error) {
	if id <= 0 {
		return nil, errors.New("account id is required")
	}

	account, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("failed to get account", "account_id", id, "error", err)
		return nil, err
	}

	return account, nil
}

func (uc *AccountUsecase) GetAccounts(ctx context.Context) ([]*model.Account, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *AccountUsecase) UpdateAccount(ctx context.Context, account *model.Account) error {
	return uc.repo.Update(ctx, account)
}

func (uc *AccountUsecase) DeleteAccount(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *AccountUsecase) DebitAccount(ctx context.Context,accountID int64,amount string,) error {

	if accountID <= 0 {
		return errors.New("account id is required")
	}

	value, err := strconv.ParseFloat(amount, 64)
	if err != nil || value <= 0 {
		return errors.New("invalid amount")
	}

	account, err := uc.repo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}

	if account.Status != "active" {
		return errors.New("account is not active")
	}

	if account.AvailableBalance < value {
		return errors.New("insufficient balance")
	}

	account.AvailableBalance -= value
	account.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, account); err != nil {
		uc.logger.Error(
			"failed to debit account",
			"account_id", accountID,
			"amount", value,
			"error", err,
		)

		return err
	}

	uc.logger.Info(
		"account debited",
		"account_id", accountID,
		"amount", value,
	)

	return nil
}

func (uc *AccountUsecase) CreditAccount(ctx context.Context,accountID int64,amount string,) (*model.Account, error) {

	if accountID <= 0 {
		return nil, errors.New("account id is required")
	}

	value, err := strconv.ParseFloat(amount, 64)
	if err != nil || value <= 0 {
		return nil, errors.New("invalid amount")
	}

	account, err := uc.repo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if account.Status != "active" {
		return nil, errors.New("account is not active")
	}

	account.AvailableBalance += value

	if err := uc.repo.Update(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}