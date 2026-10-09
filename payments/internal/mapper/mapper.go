package mapper

import (
	"strconv"
	"time"

	pb "payments/api/payments/v2"
	"payments/internal/model"
)

func ToPayment(req *pb.CreatePaymentRequest) *model.Payment {
	return &model.Payment{
		FromAccountID:     req.FromAccountId,
		ToAccountID:       req.ToAccountId,
		PaymentType:       req.PaymentType,
		Currency:          req.Currency,
		PaymentMethod:     req.PaymentMethod,
		ReferenceNumber:   req.ReferenceNumber,
		ExternalReference: req.ExternalReference,
	}
}

func ToGetPaymentResponse(
	payment *model.Payment,
) *pb.GetPaymentReply {

	response := &pb.GetPaymentReply{
		PaymentId:     payment.PaymentID,
		FromAccountId: payment.FromAccountID,
		ToAccountId:   payment.ToAccountID,
		PaymentType:   payment.PaymentType,

		Amount: strconv.FormatFloat(payment.Amount,'f',4,64,),

		Currency:          payment.Currency,
		Status:            payment.Status,
		PaymentMethod:     payment.PaymentMethod,
		ReferenceNumber:   payment.ReferenceNumber,
		ExternalReference: payment.ExternalReference,

		CreatedAt: payment.CreatedAt.Format(
			time.RFC3339,
		),
	}

	if !payment.ScheduledAt.IsZero() {
		response.ScheduledAt =
			payment.ScheduledAt.Format(time.RFC3339)
	}

	if !payment.ProcessedAt.IsZero() {
		response.ProcessedAt =
			payment.ProcessedAt.Format(time.RFC3339)
	}

	return response
}

func ToGetPaymentsResponse(
	payments []*model.Payment,
) *pb.GetPaymentsReply {

	result := make(
		[]*pb.GetPaymentReply,
		0,
		len(payments),
	)

	for _, payment := range payments {
		result = append(
			result,
			ToGetPaymentResponse(payment),
		)
	}

	return &pb.GetPaymentsReply{
		Payments: result,
	}
}

func ToUpdatePayment(
	req *pb.UpdatePaymentRequest,
	payment *model.Payment,
) {
	if req.PaymentType != nil {
		payment.PaymentType = *req.PaymentType
	}

	if req.Amount != nil {
		if value, err := strconv.ParseFloat(*req.Amount, 64); err == nil {
			payment.Amount = value
		}
	}

	if req.Currency != nil {
		payment.Currency = *req.Currency
	}

	if req.Status != nil {
		payment.Status = *req.Status
	}

	if req.PaymentMethod != nil {
		payment.PaymentMethod = *req.PaymentMethod
	}

	if req.ExternalReference != nil {
		payment.ExternalReference = *req.ExternalReference
	}

	if req.ScheduledAt != nil {
		if value, err := time.Parse(
			time.RFC3339,
			*req.ScheduledAt,
		); err == nil {
			payment.ScheduledAt = value
		}
	}
}