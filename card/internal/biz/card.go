package biz

import (
	"context"
	"errors"
	"strings"
	"time"

	pb "card/api/card/v1"
	"card/internal/mapper"
	"card/internal/model"

	"github.com/bwmarrin/snowflake"
	"log/slog"
)

type CardRepo interface {
	Create(ctx context.Context, card *model.Card) error
	GetByID(ctx context.Context, id int64) (*model.Card, error)
	GetAll(ctx context.Context) ([]*model.Card, error)
	Update(ctx context.Context, card *model.Card) error
	Delete(ctx context.Context, id int64) error
}

type CardUsecase struct {
	repo   CardRepo
	node   *snowflake.Node
	logger *slog.Logger
}

func NewCardUsecase(repo CardRepo, logger *slog.Logger) (*CardUsecase, error) {
	node, err := snowflake.NewNode(2)
	if err != nil {
		return nil, err
	}

	return &CardUsecase{
		repo:   repo,
		node:   node,
		logger: logger,
	}, nil
}

func (uc *CardUsecase) CreateCard(ctx context.Context, req *pb.CreateCardRequest) error {
	if req.AccountId <= 0 {
		return errors.New("account id is required")
	}

	cardType := strings.ToLower(strings.TrimSpace(req.CardType))

	if cardType != "debit" &&
		cardType != "credit" &&
		cardType != "prepaid" {
		return errors.New("invalid card type")
	}

	card := mapper.ToCard(req)

	card.CardID = uc.node.Generate().Int64()
	card.CardType = cardType
	card.CardStatus = "active"
	card.ExpiryDate = time.Now().AddDate(3, 0, 0)
	card.DailyLimit = 0
	card.MonthlyLimit = 0
	card.PinAttempts = 0
	card.CreatedAt = time.Now()

	if err := uc.repo.Create(ctx, card); err != nil {
		uc.logger.Error("failed to create card", "error", err)
		return err
	}

	uc.logger.Info("card created", "card_id", card.CardID)

	return nil
}

func (uc *CardUsecase) GetCard(ctx context.Context, id int64) (*model.Card, error) {
	if id <= 0 {
		return nil, errors.New("card id is required")
	}

	return uc.repo.GetByID(ctx, id)
}

func (uc *CardUsecase) GetCards(ctx context.Context) ([]*model.Card, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *CardUsecase) UpdateCard(
	ctx context.Context,
	req *pb.UpdateCardRequest,
) error {

	if req.CardId <= 0 {
		return errors.New("card id is required")
	}

	card, err := uc.repo.GetByID(ctx, req.CardId)
	if err != nil {
		return err
	}

	if req.CardType != nil {
		cardType := strings.ToLower(strings.TrimSpace(*req.CardType))

		if cardType != "debit" &&
			cardType != "credit" &&
			cardType != "prepaid" {
			return errors.New("invalid card type")
		}

		card.CardType = cardType
	}

	if req.CardStatus != nil {
		status := strings.ToLower(strings.TrimSpace(*req.CardStatus))

		if status != "active" &&
			status != "blocked" &&
			status != "expired" {
			return errors.New("invalid card status")
		}

		card.CardStatus = status
	}

	if req.ExpiryDate != nil {
		expiryDate, err := time.Parse(
			"2006-01-02",
			*req.ExpiryDate,
		)
		if err != nil {
			return errors.New("invalid expiry date")
		}

		card.ExpiryDate = expiryDate
	}

	if req.DailyLimit != nil {
	}

	if req.MonthlyLimit != nil {
		
	}

	return uc.repo.Update(ctx, card)
}

func (uc *CardUsecase) DeleteCard(
	ctx context.Context,
	id int64,
) error {

	if id <= 0 {
		return errors.New("card id is required")
	}

	card, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if card.CardStatus == "blocked" {
		return errors.New("blocked card cannot be deleted")
	}

	return uc.repo.Delete(ctx, id)
}