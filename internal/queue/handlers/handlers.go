package handlers

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/Uranury/RBK_finalProject/internal/queue/jobs"
	"github.com/Uranury/RBK_finalProject/internal/services"
	"github.com/hibiken/asynq"
)

type WorkerHandler struct {
	EmailService   *services.EmailService
	InvoiceService *services.InvoiceService
	logger         *slog.Logger
}

func NewWorkerHandler(emailService *services.EmailService, invoiceService *services.InvoiceService, logger *slog.Logger) *WorkerHandler {
	return &WorkerHandler{EmailService: emailService, InvoiceService: invoiceService, logger: logger}
}

func (h *WorkerHandler) HandleSendInvoiceTask(ctx context.Context, t *asynq.Task) error {
	var payload jobs.SendInvoicePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		h.logger.Warn("failed to unmarshal SendInvoice payload", "err", err)
		return err
	}

	pdfBytes, err := h.InvoiceService.GenerateInvoicePDF(ctx, payload.OrderID, payload.OrderItemID)
	if err != nil {
		h.logger.Warn("failed to generate PDF", "err", err)
		return err
	}

	if err := h.EmailService.SendInvoice(payload.ToEmail, pdfBytes); err != nil {
		h.logger.Warn("failed to send invoice email", "to", payload.ToEmail, "err", err)
		return err
	}

	return nil
}
