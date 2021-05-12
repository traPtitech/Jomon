package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	_ Status = iota
	Submitted
	FixRequired
	Accepted
	Completed
	Rejected
)

func (s Status) String() string {
	switch s {
	case Submitted:
		return "submitted"
	case FixRequired:
		return "fix_required"
	case Accepted:
		return "accepted"
	case Completed:
		return "completed"
	case Rejected:
		return "rejected"
	default:
		return ""
	}
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("data should be a string, got %s", data)
	}

	var st Status
	switch str {
	case "submitted":
		st = Submitted
	case "fix_required":
		st = FixRequired
	case "accepted":
		st = Accepted
	case "completed":
		st = Completed
	case "rejected":
		st = Rejected
	default:
		return fmt.Errorf("invalid Status %s", str)
	}
	*s = st
	return nil
}

type RequestStatus struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	RequestID uuid.UUID `gorm:"type:char(36);not null;index"`
	Request   *Request  `gorm:"foeignKey:RequestID"`
	CreatedBy string    `gorm:"type:varchar(32);not null"`
	Status    Status    `gorm:"type:enum('submitted','fix_required','accepted','completed','rejected');not null"`
	Reason    string    `gorm:"type:text;not null"`
	CreatedAt time.Time
}
