package user

// Request structs (DTO) for user operations

type RegisterUserRequest struct {
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8,max=64"`
    FirstName string `json:"first_name" validate:"required,min=2,max=50"`
    LastName  string `json:"last_name" validate:"required,min=2,max=50"`
    Role     string `json:"role" validate:"required,oneof=customer admin superadmin vendor support delivery"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	ID 	  string `json:"id" validate:"required,uuid4"`
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,max=50"`
}

// type ListUsersParams struct {
// 	Limit  int `json:"limit" validate:"omitempty,min=1,max=100"`
// 	Offset int `json:"offset" validate:"omitempty,min=0"`
// }



// type ListUsersParams struct {
//     Page    int    `json:"page" validate:"gte=1"`
//     PerPage int    `json:"per_page" validate:"gte=1,lte=100"`
//     Email   string `json:"email,omitempty" validate:"omitempty,email"`
// 	Query   string `json:"query,omitempty" validate:"omitempty,min=2"`
// }