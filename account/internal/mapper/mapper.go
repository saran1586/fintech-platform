package mapper

import (
	"strconv"
	"time"

	pb "account/api/account/v1"
	"account/internal/model"
)

func ToAccount(req *pb.CreateAccountRequest) *model.Account {
	return &model.Account{
		CustomerID:  req.CustomerId,
		AccountType: req.AccountType,
		Currency:    req.Currency,
	}
}

func ToAccountResponse(account *model.Account) *pb.Account {
	return &pb.Account{
		AccountId:         account.AccountID,
		CustomerId:        account.CustomerID,
		AccountNumber:     account.AccountNumber,
		AccountType:       account.AccountType,
		Currency:          account.Currency,
		Status:            account.Status,
		AvailableBalance:  strconv.FormatFloat(account.AvailableBalance, 'f', 4, 64),
		PendingBalance:    strconv.FormatFloat(account.PendingBalance, 'f', 4, 64),
		CreditLimit:       strconv.FormatFloat(account.CreditLimit, 'f', 4, 64),
		CreatedAt:         account.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         account.UpdatedAt.Format(time.RFC3339),
		LastTransactionAt: account.LastTransactionAt.Format(time.RFC3339),
	}
}

func ToGetAccountResponse(account *model.Account) *pb.GetAccountReply {
	return &pb.GetAccountReply{
		AccountId:         account.AccountID,
		CustomerId:        account.CustomerID,
		AccountNumber:     account.AccountNumber,
		AccountType:       account.AccountType,
		Currency:          account.Currency,
		Status:            account.Status,
		AvailableBalance:  strconv.FormatFloat(account.AvailableBalance, 'f', 4, 64),
		PendingBalance:    strconv.FormatFloat(account.PendingBalance, 'f', 4, 64),
		CreditLimit:       strconv.FormatFloat(account.CreditLimit, 'f', 4, 64),
		CreatedAt:         account.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         account.UpdatedAt.Format(time.RFC3339),
		LastTransactionAt: account.LastTransactionAt.Format(time.RFC3339),
	}
}

func ToGetAccountsResponse(accounts []*model.Account) *pb.GetAccountsReply {
	result := make([]*pb.Account, 0, len(accounts))

	for _, account := range accounts {
		result = append(result, ToAccountResponse(account))
	}

	return &pb.GetAccountsReply{
		Accounts: result,
	}
}

func ToUpdateAccount(req *pb.UpdateAccountRequest, account *model.Account) {
	if req.AccountType != nil {
		account.AccountType = *req.AccountType
	}

	if req.Currency != nil {
		account.Currency = *req.Currency
	}

	if req.Status != nil {
		account.Status = *req.Status
	}

	if req.CreditLimit != nil {
		value, err := strconv.ParseFloat(*req.CreditLimit, 64)
		if err == nil {
			account.CreditLimit = value
		}
	}
}