//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/requeststatus"
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

//dbにstringいれる今の実装だとMarshalJson入らなそう。
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
	CreateStatus(ctx context.Context, requestID uuid.UUID, userID uuid.UUID, status Status) (*RequestStatus, error)
}

type RequestStatus struct {
	ID        uuid.UUID
	Status    string
	CreatedAt time.Time
}

func (repo *EntRepository) CreateStatus(ctx context.Context, requestID uuid.UUID, userID uuid.UUID, status Status) (*RequestStatus, error) {
	created, err := repo.client.RequestStatus.
		Create().
		SetStatus(requeststatus.Status(status.String())).
		SetCreatedAt(time.Now()).
		SetRequestID(requestID).
		SetUserID(userID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return convertEntRequestStatusToModelRequestStatus(created), nil
}

func convertEntRequestStatusToModelRequestStatus(requestStatus *ent.RequestStatus) *RequestStatus {
	if requestStatus == nil {
		return nil
	}
	return &RequestStatus{
		ID:        requestStatus.ID,
		Status:    requestStatus.Status.String(),
		CreatedAt: requestStatus.CreatedAt,
	}
}
