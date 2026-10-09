package service

import (
	"context"
	"log/slog"

	pb "account/api/account/v1"
	"account/internal/biz"
	"account/internal/mapper"
)

type AccountService struct {
	pb.UnimplementedAccountServiceServer
	uc     *biz.AccountUsecase
	logger *slog.Logger
}

func NewAccountService(uc *biz.AccountUsecase, logger *slog.Logger) *AccountService {
	return &AccountService{
		uc:     uc,
		logger: logger,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountReply, error) {
	_, err := s.uc.CreateAccount(ctx, req)
	if err != nil {
		s.logger.Error("failed to create account", "error", err)
		return nil, err
	}

	return &pb.CreateAccountReply{
		Message: "Account created successfully",
	}, nil
}

func (s *AccountService) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountReply, error) {
	account, err := s.uc.GetAccount(ctx, req.AccountId)
	if err != nil {
		s.logger.Error("failed to get account", "account_id", req.AccountId, "error", err)
		return nil, err
	}

	return mapper.ToGetAccountResponse(account), nil
}

func (s *AccountService) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.GetAccountsReply, error) {
	accounts, err := s.uc.GetAccounts(ctx)
	if err != nil {
		s.logger.Error("failed to get accounts", "error", err)
		return nil, err
	}

	return mapper.ToGetAccountsResponse(accounts), nil
}

func (s *AccountService) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountReply, error) {
	account, err := s.uc.GetAccount(ctx, req.AccountId)
	if err != nil {
		s.logger.Error("failed to get account", "error", err)
		return nil, err
	}

	mapper.ToUpdateAccount(req, account)

	if err := s.uc.UpdateAccount(ctx, account); err != nil {
		s.logger.Error("failed to update account", "error", err)
		return nil, err
	}

	return &pb.UpdateAccountReply{
		Message: "Account updated successfully",
	}, nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountReply, error) {
	err := s.uc.DeleteAccount(ctx, req.AccountId)
	if err != nil {
		s.logger.Error("failed to delete account", "account_id", req.AccountId, "error", err)
		return nil, err
	}

	return &pb.DeleteAccountReply{
		Message: "Account deleted successfully",
	}, nil
}


func (s *AccountService) DebitAccount(
	ctx context.Context,
	req *pb.DebitAccountRequest,
) (*pb.DebitAccountReply, error) {

	if err := s.uc.DebitAccount(
		ctx,
		req.AccountId,
		req.Amount,
	); err != nil {
		s.logger.Error(
			"failed to debit account",
			"account_id",
			req.AccountId,
			"amount",
			req.Amount,
			"error",
			err,
		)

		return nil, err
	}

	return &pb.DebitAccountReply{
		Message: "Account debited successfully",
	}, nil
}
func (s *AccountService) CreditAccount(
	ctx context.Context,
	req *pb.CreditAccountRequest,
) (*pb.CreditAccountReply, error) {

	_, err := s.uc.CreditAccount(
		ctx,
		req.AccountId,
		req.Amount,
	)

	if err != nil {
		s.logger.Error(
			"failed to credit account",
			"account_id", req.AccountId,
			"amount", req.Amount,
			"error", err,
		)

		return nil, err
	}

	return &pb.CreditAccountReply{
		Message: "Account credited successfully",
	}, nil
}