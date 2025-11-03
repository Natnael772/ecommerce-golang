package user

import (
	"ecommerce-app/internal/pkg/validator"

	"github.com/go-chi/chi/v5"
)

func Routes(svc Service) chi.Router {
	h := NewHandler(svc)
	r := chi.NewRouter()

	r.With(validator.Validate[RegisterUserRequest]()).Post("/register", h.Register)

	r.With(validator.Validate[LoginRequest]()).Post("/auth/login", h.Login)
	r.Get("/{id}", h.GetUser)
	r.Get("/", h.ListUsers)
	

	return r
}
