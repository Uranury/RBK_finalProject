package jobs

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	SendEmail   = "send:email"
	GeneratePDF = "generate:pdf"
)

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type PDFPayload struct {
	OrderID     uuid.UUID `json:"order_id" db:"order_id"`
	OrderItemID uuid.UUID `json:"order_item_id" db:"order_item_id"`
}

func NewSendEmailTask(to, subject, body string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(SendEmail, payload), nil
}

func NewGeneratePDFTask(orderID uuid.UUID, orderItemID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(PDFPayload{
		OrderID:     orderID,
		OrderItemID: orderItemID,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(GeneratePDF, payload), nil
}
