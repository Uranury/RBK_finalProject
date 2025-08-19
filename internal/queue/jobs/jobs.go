package jobs

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	SendInvoice = "invoice:send"
)

type SendInvoicePayload struct {
	OrderID     uuid.UUID `json:"order_id"`
	OrderItemID uuid.UUID `json:"order_item_id"`
	ToEmail     string    `json:"to_email"`
}

func NewSendInvoiceTask(orderID uuid.UUID, orderItemID uuid.UUID, toEmail string) (*asynq.Task, error) {
	payload, err := json.Marshal(SendInvoicePayload{
		OrderID:     orderID,
		OrderItemID: orderItemID,
		ToEmail:     toEmail,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(SendInvoice, payload), nil
}
