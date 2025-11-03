package idgen

import "github.com/google/uuid"

// Generic ID generation
func NewUUID() string {
    return uuid.New().String()
}

func IsValidUUID(id string) bool {
    _, err := uuid.Parse(id)
    return err == nil
}