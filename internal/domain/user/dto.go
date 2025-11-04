package user

// --- Request Dto ---
type RegisterUserRequest struct {
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8,max=64"`
    FirstName string `json:"first_name" validate:"required,min=2,max=50"`
    LastName  string `json:"last_name" validate:"required,min=2,max=50"`
    Role     string `json:"role" validate:"required,oneof=customer admin superadmin vendor support delivery"`
}

type UpdateUserRequest struct {
	ID 	  string `json:"id" validate:"required,uuid4"`
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,max=50"`
}
