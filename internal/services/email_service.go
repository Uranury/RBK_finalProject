package services

import (
	"context"
	"github.com/mailgun/mailgun-go/v4"
	"log/slog"
	"time"
)

type EmailService struct {
	mg     mailgun.Mailgun
	domain string
	logger *slog.Logger
}

func NewEmailService(mg mailgun.Mailgun, domain string, logger *slog.Logger) *EmailService {
	return &EmailService{mg: mg, domain: domain, logger: logger}
}

func (s *EmailService) SendInvoice(to string, pdf []byte) error {
	msg := mailgun.NewMessage(
		"noreply@"+s.domain,
		"Your Invoice",
		"Hello, please find attached your invoice.",
		to,
	)

	msg.AddBufferAttachment("invoice.pdf", pdf)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, _, err := s.mg.Send(ctx, msg)
	if err != nil {
		s.logger.Warn("failed to send email", "to", to, "err", err)
		return err
	}

	return nil
}
