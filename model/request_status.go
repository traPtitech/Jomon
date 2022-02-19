//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

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
//multipart/form-dataじゃないから使わなそう
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

func ConvertStrStatusToStatus(str string) (Status, error) {
	var status Status
	switch str {
	case "submitted":
		status = Submitted
	case "fix_required":
		status = FixRequired
	case "accepted":
		status = Accepted
	case "completed":
		status = Completed
	case "rejected":
		status = Rejected
	default:
		return Status(0), fmt.Errorf("invalid Status %s", str)
	}
	return status, nil
}
type RequestStatusRepository interface {
	CreateStatus(ctx context.Context, requestID uuid.UUID, userID uuid.UUID, status Status)
}

type RequestStatus struct {
	ID        uuid.UUID
	Status    string
	CreatedAt time.Time
}
