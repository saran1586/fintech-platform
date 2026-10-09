package biz
import (
	"context"
	"errors"
	"log/slog"
	"net/mail"
	"strings"
	"time"
	"unicode"
   
	"customer/internal/model"

	"github.com/bwmarrin/snowflake"
)

type CustomerRepo interface {
	Create(ctx context.Context, customer *model.Customer) error
	GetByID(ctx context.Context, id int64) (*model.Customer, error)
	GetAll(ctx context.Context) ([]*model.Customer, error)
	Update(ctx context.Context, customer *model.Customer) error
	
}


type CustomerUsecase struct {
	repo   CustomerRepo
	node   *snowflake.Node
	logger *slog.Logger
}

func NewCustomerUsecase(repo CustomerRepo, logger *slog.Logger) (*CustomerUsecase, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}

	return &CustomerUsecase{
		repo:   repo,
		node:   node,
		logger: logger,
	}, nil
}

func (uc *CustomerUsecase) CreateCustomer(ctx context.Context, customer *model.Customer) error {
	if err := validateCustomer(customer); err != nil {
		return err
	}

	now := time.Now()

	customer.CustomerID = uc.node.Generate().Int64()
	customer.CustomerNumber = "CUS" + now.Format("20060102150405")
	customer.Status = "active"
	customer.KycStatus = "pending"
	customer.CreatedAt = now
	customer.UpdatedAt = now

	if err := uc.repo.Create(ctx, customer); err != nil {
		uc.logger.Error("failed to create customer", "error", err)
		return err
	}

	uc.logger.Info("customer created", "customer_id", customer.CustomerID)

	return nil
}

func validateCustomer(customer *model.Customer) error {
	if !isValidName(customer.FirstName) {
		return errors.New("first name must contain only letters")
	}

	if !isValidName(customer.LastName) {
		return errors.New("last name must contain only letters")
	}

	if !isValidEmail(customer.Email) {
		return errors.New("invalid email format")
	}

	if !isValidPhone(customer.Phone) {
		return errors.New("phone number must be 10 digits")
	}

	if customer.DateOfBirth.IsZero() {
		return errors.New("date of birth is required")
	}

	return nil
}
func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)

	if email == "" {
		return false
	}

	_, err := mail.ParseAddress(email)
	return err == nil
}
func isValidPhone(phone string) bool {
	for _, char := range phone {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}
func isValidName(name string) bool {
	name = strings.TrimSpace(name)

	if name == "" {
		return false
	}

	for _, char := range name {
		if !unicode.IsLetter(char) {
			return false
		}
	}

	return true
}


func (uc *CustomerUsecase) GetCustomer(ctx context.Context, id int64) (*model.Customer, error) {
	if id <= 0 {
		return nil, errors.New("invalid customer id")
	}

	customer, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("failed to get customer", "customer_id", id, "error", err)
		return nil, err
	}

	return customer, nil
}


func (uc *CustomerUsecase) GetCustomers(ctx context.Context) ([]*model.Customer, error) {
	customers, err := uc.repo.GetAll(ctx)
	if err != nil {
		uc.logger.Error("failed to get customers", "error", err)
		return nil, err
	}

	return customers, nil
}

func (uc *CustomerUsecase) UpdateCustomer(ctx context.Context, customer *model.Customer) error {
	if err := validateCustomer(customer); err != nil {
		return err
	}

	existing, err := uc.repo.GetByID(ctx, customer.CustomerID)
	if err != nil {
		uc.logger.Error("failed to get customer", "customer_id", customer.CustomerID, "error", err)
		return err
	}

	existing.FirstName = customer.FirstName
	existing.LastName = customer.LastName
	existing.Email = customer.Email
	existing.Phone = customer.Phone
	existing.DateOfBirth = customer.DateOfBirth
	existing.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, existing); err != nil {
		uc.logger.Error("failed to update customer", "customer_id", customer.CustomerID, "error", err)
		return err
	}

	uc.logger.Info("customer updated", "customer_id", customer.CustomerID)

	return nil
}

func (uc *CustomerUsecase) DeleteCustomer(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid customer id")
	}

	customer, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("failed to get customer", "customer_id", id, "error", err)
		return err
	}

	customer.Status = "closed"
	customer.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, customer); err != nil {
		uc.logger.Error("failed to close customer", "customer_id", id, "error", err)
		return err
	}

	uc.logger.Info("customer closed", "customer_id", id)

	return nil
}