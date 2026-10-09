package biz

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"time"

	pb "transaction/api/transaction/v1"
	"transaction/internal/mapper"
	"transaction/internal/model"

	"github.com/bwmarrin/snowflake"
)

type TransactionRepo interface {
	Create(ctx context.Context, transaction *model.Transaction) error
	GetByID(ctx context.Context, id int64) (*model.Transaction, error)
	GetAll(ctx context.Context) ([]*model.Transaction, error)
	Update(ctx context.Context, transaction *model.Transaction) error
	Delete(ctx context.Context, id int64) error
}

type TransactionUsecase struct {
	repo   TransactionRepo
	node   *snowflake.Node
	logger *slog.Logger
}

func NewTransactionUsecase(
	repo TransactionRepo,
	logger *slog.Logger,
) (*TransactionUsecase, error) {

	node, err := snowflake.NewNode(3)
	if err != nil {
		return nil, err
	}

	return &TransactionUsecase{
		repo:   repo,
		node:   node,
		logger: logger,
	}, nil
}
func (uc *TransactionUsecase) CreateTransaction(
	ctx context.Context,
	req *pb.CreateTransactionRequest,
) (int64, error) {

	if req.PaymentId <= 0 {
		return 0, errors.New("payment id is required")
	}

	if req.AccountId <= 0 {
		return 0, errors.New("account id is required")
	}

	transactionType := strings.ToLower(
		strings.TrimSpace(req.TransactionType),
	)

	if transactionType != "debit" &&
		transactionType != "credit" &&
		transactionType != "fee" &&
		transactionType != "reversal" {
		return 0, errors.New("invalid transaction type")
	}

	amount, err := strconv.ParseFloat(
		req.Amount,
		64,
	)
	if err != nil || amount <= 0 {
		return 0, errors.New("invalid amount")
	}

	currency := strings.ToUpper(
		strings.TrimSpace(req.Currency),
	)

	if len(currency) != 3 {
		return 0, errors.New("currency must contain 3 characters")
	}

	referenceNumber := strings.TrimSpace(
		req.ReferenceNumber,
	)

	if referenceNumber == "" {
		return 0, errors.New("reference number is required")
	}

	if strings.TrimSpace(req.MerchantInfo) != "" {

		var merchant map[string]interface{}

		if err := json.Unmarshal(
			[]byte(req.MerchantInfo),
			&merchant,
		); err != nil {
			return 0, errors.New("invalid merchant info")
		}
	}

	transaction := mapper.ToTransaction(req)

	transaction.TransactionID = uc.node.Generate().Int64()
	transaction.PaymentID = req.PaymentId
	transaction.AccountID = req.AccountId
	transaction.TransactionType = transactionType
	transaction.Amount = amount
	transaction.Currency = currency
	transaction.Status = "pending"
	transaction.ReferenceNumber = referenceNumber
	transaction.TransactionDate = time.Now()
	transaction.CreatedAt = time.Now()



	if strings.TrimSpace(transaction.MerchantInfo) == "" {
    transaction.MerchantInfo = "{}"
}

	if req.RelatedTransactionId > 0 {
		transaction.RelatedTransactionID =
			req.RelatedTransactionId
	}

	if req.PostingDate != "" {

		postingDate, err := time.Parse(
			"2006-01-02",
			req.PostingDate,
		)

		if err != nil {
			return 0, errors.New("invalid posting date")
		}

		transaction.PostingDate = postingDate
	}

	if err := uc.repo.Create(
		ctx,
		transaction,
	); err != nil {

		uc.logger.Error(
			"failed to create transaction",
			"transaction_id",
			transaction.TransactionID,
			"payment_id",
			transaction.PaymentID,
			"error",
			err,
		)

		return 0, err
	}

	uc.logger.Info(
		"transaction created",
		"transaction_id",
		transaction.TransactionID,
		"payment_id",
		transaction.PaymentID,
	)

	return transaction.TransactionID, nil
}

func (uc *TransactionUsecase) GetTransaction(
	ctx context.Context,
	id int64,
) (*model.Transaction, error) {

	if id <= 0 {
		return nil, errors.New(
			"transaction id is required",
		)
	}

	transaction, err := uc.repo.GetByID(
		ctx,
		id,
	)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (uc *TransactionUsecase) GetTransactions(
	ctx context.Context,
) ([]*model.Transaction, error) {

	transactions, err := uc.repo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (uc *TransactionUsecase) UpdateTransaction(
	ctx context.Context,
	req *pb.UpdateTransactionRequest,
) error {

	if req.TransactionId <= 0 {
		return errors.New(
			"transaction id is required",
		)
	}

	transaction, err := uc.repo.GetByID(
		ctx,
		req.TransactionId,
	)

	if err != nil {
		return err
	}

	if req.TransactionType != nil {

		transactionType := strings.ToLower(
			strings.TrimSpace(
				*req.TransactionType,
			),
		)

		if transactionType != "debit" &&
			transactionType != "credit" &&
			transactionType != "fee" &&
			transactionType != "reversal" {

			return errors.New(
				"invalid transaction type",
			)
		}

		transaction.TransactionType = transactionType
	}

	if req.Amount != nil {

		amount, err := strconv.ParseFloat(
			*req.Amount,
			64,
		)

		if err != nil || amount <= 0 {
			return errors.New(
				"invalid amount",
			)
		}

		transaction.Amount = amount
	}

	if req.Currency != nil {

		currency := strings.ToUpper(
			strings.TrimSpace(
				*req.Currency,
			),
		)

		if len(currency) != 3 {
			return errors.New(
				"currency must contain 3 characters",
			)
		}

		transaction.Currency = currency
	}

	if req.Status != nil {

		status := strings.ToLower(
			strings.TrimSpace(
				*req.Status,
			),
		)

		if status != "pending" &&
			status != "posted" &&
			status != "archived" &&
			status != "failed" {

			return errors.New(
				"invalid transaction status",
			)
		}

		transaction.Status = status

		if status == "posted" {
			transaction.PostedAt = time.Now()
		}

		if status == "archived" {
			transaction.ArchivedAt = time.Now()
		}
	}

	if req.Description != nil {
		transaction.Description = *req.Description
	}

	if req.MerchantInfo != nil {

		merchantInfo := strings.TrimSpace(
			*req.MerchantInfo,
		)

		if merchantInfo != "" {

			var merchant map[string]interface{}

			if err := json.Unmarshal(
				[]byte(merchantInfo),
				&merchant,
			); err != nil {
				return errors.New(
					"invalid merchant info",
				)
			}
		}

		transaction.MerchantInfo = merchantInfo
	}

	if req.PostingDate != nil {

		postingDate, err := time.Parse(
			"2006-01-02",
			*req.PostingDate,
		)

		if err != nil {
			return errors.New(
				"invalid posting date",
			)
		}

		transaction.PostingDate = postingDate
	}

	if err := uc.repo.Update(
		ctx,
		transaction,
	); err != nil {

		uc.logger.Error(
			"failed to update transaction",
			"transaction_id",
			transaction.TransactionID,
			"error",
			err,
		)

		return err
	}

	uc.logger.Info(
		"transaction updated",
		"transaction_id",
		transaction.TransactionID,
	)

	return nil
}

func (uc *TransactionUsecase) DeleteTransaction(
	ctx context.Context,
	id int64,
) error {

	if id <= 0 {
		return errors.New(
			"transaction id is required",
		)
	}

	transaction, err := uc.repo.GetByID(
		ctx,
		id,
	)

	if err != nil {
		return err
	}

	if transaction.Status == "posted" {
		return errors.New(
			"posted transaction cannot be deleted",
		)
	}

	if err := uc.repo.Delete(
		ctx,
		id,
	); err != nil {

		uc.logger.Error(
			"failed to delete transaction",
			"transaction_id",
			id,
			"error",
			err,
		)

		return err
	}

	uc.logger.Info(
		"transaction deleted",
		"transaction_id",
		id,
	)

	return nil
}