//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	_ Status = iota
	PendingReview
	ChangeRequested
	Approved
	PaymentFinished
	Rejected
)

func (s Status) String() string {
	switch s {
	case PendingReview:
		return "pending_review"
	case ChangeRequested:
		return "change_requested"
	case Approved:
		return "approved"
	case PaymentFinished:
		return "payment_finished"
	case Rejected:
		return "rejected"
	default:
		// FIXME: ここerrorにしたい
		return ""
	}
}

// dbにstringいれる今の実装だとMarshalJson入らなそう。
func (s Status) MarshalJSON() ([]byte, error) {
	switch s {
	case PendingReview:
		return json.Marshal("pending_review")
	case ChangeRequested:
		return json.Marshal("change_requested")
	case Approved:
		return json.Marshal("approved")
	case PaymentFinished:
		return json.Marshal("payment_finished")
	case Rejected:
		return json.Marshal("rejected")
	default:
		return nil, fmt.Errorf("invalid status: %d", s)
	}
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("data should be a string, got %s", data)
	}

	var st Status
	switch str {
	case "pending_review":
		st = PendingReview
	case "change_requested":
		st = ChangeRequested
	case "approved":
		st = Approved
	case "payment_finished":
		st = PaymentFinished
	case "rejected":
		st = Rejected
	default:
		return fmt.Errorf("invalid Status %s", str)
	}
	*s = st
	return nil
}

type ApplicationStatusRepository interface {
	CreateStatus(
		ctx context.Context, applicationID uuid.UUID, userID uuid.UUID, status Status,
	) (*ApplicationStatus, error)
}

type ApplicationStatus struct {
	ID        uuid.UUID
	CreatedBy uuid.UUID
	Status    Status
	CreatedAt time.Time
}
