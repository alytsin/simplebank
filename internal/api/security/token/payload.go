package token

import (
	"github.com/google/uuid"
)

type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func NewPayload(username string) *Payload {
	return &Payload{
		ID:       uuid.New(),
		Username: username,
	}
}
