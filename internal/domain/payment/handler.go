package payment

import (
	"ecommerce-app/internal/pkg/logger"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type PaymentHandler struct {
	svc PaymentService
}

func NewPaymentHandler(svc PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}


func (h *PaymentHandler) GetPaymentByOrderID(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")
	// implement using repo
	w.Write([]byte("TODO: Get payment by order ID: " + orderID))
}

func (h *PaymentHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("Failed to read webhook payload: %v", err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	sigHeader := r.Header.Get("Stripe-Signature")
	logger.Info("stripe sigHeader: %s", sigHeader)

	if err := h.svc.HandleWebhook(r.Context(), provider, payload, sigHeader); err != nil {
		logger.Error("Failed to handle webhook: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook processed successfully"))
}

