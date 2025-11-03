package payment

import (
	"github.com/go-chi/chi/v5"
)
type PaymentRoutes struct{}

func Routes(svc PaymentService) chi.Router {
	h := NewPaymentHandler(svc)
	r := chi.NewRouter()
	
	r.Get("/order/{orderID}", h.GetPaymentByOrderID)

	r.Post("/webhook/{provider}", h.HandleWebhook)

	return r
}