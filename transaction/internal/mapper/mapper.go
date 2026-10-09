package mapper

import (
	"strconv"
	"time"

	pb "transaction/api/transaction/v1"
	"transaction/internal/model"
)

func ToTransaction(req *pb.CreateTransactionRequest) *model.Transaction {
	return &model.Transaction{
		PaymentID:       req.PaymentId,
		AccountID:       req.AccountId,
		TransactionType: req.TransactionType,
		Currency:        req.Currency,
		Description:     req.Description,
		ReferenceNumber: req.ReferenceNumber,
		MerchantInfo:    req.MerchantInfo,
	}
}

func ToGetTransactionResponse(
	transaction *model.Transaction,
) *pb.GetTransactionReply {

	response := &pb.GetTransactionReply{
		TransactionId: transaction.TransactionID,
		PaymentId:     transaction.PaymentID,
		AccountId:     transaction.AccountID,
		TransactionType: transaction.TransactionType,

		Amount: strconv.FormatFloat(
			transaction.Amount,
			'f',
			4,
			64,
		),

		Currency:       transaction.Currency,
		Status:         transaction.Status,
		Description:    transaction.Description,
		ReferenceNumber: transaction.ReferenceNumber,
		MerchantInfo:   transaction.MerchantInfo,

		TransactionDate: transaction.TransactionDate.Format(
			time.RFC3339,
		),

		CreatedAt: transaction.CreatedAt.Format(
			time.RFC3339,
		),
	}

	if transaction.RelatedTransactionID > 0 {
		response.RelatedTransactionId =
			transaction.RelatedTransactionID
	}

	if !transaction.PostingDate.IsZero() {
		response.PostingDate =
			transaction.PostingDate.Format("2006-01-02")
	}

	if !transaction.PostedAt.IsZero() {
		response.PostedAt =
			transaction.PostedAt.Format(time.RFC3339)
	}

	if !transaction.ArchivedAt.IsZero() {
		response.ArchivedAt =
			transaction.ArchivedAt.Format(time.RFC3339)
	}

	return response
}

func ToGetTransactionsResponse(
	transactions []*model.Transaction,
) *pb.GetTransactionsReply {

	result := make(
		[]*pb.GetTransactionReply,
		0,
		len(transactions),
	)

	for _, transaction := range transactions {
		result = append(
			result,
			ToGetTransactionResponse(transaction),
		)
	}

	return &pb.GetTransactionsReply{
		Transactions: result,
	}
}