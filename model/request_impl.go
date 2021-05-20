package model

import (
	"github.com/traPtitech/Jomon/ent"
)

func ConvertEntRequestToModelRequest(request *ent.Request) *Request {
	return &Request{
		ID:        request.ID,
		Amount:    request.Amount,
		CreatedAt: request.CreatedAt,
	}
}
