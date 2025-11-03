package address

import (
	"ecommerce-app/internal/pkg/validator"

	"github.com/go-chi/chi/v5"
)

func Routes(svc Service) chi.Router {
	h := NewHandler(svc)
	r := chi.NewRouter()

	r.With(validator.Validate[CreateAddressRequest]()).Post("/", h.CreateAddress)

	r.Get("/", h.getAddressesByUser)
	
	r.Get("/{id}", h.GetAddressByID)

	r.With(validator.Validate[UpdateAddressRequest]()).Put("/{id}", h.UpdateAddress)
	r.Delete("/{id}", h.DeleteAddress)
	return r
}