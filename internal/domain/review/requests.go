package review

import "github.com/google/uuid"

type CreateReviewRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Rating    int32     `json:"rating" validate:"required,gte=1,lte=5"`
	Comment   string    `json:"comment" validate:"required,min=2"`
}

type UpdateReviewRequest struct {
	Rating  *int32  `json:"rating,omitempty" validate:"omitempty,gte=1,lte=5"`
	Comment *string `json:"comment,omitempty" validate:"omitempty,min=2"`
}