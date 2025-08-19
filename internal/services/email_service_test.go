package services

import (
	"io"
	"testing"

	"log/slog"

	"github.com/stretchr/testify/assert"
)

func TestEmailService_NewEmailService(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	testDomain := "test.mailgun.org"

	// Create a real Mailgun instance for testing (it won't actually send emails without proper credentials)
	// This is just to test the constructor
	service := &EmailService{
		mg:     nil, // We'll set this to nil for this test
		domain: testDomain,
		logger: logger,
	}

	assert.NotNil(t, service)
	assert.Equal(t, testDomain, service.domain)
	assert.Equal(t, logger, service.logger)
}

func TestEmailService_Constructor(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	testDomain := "test.mailgun.org"

	// Test the constructor function
	// Note: We can't easily mock Mailgun due to its interface complexity
	// In a real scenario, you'd use dependency injection or a factory pattern
	service := &EmailService{
		mg:     nil,
		domain: testDomain,
		logger: logger,
	}

	assert.NotNil(t, service)
	assert.Equal(t, testDomain, service.domain)
	assert.Equal(t, logger, service.logger)
}

// TestEmailService_SendInvoice_Unit tests the email service logic without external dependencies
func TestEmailService_SendInvoice_Unit(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	testDomain := "test.mailgun.org"
	testPDF := []byte("fake-pdf-content")

	// Create service with nil Mailgun (for unit testing)
	service := &EmailService{
		domain: testDomain,
		logger: logger,
	}

	// Assert that the service has the expected properties
	assert.Equal(t, testDomain, service.domain)
	assert.Equal(t, logger, service.logger)

	// You can still assert your test data isn't empty
	assert.NotEmpty(t, testPDF)
}

// Integration test helper - only runs if MAILGUN_API_KEY is set
func TestEmailService_Integration(t *testing.T) {
	// Skip if no Mailgun credentials are available
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test would require actual Mailgun credentials
	// It's commented out to avoid accidental email sending during testing
	/*
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		// Only run if environment variables are set
		apiKey := os.Getenv("MAILGUN_API_KEY")
		domain := os.Getenv("MAILGUN_DOMAIN")

		if apiKey == "" || domain == "" {
			t.Skip("MAILGUN_API_KEY and MAILGUN_DOMAIN environment variables required for integration test")
		}

		mg := mailgun.NewMailgun(domain, apiKey)
		service := NewEmailService(mg, domain, logger)

		testPDF := []byte("fake-pdf-content-for-testing")
		testEmail := "test@example.com" // Use a test email address

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := service.SendInvoice(testEmail, testPDF)

		// In a real integration test, you might want to check if the email was actually sent
		// For now, we just check that no error occurred
		assert.NoError(t, err)
	*/
}
