package service

import (
	"context"
	"log/slog"

	pb "card/api/card/v1"
	"card/internal/biz"
	"card/internal/mapper"
)

type CardService struct {
	pb.UnimplementedCardServiceServer

	uc     *biz.CardUsecase
	logger *slog.Logger
}

func NewCardService(uc *biz.CardUsecase,logger *slog.Logger,) *CardService {
	return &CardService{
		uc:     uc,
		logger: logger,
	}
}

func (s *CardService) CreateCard(ctx context.Context,req *pb.CreateCardRequest,
) (*pb.CreateCardReply, error) {

	if err := s.uc.CreateCard(ctx, req); err != nil {
		s.logger.Error("failed to create card", "error", err)
		return nil, err
	}

	return &pb.CreateCardReply{
		Message: "Card created successfully",
	}, nil
}

func (s *CardService) GetCard(ctx context.Context,req *pb.GetCardRequest,) (*pb.GetCardReply, error) {

	card, err := s.uc.GetCard(ctx, req.CardId)
	if err != nil {
		s.logger.Error("failed to get card","card_id", req.CardId,"error", err,)
		return nil, err
	}

	return mapper.ToGetCardResponse(card), nil
}


func (s *CardService) GetCards(ctx context.Context,req *pb.GetCardsRequest) (*pb.GetCardsReply, error) {

	cards, err := s.uc.GetCards(ctx)
	if err != nil {
		s.logger.Error(
			"failed to get cards",
			"error", err,
		)
		return nil, err
	}

	return mapper.ToGetCardsResponse(cards), nil
}


func (s *CardService) UpdateCard(	ctx context.Context,req *pb.UpdateCardRequest,) (*pb.UpdateCardReply, error) {

	if err := s.uc.UpdateCard(ctx, req); err != nil {
		s.logger.Error(
			"failed to update card",
			"card_id", req.CardId,
			"error", err,
		)
		return nil, err
	}

	return &pb.UpdateCardReply{
		Message: "Card updated successfully",
	}, nil
}

func (s *CardService) DeleteCard(ctx context.Context,req *pb.DeleteCardRequest,) (*pb.DeleteCardReply, error) {

	if err := s.uc.DeleteCard(ctx, req.CardId); err != nil {
		s.logger.Error(
			"failed to delete card",
			"card_id", req.CardId,
			"error", err,
		)
		return nil, err
	}

	return &pb.DeleteCardReply{
		Message: "Card deleted successfully",
	}, nil
}