package auth

// --- Request structs (DTO) for auth operations

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// --- Response structs for auth operations

type UserWithToken struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}