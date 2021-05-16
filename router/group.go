package router

import "github.com/google/uuid"

type group struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Budget      int          `json:"budget"`
	Owner       []*uuid.UUID `json:"owner"`
}
