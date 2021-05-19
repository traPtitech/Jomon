package model

import (
	"encoding/json"
	"fmt"
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

type RequestStatusRepository interface {
}
