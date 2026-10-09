package service

import (
	"context"
	"log/slog"

	pb "payments/api/payments/v2"
	"payments/internal/biz"
	"payments/internal/mapper"
)

type PaymentsService struct {
	pb.UnimplementedPaymentsServiceServer

	uc     *biz.PaymentUsecase
	logger *slog.Logger
}

func NewPaymentsService(
	uc *biz.PaymentUsecase,
	logger *slog.Logger,
) *PaymentsService {

	return &PaymentsService{
		uc:     uc,
		logger: logger,
	}
}

func (s *PaymentsService) CreatePayment(ctx context.Context,req *pb.CreatePaymentRequest,) (*pb.CreatePaymentReply, error) {
	if err := s.uc.CreatePayment(
		ctx,
		req,
	); err != nil {

		s.logger.Error(
			"failed to create payment",
			"error",
			err,
		)

		return nil, err
	}

	return &pb.CreatePaymentReply{
		Message: "Payment created successfully",
	}, nil
}

func (s *PaymentsService) GetPayment(
	ctx context.Context,
	req *pb.GetPaymentRequest,
) (*pb.GetPaymentReply, error) {

	payment, err := s.uc.GetPayment(
		ctx,
		req.PaymentId,
	)

	if err != nil {

		s.logger.Error(
			"failed to get payment",
			"payment_id",
			req.PaymentId,
			"error",
			err,
		)

		return nil, err
	}

	return mapper.ToGetPaymentResponse(
		payment,
	), nil
}

func (s *PaymentsService) GetPayments(
	ctx context.Context,
	req *pb.GetPaymentsRequest,
) (*pb.GetPaymentsReply, error) {

	payments, err := s.uc.GetPayments(ctx)

	if err != nil {

		s.logger.Error(
			"failed to get payments",
			"error",
			err,
		)

		return nil, err
	}

	return mapper.ToGetPaymentsResponse(
		payments,
	), nil
}

func (s *PaymentsService) UpdatePayment(
	ctx context.Context,
	req *pb.UpdatePaymentRequest,
) (*pb.UpdatePaymentReply, error) {

	if err := s.uc.UpdatePayment(
		ctx,
		req,
	); err != nil {

		s.logger.Error(
			"failed to update payment",
			"payment_id",
			req.PaymentId,
			"error",
			err,
		)

		return nil, err
	}

	return &pb.UpdatePaymentReply{
		Message: "Payment updated successfully",
	}, nil
}

func (s *PaymentsService) DeletePayment(
	ctx context.Context,
	req *pb.DeletePaymentRequest,
) (*pb.DeletePaymentReply, error) {

	if err := s.uc.DeletePayment(
		ctx,
		req.PaymentId,
	); err != nil {

		s.logger.Error(
			"failed to delete payment",
			"payment_id",
			req.PaymentId,
			"error",
			err,
		)

		return nil, err
	}

	return &pb.DeletePaymentReply{
		Message: "Payment deleted successfully",
	}, nil
}