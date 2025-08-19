package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/mailgun/mailgun-go/v4"
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
	s.logger.Info("attempting to send invoice email", "to", to, "domain", s.domain, "pdf_size", len(pdf))

	msg := mailgun.NewMessage(
		"noreply@"+s.domain,
		"Your Invoice",
		"Hello, please find attached your invoice.",
		to,
	)

	msg.AddBufferAttachment("invoice.pdf", pdf)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.logger.Debug("sending email via Mailgun", "from", "noreply@"+s.domain, "to", to)

	_, _, err := s.mg.Send(ctx, msg)
	if err != nil {
		s.logger.Error("failed to send email", "to", to, "err", err)
		return err
	}

	s.logger.Info("email sent successfully", "to", to)
	return nil
}
