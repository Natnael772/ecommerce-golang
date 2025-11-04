package stripe

import (
	"context"
	"ecommerce-app/configs"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/paymentintent"
	"github.com/stripe/stripe-go/v83/webhook"
)

const ProviderName = "stripe"

func NewStripeProvider() *StripeProvider {
	cfg := configs.Load()
	apiKey := cfg.StripeAPIKey
	webhookSecret := cfg.StripeWebhookSecret

	stripe.Key = apiKey
	return &StripeProvider{webhookSecret: webhookSecret}
}

// CreatePaymentIntent creates a Stripe PaymentIntent
func (s *StripeProvider) CreatePaymentIntent(ctx context.Context, amountCents int64, currency string, metadata map[string]string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(currency),
		Metadata: metadata,
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		return nil, err
	}

	return intent, nil
}

func (s *StripeProvider) HandleWebhook(payload []byte, sigHeader string) (*ProviderWebhookEvent, error) {
	event, err := s.VerifyWebhookSignature(payload, sigHeader)
	if err != nil {
		return nil, fmt.Errorf("invalid webhook signature: %w", err)
	}

	// Map PaymentIntent event types to internal statuses
	piEvents := map[string]string{
		"payment_intent.created":        "INITIATED",
		"payment_intent.succeeded":      "SUCCESS",
		"payment_intent.payment_failed": "FAILED",
	}

	// Map Charge event types to internal statuses
	chargeEvents := map[string]string{
		"charge.succeeded": "SUCCESS",
		"charge.failed":    "FAILED",
		"charge.pending":   "INITIATED", // optional
		"charge.refunded":  "REFUNDED",
		"charge.captured":  "SUCCESS", // for delayed capture flows
		"charge.updated":   "UPDATED", // optional: for charge status updates
	}

	var status string
	if s, ok := piEvents[string(event.Type)]; ok {
		status = s
		// Parse PaymentIntent
		var intent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
			return nil, err
		}

		orderID, ok := intent.Metadata["order_id"]
		if !ok {
			return nil, errors.New("missing order_id in metadata")
		}

		failureReason := ""
		if status == "FAILED" && intent.LastPaymentError != nil {
			failureReason = intent.LastPaymentError.Msg
		}

		return &ProviderWebhookEvent{
			Provider:      ProviderName,
			ProviderTxnID: intent.ID,
			OrderID:       orderID,
			Status:        status,
			FailureReason: failureReason,
			RawEvent:      intent,
		}, nil
	} else if s, ok := chargeEvents[string(event.Type)]; ok {
		status = s
		// Parse Charge
		var ch stripe.Charge
		if err := json.Unmarshal(event.Data.Raw, &ch); err != nil {
			return nil, err
		}

		// Link to PaymentIntent for metadata if available
		var orderID string
		if ch.PaymentIntent != nil {
			// Retrieve metadata from PaymentIntent if needed
			pi, err := paymentintent.Get(ch.PaymentIntent.ID, nil)
			if err == nil {
				orderID = pi.Metadata["order_id"]
			}
		}

		// Special handling for charge.updated event
		if event.Type == "charge.updated" && ch.Status != "" {
			switch ch.Status {
			case "succeeded":
				status = "SUCCESS"
			case "failed":
				status = "FAILED"
			case "pending":
				status = "INITIATED"
			case "canceled":
				status = "CANCELLED"
			default:
				return nil, fmt.Errorf("unhandled charge status: %s", ch.Status)
			}
		}

		return &ProviderWebhookEvent{
			Provider:      ProviderName,
			ProviderTxnID: ch.ID,
			OrderID:       orderID,
			Status:        status,
			FailureReason: ch.FailureMessage,
			RawEvent:      ch,
		}, nil
	}

	// Catch-all for unhandled event types
	fmt.Printf("Stripe webhook ignored: %s\n", event.Type)
	return nil, nil
}

func (s *StripeProvider) VerifyWebhookSignature(payload []byte, sigHeader string) (stripe.Event, error) {
	event, err := webhook.ConstructEventWithOptions(
		payload,
		sigHeader,
		s.webhookSecret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
	if err != nil {
		return stripe.Event{}, fmt.Errorf("invalid webhook signature: %w", err)
	}
	return event, nil
}
