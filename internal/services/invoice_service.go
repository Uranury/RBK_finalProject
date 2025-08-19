package services

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Uranury/RBK_finalProject/internal/repositories/order"
	"github.com/Uranury/RBK_finalProject/pkg/apperrors"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"log/slog"
)

type InvoiceService struct {
	orderRepo order.Repository
	logger    *slog.Logger
}

func NewInvoiceService(orderRepo order.Repository, logger *slog.Logger) *InvoiceService {
	return &InvoiceService{orderRepo: orderRepo, logger: logger}
}

func (s *InvoiceService) GenerateInvoicePDF(ctx context.Context, orderID, orderItemID uuid.UUID) ([]byte, error) {
	ord, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		s.logger.Warn("failed to retrieve order", "order_id", orderID, "err", err)
		return nil, apperrors.NewInternalError("failed to get order by id", err)
	}
	if ord == nil {
		return nil, apperrors.NewNotFoundError("order not found")
	}

	ordItem, err := s.orderRepo.GetOrderItemByID(ctx, orderItemID)
	if err != nil {
		s.logger.Warn("failed to retrieve order item", "order_item_id", orderItemID, "err", err)
		return nil, apperrors.NewInternalError("failed to get order item by id", err)
	}
	if ordItem == nil {
		return nil, apperrors.NewNotFoundError("order item not found")
	}

	// --- PDF generation ---
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(40, 10, "Invoice")
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 10,
		fmt.Sprintf("Order ID: %s\nUser ID: %s\nTotal Amount: %.2f\nStatus: %s\nCreated At: %s\n\nOrder Item ID: %s\nSkin ID: %s\nPrice: %.2f\nCreated At: %s",
			ord.ID.String(),
			ord.UserID.String(),
			ord.TotalAmount,
			string(ord.Status),
			ord.CreatedAt.Format("2006-01-02 15:04:05"),
			ordItem.ID.String(),
			ordItem.SkinID.String(),
			ordItem.Price,
			ordItem.CreatedAt.Format("2006-01-02 15:04:05"),
		),
		"", "", false,
	)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, apperrors.NewInternalError("failed to generate PDF", err)
	}

	return buf.Bytes(), nil
}
