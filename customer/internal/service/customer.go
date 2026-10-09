package service

import (
	"context"
	"log/slog"

	pb "customer/api/customer/v1"
	"customer/internal/biz"
	"customer/internal/mapper"
)

type CustomerService struct {
	pb.UnimplementedCustomerServiceServer
	uc     *biz.CustomerUsecase
	logger *slog.Logger
}

func NewCustomerService(uc *biz.CustomerUsecase, logger *slog.Logger) *CustomerService {
	return &CustomerService{
		uc:     uc,
		logger: logger,
	}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerReply, error) {
	customer, err := mapper.ToCustomer(req)
	if err != nil {
		s.logger.Error("failed to map customer request", "error", err)
		return nil, err
	}

	if err := s.uc.CreateCustomer(ctx, customer); err != nil {
		s.logger.Error("failed to create customer", "error", err)
		return nil, err
	}

	s.logger.Info("customer created successfully", "email", customer.Email)

	return &pb.CreateCustomerReply{
		Message: "Customer created successfully",
	}, nil
}



func (s *CustomerService) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.GetCustomerReply, error) {
	customer, err := s.uc.GetCustomer(ctx, req.CustomerId)
	if err != nil {
		s.logger.Error("failed to get customer", "customer_id", req.CustomerId, "error", err)
		return nil, err
	}

	return mapper.ToGetCustomerResponse(customer), nil
}

func (s *CustomerService) GetCustomers(ctx context.Context, req *pb.GetCustomersRequest) (*pb.GetCustomersReply, error) {
	customers, err := s.uc.GetCustomers(ctx)
	if err != nil {
		s.logger.Error("failed to get customers", "error", err)
		return nil, err
	}

	return mapper.ToGetCustomersResponse(customers), nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.UpdateCustomerReply, error) {
	customer, err := s.uc.GetCustomer(ctx, req.CustomerId)
	if err != nil {
		return nil, err
	}

	if err := mapper.ToUpdateCustomer(req, customer); err != nil {
		return nil, err
	}

	if err := s.uc.UpdateCustomer(ctx, customer); err != nil {
		return nil, err
	}

	return &pb.UpdateCustomerReply{
		Message: "Customer updated successfully",
	}, nil
}
func (s *CustomerService) DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest) (*pb.DeleteCustomerReply, error) {
	if err := s.uc.DeleteCustomer(ctx, req.CustomerId); err != nil {
		s.logger.Error("failed to delete customer", "customer_id", req.CustomerId, "error", err)
		return nil, err
	}

	return &pb.DeleteCustomerReply{
		Message: "Customer deleted successfully",
	}, nil
}