package service

import (
	"context"
	"log/slog"

	pb "transaction/api/transaction/v1"
	"transaction/internal/biz"
	"transaction/internal/mapper"
)

type TransactionService struct {
	pb.UnimplementedTransactionServiceServer

	uc     *biz.TransactionUsecase
	logger *slog.Logger
}

func NewTransactionService(uc *biz.TransactionUsecase,logger *slog.Logger,) *TransactionService {
	return &TransactionService{
		uc:     uc,
		logger: logger,
	}
}

func (s *TransactionService) CreateTransaction(
	ctx context.Context,
	req *pb.CreateTransactionRequest,
) (*pb.CreateTransactionReply, error) {

	transactionID, err := s.uc.CreateTransaction(
		ctx,
		req,
	)

	if err != nil {
		s.logger.Error(
			"failed to create transaction",
			"error", err,
		)

		return nil, err
	}

	return &pb.CreateTransactionReply{
		Message:       "Transaction created successfully",
		TransactionId: transactionID,
	}, nil
}
func (s *TransactionService) GetTransaction(ctx context.Context,req *pb.GetTransactionRequest,) (*pb.GetTransactionReply, error) {

	transaction, err := s.uc.GetTransaction(
		ctx,
		req.TransactionId,
	)
	if err != nil {
		s.logger.Error(
			"failed to get transaction",
			"transaction_id", req.TransactionId,
			"error", err,
		)
		return nil, err
	}

	return mapper.ToGetTransactionResponse(transaction), nil
}


func (s *TransactionService) GetTransactions(ctx context.Context,req *pb.GetTransactionsRequest,) (*pb.GetTransactionsReply, error) {

	transactions, err := s.uc.GetTransactions(ctx)
	if err != nil {
		s.logger.Error(
			"failed to get transactions",
			"error", err,
		)
		return nil, err
	}

	return mapper.ToGetTransactionsResponse(transactions), nil
}

func (s *TransactionService) UpdateTransaction(ctx context.Context,req *pb.UpdateTransactionRequest,) (*pb.UpdateTransactionReply, error) {

	if err := s.uc.UpdateTransaction(ctx, req); err != nil {
		s.logger.Error(
			"failed to update transaction",
			"transaction_id", req.TransactionId,
			"error", err,
		)

		return nil, err
	}

	return &pb.UpdateTransactionReply{
		Message: "Transaction updated successfully",
	}, nil
}

func (s *TransactionService) DeleteTransaction(ctx context.Context,req *pb.DeleteTransactionRequest,) (*pb.DeleteTransactionReply, error) {

	if err := s.uc.DeleteTransaction(
		ctx,
		req.TransactionId,
	); err != nil {
		s.logger.Error(
			"failed to delete transaction",
			"transaction_id", req.TransactionId,
			"error", err,
		)
		return nil, err
	}

	return &pb.DeleteTransactionReply{
		Message: "Transaction deleted successfully",
	}, nil
}