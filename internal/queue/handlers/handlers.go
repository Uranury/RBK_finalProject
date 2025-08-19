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
	h.logger.Info("starting to handle send-invoice task", "payload_size", len(t.Payload()))

	var payload jobs.SendInvoicePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		h.logger.Error("failed to unmarshal SendInvoice payload", "err", err)
		return err
	}

	h.logger.Info("unmarshalled payload successfully", "order_id", payload.OrderID, "order_item_id", payload.OrderItemID, "to_email", payload.ToEmail)

	pdfBytes, err := h.InvoiceService.GenerateInvoicePDF(ctx, payload.OrderID, payload.OrderItemID)
	if err != nil {
		h.logger.Error("failed to generate PDF", "err", err)
		return err
	}

	h.logger.Info("PDF generated successfully", "pdf_size", len(pdfBytes))

	if err := h.EmailService.SendInvoice(payload.ToEmail, pdfBytes); err != nil {
		h.logger.Error("failed to send invoice email", "to", payload.ToEmail, "err", err)
		return err
	}

	h.logger.Info("send-invoice task completed successfully", "to", payload.ToEmail)
	return nil
}
