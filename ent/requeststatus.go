// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent/request"
	"github.com/traPtitech/Jomon/ent/requeststatus"
	"github.com/traPtitech/Jomon/ent/user"
)

// RequestStatus is the model entity for the RequestStatus schema.
type RequestStatus struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Status holds the value of the "status" field.
	Status requeststatus.Status `json:"status,omitempty"`
	// Reason holds the value of the "reason" field.
	Reason string `json:"reason,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the RequestStatusQuery when eager-loading is set.
	Edges               RequestStatusEdges `json:"edges"`
	request_status      *uuid.UUID
	request_status_user *uuid.UUID
}

// RequestStatusEdges holds the relations/edges for other nodes in the graph.
type RequestStatusEdges struct {
	// Request holds the value of the request edge.
	Request *Request `json:"request,omitempty"`
	// User holds the value of the user edge.
	User *User `json:"user,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// RequestOrErr returns the Request value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e RequestStatusEdges) RequestOrErr() (*Request, error) {
	if e.loadedTypes[0] {
		if e.Request == nil {
			// The edge request was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: request.Label}
		}
		return e.Request, nil
	}
	return nil, &NotLoadedError{edge: "request"}
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e RequestStatusEdges) UserOrErr() (*User, error) {
	if e.loadedTypes[1] {
		if e.User == nil {
			// The edge user was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.User, nil
	}
	return nil, &NotLoadedError{edge: "user"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*RequestStatus) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case requeststatus.FieldStatus, requeststatus.FieldReason:
			values[i] = new(sql.NullString)
		case requeststatus.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		case requeststatus.FieldID:
			values[i] = new(uuid.UUID)
		case requeststatus.ForeignKeys[0]: // request_status
			values[i] = new(uuid.UUID)
		case requeststatus.ForeignKeys[1]: // request_status_user
			values[i] = new(uuid.UUID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type RequestStatus", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the RequestStatus fields.
func (rs *RequestStatus) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case requeststatus.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				rs.ID = *value
			}
		case requeststatus.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				rs.Status = requeststatus.Status(value.String)
			}
		case requeststatus.FieldReason:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field reason", values[i])
			} else if value.Valid {
				rs.Reason = value.String
			}
		case requeststatus.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				rs.CreatedAt = value.Time
			}
		case requeststatus.ForeignKeys[0]:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field request_status", values[i])
			} else if value != nil {
				rs.request_status = value
			}
		case requeststatus.ForeignKeys[1]:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field request_status_user", values[i])
			} else if value != nil {
				rs.request_status_user = value
			}
		}
	}
	return nil
}

// QueryRequest queries the "request" edge of the RequestStatus entity.
func (rs *RequestStatus) QueryRequest() *RequestQuery {
	return (&RequestStatusClient{config: rs.config}).QueryRequest(rs)
}

// QueryUser queries the "user" edge of the RequestStatus entity.
func (rs *RequestStatus) QueryUser() *UserQuery {
	return (&RequestStatusClient{config: rs.config}).QueryUser(rs)
}

// Update returns a builder for updating this RequestStatus.
// Note that you need to call RequestStatus.Unwrap() before calling this method if this RequestStatus
// was returned from a transaction, and the transaction was committed or rolled back.
func (rs *RequestStatus) Update() *RequestStatusUpdateOne {
	return (&RequestStatusClient{config: rs.config}).UpdateOne(rs)
}

// Unwrap unwraps the RequestStatus entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (rs *RequestStatus) Unwrap() *RequestStatus {
	tx, ok := rs.config.driver.(*txDriver)
	if !ok {
		panic("ent: RequestStatus is not a transactional entity")
	}
	rs.config.driver = tx.drv
	return rs
}

// String implements the fmt.Stringer.
func (rs *RequestStatus) String() string {
	var builder strings.Builder
	builder.WriteString("RequestStatus(")
	builder.WriteString(fmt.Sprintf("id=%v", rs.ID))
	builder.WriteString(", status=")
	builder.WriteString(fmt.Sprintf("%v", rs.Status))
	builder.WriteString(", reason=")
	builder.WriteString(rs.Reason)
	builder.WriteString(", created_at=")
	builder.WriteString(rs.CreatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// RequestStatusSlice is a parsable slice of RequestStatus.
type RequestStatusSlice []*RequestStatus

func (rs RequestStatusSlice) config(cfg config) {
	for _i := range rs {
		rs[_i].config = cfg
	}
}
