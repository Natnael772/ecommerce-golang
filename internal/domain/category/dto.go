package category

import "github.com/google/uuid"

type CreateCategoryRequest struct {
	Name     string     `json:"name" validate:"required,min=2,max=100"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
}

type UpdateCategoryRequest struct {
	Name     *string     `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	ParentID *uuid.UUID  `json:"parent_id,omitempty"`
}

