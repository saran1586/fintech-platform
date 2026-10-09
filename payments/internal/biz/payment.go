package biz

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"strings"
   "fmt"
	"time"
    accountpb "account/api/account/v1"
	transactionpb "transaction/api/transaction/v1"
	"payments/internal/event"
"payments/internal/messaging"
	pb "payments/api/payments/v2"
	"payments/internal/mapper"
	"payments/internal/model"
	

	"github.com/bwmarrin/snowflake"
)

type PaymentRepo interface {
	Create(ctx context.Context, payment *model.Payment) error
	GetByID(ctx context.Context, id int64) (*model.Payment, error)
	GetAll(ctx context.Context) ([]*model.Payment, error)
	Update(ctx context.Context, payment *model.Payment) error
	Delete(ctx context.Context, id int64) error
}

type AccountClient interface {
	GetAccount(
		ctx context.Context,
		accountID int64,
	) (*accountpb.GetAccountReply, error)

	DebitAccount(
		ctx context.Context,
		accountID int64,
		amount string,
	) (*accountpb.DebitAccountReply, error)

	CreditAccount(
		ctx context.Context,
		accountID int64,
		amount string,
	) (*accountpb.CreditAccountReply, error)
}

type TransactionClient interface {
	CreateTransaction(
		ctx context.Context,
		req *transactionpb.CreateTransactionRequest,
	) (*transactionpb.CreateTransactionReply, error)
}

type PaymentUsecase struct {
	repo   PaymentRepo
		accountClient AccountClient
		transactionClient TransactionClient
		 publisher         *messaging.Publisher
	node   *snowflake.Node
	logger *slog.Logger
}

func NewPaymentUsecase(
	repo PaymentRepo,
	accountClient AccountClient,
	transactionClient TransactionClient,
	 publisher *messaging.Publisher,
	logger *slog.Logger,
) (*PaymentUsecase, error) {

	node, err := snowflake.NewNode(4)
	if err != nil {
		return nil, err
	}

	return &PaymentUsecase{
		repo:   repo,
		accountClient: accountClient,
		transactionClient: transactionClient,
		 publisher:         publisher,
		node:   node,
		logger: logger,
	}, nil
}
func (uc *PaymentUsecase) CreatePayment(
	ctx context.Context,
	req *pb.CreatePaymentRequest,
) error {

	paymentType := strings.ToLower(
		strings.TrimSpace(req.PaymentType),
	)

	if paymentType != "transfer" &&
		paymentType != "payment" &&
		paymentType != "topup" {
		return errors.New("invalid payment type")
	}

	if paymentType == "transfer" ||
		paymentType == "payment" {

		if req.FromAccountId <= 0 {
			return errors.New("from account id is required")
		}

		if req.ToAccountId <= 0 {
			return errors.New("to account id is required")
		}

		if req.FromAccountId == req.ToAccountId {
			return errors.New("from and to account cannot be same")
		}
	}

	if paymentType == "topup" &&
		req.ToAccountId <= 0 {
		return errors.New("to account id is required")
	}

	amount, err := strconv.ParseFloat(
		req.Amount,
		64,
	)
	if err != nil || amount <= 0 {
		return errors.New("invalid amount")
	}

	currency := strings.ToUpper(
		strings.TrimSpace(req.Currency),
	)

	if len(currency) != 3 {
		return errors.New("currency must contain 3 characters")
	}

	paymentMethod := strings.ToLower(
		strings.TrimSpace(req.PaymentMethod),
	)

	if paymentMethod != "card" &&
		paymentMethod != "bank_transfer" &&
		paymentMethod != "wallet" {
		return errors.New("invalid payment method")
	}

	referenceNumber := strings.TrimSpace(
		req.ReferenceNumber,
	)

	if referenceNumber == "" {
		return errors.New("reference number is required")
	}

	payment := mapper.ToPayment(req)

	payment.PaymentID = uc.node.Generate().Int64()
	payment.PaymentType = paymentType
	payment.Amount = amount
	payment.Currency = currency
	payment.Status = "initiated"
	payment.PaymentMethod = paymentMethod
	payment.ReferenceNumber = referenceNumber
	payment.CreatedAt = time.Now()

	if req.ScheduledAt != "" {

		scheduledAt, err := time.Parse(
			time.RFC3339,
			req.ScheduledAt,
		)

		if err != nil {
			return errors.New("invalid scheduled at")
		}

		payment.ScheduledAt = scheduledAt
	}

	if err := uc.repo.Create(ctx, payment); err != nil {

		uc.logger.Error(
			"failed to create payment",
			"payment_id",
			payment.PaymentID,
			"error",
			err,
		)

		return err
	}

	err = uc.publisher.PublishPaymentInitiated(
		event.PaymentInitiatedEvent{
			PaymentID:       payment.PaymentID,
			FromAccountID:   req.FromAccountId,
			ToAccountID:     req.ToAccountId,
			Amount:          req.Amount,
			Currency:        payment.Currency,
			ReferenceNumber: payment.ReferenceNumber,
		},
	)

	if err != nil {
		uc.logger.Error(
			"failed to publish payment initiated event",
			"payment_id",
			payment.PaymentID,
			"error",
			err,
		)
	}

	payment.Status = "processing"

	if err := uc.repo.Update(ctx, payment); err != nil {
		return err
	}

	if paymentType == "transfer" ||
		paymentType == "payment" {

		_, err := uc.accountClient.GetAccount(
			ctx,
			req.FromAccountId,
		)

		if err != nil {
			payment.Status = "failed"
			uc.repo.Update(ctx, payment)

			return errors.New("from account not found")
		}

		_, err = uc.accountClient.GetAccount(
			ctx,
			req.ToAccountId,
		)

		if err != nil {
			payment.Status = "failed"
			uc.repo.Update(ctx, payment)

			return errors.New("to account not found")
		}
	}

	if paymentType == "topup" {

		_, err := uc.accountClient.GetAccount(
			ctx,
			req.ToAccountId,
		)

		if err != nil {
			payment.Status = "failed"
			uc.repo.Update(ctx, payment)

			return errors.New("to account not found")
		}
	}

	if paymentType == "transfer" ||
		paymentType == "payment" {

		_, err = uc.accountClient.DebitAccount(
			ctx,
			req.FromAccountId,
			req.Amount,
		)

		if err != nil {

			uc.logger.Error(
				"failed to debit from account",
				"account_id",
				req.FromAccountId,
				"amount",
				req.Amount,
				"error",
				err,
			)

			payment.Status = "failed"

			if updateErr := uc.repo.Update(ctx, payment); updateErr != nil {
				uc.logger.Error(
					"failed to update payment status",
					"payment_id",
					payment.PaymentID,
					"error",
					updateErr,
				)
			}

			return errors.New("failed to debit from account")
		}

		debitTransaction, err := uc.transactionClient.CreateTransaction(
			ctx,
			&transactionpb.CreateTransactionRequest{
				PaymentId:       payment.PaymentID,
				AccountId:       req.FromAccountId,
				TransactionType: "debit",
				Amount:          req.Amount,
				Currency:        payment.Currency,
				ReferenceNumber: payment.ReferenceNumber,
			},
		)

		if err != nil {

			fmt.Println("TRANSACTION GRPC ERROR:", err)

			uc.logger.Error(
				"failed to create debit transaction",
				"payment_id",
				payment.PaymentID,
				"account_id",
				req.FromAccountId,
				"error",
				err,
			)

			_, rollbackErr := uc.accountClient.CreditAccount(
				ctx,
				req.FromAccountId,
				req.Amount,
			)

			if rollbackErr != nil {
				uc.logger.Error(
					"failed to rollback debit",
					"account_id",
					req.FromAccountId,
					"amount",
					req.Amount,
					"error",
					rollbackErr,
				)
			}

			payment.Status = "failed"

			if updateErr := uc.repo.Update(ctx, payment); updateErr != nil {
				uc.logger.Error(
					"failed to update payment status",
					"payment_id",
					payment.PaymentID,
					"error",
					updateErr,
				)
			}

			return errors.New("failed to create debit transaction")
		}

		_, err = uc.accountClient.CreditAccount(
			ctx,
			req.ToAccountId,
			req.Amount,
		)

		if err != nil {

			uc.logger.Error(
				"failed to credit to account",
				"account_id",
				req.ToAccountId,
				"amount",
				req.Amount,
				"error",
				err,
			)

			_, rollbackErr := uc.accountClient.CreditAccount(
				ctx,
				req.FromAccountId,
				req.Amount,
			)

			if rollbackErr != nil {
				uc.logger.Error(
					"failed to rollback debit",
					"account_id",
					req.FromAccountId,
					"amount",
					req.Amount,
					"error",
					rollbackErr,
				)
			}

			payment.Status = "failed"

			if updateErr := uc.repo.Update(ctx, payment); updateErr != nil {
				uc.logger.Error(
					"failed to update payment status",
					"payment_id",
					payment.PaymentID,
					"error",
					updateErr,
				)
			}

			return errors.New("failed to credit to account")
		}

		_, err = uc.transactionClient.CreateTransaction(
			ctx,
			&transactionpb.CreateTransactionRequest{
				PaymentId:            payment.PaymentID,
				AccountId:            req.ToAccountId,
				RelatedTransactionId: debitTransaction.TransactionId,
				TransactionType:      "credit",
				Amount:               req.Amount,
				Currency:             payment.Currency,
				ReferenceNumber:      payment.ReferenceNumber,
			},
		)

		if err != nil {

			uc.logger.Error(
				"failed to create credit transaction",
				"payment_id",
				payment.PaymentID,
				"account_id",
				req.ToAccountId,
				"error",
				err,
			)

			_, rollbackErr := uc.accountClient.DebitAccount(
				ctx,
				req.ToAccountId,
				req.Amount,
			)

			if rollbackErr != nil {
				uc.logger.Error(
					"failed to rollback credit",
					"account_id",
					req.ToAccountId,
					"amount",
					req.Amount,
					"error",
					rollbackErr,
				)
			}

			_, rollbackErr = uc.accountClient.CreditAccount(
				ctx,
				req.FromAccountId,
				req.Amount,
			)

			if rollbackErr != nil {
				uc.logger.Error(
					"failed to rollback debit",
					"account_id",
					req.FromAccountId,
					"amount",
					req.Amount,
					"error",
					rollbackErr,
				)
			}

			payment.Status = "failed"

			if updateErr := uc.repo.Update(ctx, payment); updateErr != nil {
				uc.logger.Error(
					"failed to update payment status",
					"payment_id",
					payment.PaymentID,
					"error",
					updateErr,
				)
			}

			return errors.New("failed to create credit transaction")
		}
	}

	payment.Status = "completed"
	payment.ProcessedAt = time.Now()

	if err := uc.repo.Update(ctx, payment); err != nil {

		uc.logger.Error(
			"failed to update completed payment",
			"payment_id",
			payment.PaymentID,
			"error",
			err,
		)

		return err
	}

	uc.logger.Info(
		"payment completed",
		"payment_id",
		payment.PaymentID,
	)

	return nil
}
func (uc *PaymentUsecase) GetPayment(ctx context.Context,id int64,) (*model.Payment, error) {

	if id <= 0 {
		return nil, errors.New(
			"payment id is required",
		)
	}

	return uc.repo.GetByID(ctx, id)
}

func (uc *PaymentUsecase) GetPayments(ctx context.Context,) ([]*model.Payment, error) {

	return uc.repo.GetAll(ctx)
}

func (uc *PaymentUsecase) UpdatePayment(ctx context.Context,req *pb.UpdatePaymentRequest,) error {

	if req.PaymentId <= 0 {
		return errors.New(
			"payment id is required",
		)
	}

	payment, err := uc.repo.GetByID(
		ctx,
		req.PaymentId,
	)

	if err != nil {
		return err
	}

	if req.PaymentType != nil {

		paymentType := strings.ToLower(
			strings.TrimSpace(
				*req.PaymentType,
			),
		)

		if paymentType != "transfer" &&
			paymentType != "payment" &&
			paymentType != "topup" {
			return errors.New(
				"invalid payment type",
			)
		}
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
	}

	if req.Status != nil {

		status := strings.ToLower(
			strings.TrimSpace(
				*req.Status,
			),
		)

		if status != "initiated" &&
			status != "processing" &&
			status != "completed" &&
			status != "failed" &&
			status != "reversed" {

			return errors.New(
				"invalid payment status",
			)
		}
	}

	if req.PaymentMethod != nil {

		paymentMethod := strings.ToLower(
			strings.TrimSpace(
				*req.PaymentMethod,
			),
		)

		if paymentMethod != "card" &&
			paymentMethod != "bank_transfer" &&
			paymentMethod != "wallet" {

			return errors.New(
				"invalid payment method",
			)
		}
	}

	if req.ScheduledAt != nil {

		_, err := time.Parse(
			time.RFC3339,
			*req.ScheduledAt,
		)

		if err != nil {
			return errors.New(
				"invalid scheduled at",
			)
		}
	}

	mapper.ToUpdatePayment(
		req,
		payment,
	)

	if payment.Status == "completed" &&
		payment.ProcessedAt.IsZero() {

		payment.ProcessedAt = time.Now()
	}

	if err := uc.repo.Update(
		ctx,
		payment,
	); err != nil {

		uc.logger.Error(
			"failed to update payment",
			"payment_id",
			payment.PaymentID,
			"error",
			err,
		)

		return err
	}

	return nil
}

func (uc *PaymentUsecase) DeletePayment(
	ctx context.Context,
	id int64,
) error {

	if id <= 0 {
		return errors.New(
			"payment id is required",
		)
	}

	payment, err := uc.repo.GetByID(
		ctx,
		id,
	)

	if err != nil {
		return err
	}

	if payment.Status == "completed" {
		return errors.New(
			"completed payment cannot be deleted",
		)
	}

	return uc.repo.Delete(ctx, id)
}