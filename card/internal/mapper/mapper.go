package mapper

import (
	pb "card/api/card/v1"
	"card/internal/model"
	"strconv"
	"time"
)

func ToCard(req *pb.CreateCardRequest) *model.Card {
	return &model.Card{
		AccountID: req.AccountId,
		CardType:  req.CardType,
	}
}

func ToGetCardResponse(card *model.Card) *pb.GetCardReply {
	return &pb.GetCardReply{
		CardId:         card.CardID,
		AccountId:      card.AccountID,
		CardNumberHash: card.CardNumberHash,
		CardType:       card.CardType,
		CardStatus:     card.CardStatus,
		ExpiryDate:     card.ExpiryDate.Format("2006-01-02"),
		DailyLimit:     strconv.FormatFloat(card.DailyLimit, 'f', 4, 64),
		MonthlyLimit:   strconv.FormatFloat(card.MonthlyLimit, 'f', 4, 64),
		PinAttempts:    card.PinAttempts,
		CreatedAt:      card.CreatedAt.Format(time.RFC3339),
		LastUsedAt:     card.LastUsedAt.Format(time.RFC3339),
	}
}

func ToUpdateCard(req *pb.UpdateCardRequest,card *model.Card,
) {
	if req.CardType != nil {
		card.CardType = *req.CardType
	}

	if req.CardStatus != nil {
		card.CardStatus = *req.CardStatus
	}

	if req.ExpiryDate != nil {
		if value, err := time.Parse("2006-01-02", *req.ExpiryDate); err == nil {
			card.ExpiryDate = value
		}
	}

	if req.DailyLimit != nil {
		if value, err := strconv.ParseFloat(*req.DailyLimit, 64); err == nil {
			card.DailyLimit = value
		}
	}

	if req.MonthlyLimit != nil {
		if value, err := strconv.ParseFloat(*req.MonthlyLimit, 64); err == nil {
			card.MonthlyLimit = value
		}
	}
}

func ToGetCardsResponse(cards []*model.Card) *pb.GetCardsReply {
	result := make([]*pb.GetCardReply, 0, len(cards))

	for _, card := range cards {
		result = append(result, ToGetCardResponse(card))
	}

	return &pb.GetCardsReply{
		Cards: result,
	}
}