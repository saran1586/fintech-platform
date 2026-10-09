package mapper

import (
	"time"

	pb "customer/api/customer/v1"
	"customer/internal/model"
	"errors"
	
)

func ToCustomer(req *pb.CreateCustomerRequest) (*model.Customer, error) {
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, err
	}

	return &model.Customer{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Phone:       req.Phone,
		DateOfBirth: dob,
	}, nil
}

func ToGetCustomerResponse(customer *model.Customer) *pb.GetCustomerReply {
	return &pb.GetCustomerReply{
		CustomerId:     customer.CustomerID,
		CustomerNumber: customer.CustomerNumber,
		FirstName:      customer.FirstName,
		LastName:       customer.LastName,
		Email:          customer.Email,
		Phone:          customer.Phone,
		DateOfBirth:    customer.DateOfBirth.Format("2006-01-02"),
		Status:         customer.Status,
		KycStatus:      customer.KycStatus,
	}
}

func ToGetCustomersResponse(customers []*model.Customer) *pb.GetCustomersReply {
	result := make([]*pb.Customer, 0, len(customers))

	for _, customer := range customers {
		result = append(result, &pb.Customer{
			CustomerId:     customer.CustomerID,
			CustomerNumber: customer.CustomerNumber,
			FirstName:      customer.FirstName,
			LastName:       customer.LastName,
			Email:          customer.Email,
			Phone:          customer.Phone,
			DateOfBirth:    customer.DateOfBirth.Format("2006-01-02"),
			Status:         customer.Status,
			KycStatus:      customer.KycStatus,
		})
	}

	return &pb.GetCustomersReply{
		Customers: result,
	}
}

func ToUpdateCustomer(req *pb.UpdateCustomerRequest, customer *model.Customer) error {
	if req.FirstName != nil {
		customer.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		customer.LastName = *req.LastName
	}

	if req.Email != nil {
		customer.Email = *req.Email
	}

	if req.Phone != nil {
		customer.Phone = *req.Phone
	}

	if req.DateOfBirth != nil {
		dateOfBirth, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			return errors.New("invalid date of birth format")
		}

		customer.DateOfBirth = dateOfBirth
	}

	return nil
}