package router

import "github.com/google/uuid"

type Transaction struct {
	Amount int          `json:"amount"`
	Target string       `json:"target"`
	Tags   []*uuid.UUID `json:"tags"`
	Group  *uuid.UUID   `json:"group"`
}
