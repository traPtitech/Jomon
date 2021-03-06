// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent/predicate"
	"github.com/traPtitech/Jomon/ent/request"
	"github.com/traPtitech/Jomon/ent/requeststatus"
	"github.com/traPtitech/Jomon/ent/user"
)

// RequestStatusQuery is the builder for querying RequestStatus entities.
type RequestStatusQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.RequestStatus
	// eager-loading edges.
	withRequest *RequestQuery
	withUser    *UserQuery
	withFKs     bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the RequestStatusQuery builder.
func (rsq *RequestStatusQuery) Where(ps ...predicate.RequestStatus) *RequestStatusQuery {
	rsq.predicates = append(rsq.predicates, ps...)
	return rsq
}

// Limit adds a limit step to the query.
func (rsq *RequestStatusQuery) Limit(limit int) *RequestStatusQuery {
	rsq.limit = &limit
	return rsq
}

// Offset adds an offset step to the query.
func (rsq *RequestStatusQuery) Offset(offset int) *RequestStatusQuery {
	rsq.offset = &offset
	return rsq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (rsq *RequestStatusQuery) Unique(unique bool) *RequestStatusQuery {
	rsq.unique = &unique
	return rsq
}

// Order adds an order step to the query.
func (rsq *RequestStatusQuery) Order(o ...OrderFunc) *RequestStatusQuery {
	rsq.order = append(rsq.order, o...)
	return rsq
}

// QueryRequest chains the current query on the "request" edge.
func (rsq *RequestStatusQuery) QueryRequest() *RequestQuery {
	query := &RequestQuery{config: rsq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rsq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rsq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(requeststatus.Table, requeststatus.FieldID, selector),
			sqlgraph.To(request.Table, request.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, requeststatus.RequestTable, requeststatus.RequestColumn),
		)
		fromU = sqlgraph.SetNeighbors(rsq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryUser chains the current query on the "user" edge.
func (rsq *RequestStatusQuery) QueryUser() *UserQuery {
	query := &UserQuery{config: rsq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rsq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rsq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(requeststatus.Table, requeststatus.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, requeststatus.UserTable, requeststatus.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(rsq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first RequestStatus entity from the query.
// Returns a *NotFoundError when no RequestStatus was found.
func (rsq *RequestStatusQuery) First(ctx context.Context) (*RequestStatus, error) {
	nodes, err := rsq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{requeststatus.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (rsq *RequestStatusQuery) FirstX(ctx context.Context) *RequestStatus {
	node, err := rsq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first RequestStatus ID from the query.
// Returns a *NotFoundError when no RequestStatus ID was found.
func (rsq *RequestStatusQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = rsq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{requeststatus.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (rsq *RequestStatusQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := rsq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single RequestStatus entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when exactly one RequestStatus entity is not found.
// Returns a *NotFoundError when no RequestStatus entities are found.
func (rsq *RequestStatusQuery) Only(ctx context.Context) (*RequestStatus, error) {
	nodes, err := rsq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{requeststatus.Label}
	default:
		return nil, &NotSingularError{requeststatus.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (rsq *RequestStatusQuery) OnlyX(ctx context.Context) *RequestStatus {
	node, err := rsq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only RequestStatus ID in the query.
// Returns a *NotSingularError when exactly one RequestStatus ID is not found.
// Returns a *NotFoundError when no entities are found.
func (rsq *RequestStatusQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = rsq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{requeststatus.Label}
	default:
		err = &NotSingularError{requeststatus.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (rsq *RequestStatusQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := rsq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of RequestStatusSlice.
func (rsq *RequestStatusQuery) All(ctx context.Context) ([]*RequestStatus, error) {
	if err := rsq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return rsq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (rsq *RequestStatusQuery) AllX(ctx context.Context) []*RequestStatus {
	nodes, err := rsq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of RequestStatus IDs.
func (rsq *RequestStatusQuery) IDs(ctx context.Context) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	if err := rsq.Select(requeststatus.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (rsq *RequestStatusQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := rsq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (rsq *RequestStatusQuery) Count(ctx context.Context) (int, error) {
	if err := rsq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return rsq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (rsq *RequestStatusQuery) CountX(ctx context.Context) int {
	count, err := rsq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (rsq *RequestStatusQuery) Exist(ctx context.Context) (bool, error) {
	if err := rsq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return rsq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (rsq *RequestStatusQuery) ExistX(ctx context.Context) bool {
	exist, err := rsq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the RequestStatusQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (rsq *RequestStatusQuery) Clone() *RequestStatusQuery {
	if rsq == nil {
		return nil
	}
	return &RequestStatusQuery{
		config:      rsq.config,
		limit:       rsq.limit,
		offset:      rsq.offset,
		order:       append([]OrderFunc{}, rsq.order...),
		predicates:  append([]predicate.RequestStatus{}, rsq.predicates...),
		withRequest: rsq.withRequest.Clone(),
		withUser:    rsq.withUser.Clone(),
		// clone intermediate query.
		sql:  rsq.sql.Clone(),
		path: rsq.path,
	}
}

// WithRequest tells the query-builder to eager-load the nodes that are connected to
// the "request" edge. The optional arguments are used to configure the query builder of the edge.
func (rsq *RequestStatusQuery) WithRequest(opts ...func(*RequestQuery)) *RequestStatusQuery {
	query := &RequestQuery{config: rsq.config}
	for _, opt := range opts {
		opt(query)
	}
	rsq.withRequest = query
	return rsq
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (rsq *RequestStatusQuery) WithUser(opts ...func(*UserQuery)) *RequestStatusQuery {
	query := &UserQuery{config: rsq.config}
	for _, opt := range opts {
		opt(query)
	}
	rsq.withUser = query
	return rsq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Status requeststatus.Status `json:"status,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.RequestStatus.Query().
//		GroupBy(requeststatus.FieldStatus).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (rsq *RequestStatusQuery) GroupBy(field string, fields ...string) *RequestStatusGroupBy {
	group := &RequestStatusGroupBy{config: rsq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := rsq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return rsq.sqlQuery(ctx), nil
	}
	return group
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Status requeststatus.Status `json:"status,omitempty"`
//	}
//
//	client.RequestStatus.Query().
//		Select(requeststatus.FieldStatus).
//		Scan(ctx, &v)
//
func (rsq *RequestStatusQuery) Select(field string, fields ...string) *RequestStatusSelect {
	rsq.fields = append([]string{field}, fields...)
	return &RequestStatusSelect{RequestStatusQuery: rsq}
}

func (rsq *RequestStatusQuery) prepareQuery(ctx context.Context) error {
	for _, f := range rsq.fields {
		if !requeststatus.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if rsq.path != nil {
		prev, err := rsq.path(ctx)
		if err != nil {
			return err
		}
		rsq.sql = prev
	}
	return nil
}

func (rsq *RequestStatusQuery) sqlAll(ctx context.Context) ([]*RequestStatus, error) {
	var (
		nodes       = []*RequestStatus{}
		withFKs     = rsq.withFKs
		_spec       = rsq.querySpec()
		loadedTypes = [2]bool{
			rsq.withRequest != nil,
			rsq.withUser != nil,
		}
	)
	if rsq.withRequest != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, requeststatus.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &RequestStatus{config: rsq.config}
		nodes = append(nodes, node)
		return node.scanValues(columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if err := sqlgraph.QueryNodes(ctx, rsq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := rsq.withRequest; query != nil {
		ids := make([]uuid.UUID, 0, len(nodes))
		nodeids := make(map[uuid.UUID][]*RequestStatus)
		for i := range nodes {
			if nodes[i].request_status == nil {
				continue
			}
			fk := *nodes[i].request_status
			if _, ok := nodeids[fk]; !ok {
				ids = append(ids, fk)
			}
			nodeids[fk] = append(nodeids[fk], nodes[i])
		}
		query.Where(request.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "request_status" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Request = n
			}
		}
	}

	if query := rsq.withUser; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[uuid.UUID]*RequestStatus)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.User(func(s *sql.Selector) {
			s.Where(sql.InValues(requeststatus.UserColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.request_status_user
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "request_status_user" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "request_status_user" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.User = n
		}
	}

	return nodes, nil
}

func (rsq *RequestStatusQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := rsq.querySpec()
	return sqlgraph.CountNodes(ctx, rsq.driver, _spec)
}

func (rsq *RequestStatusQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := rsq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (rsq *RequestStatusQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   requeststatus.Table,
			Columns: requeststatus.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: requeststatus.FieldID,
			},
		},
		From:   rsq.sql,
		Unique: true,
	}
	if unique := rsq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := rsq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, requeststatus.FieldID)
		for i := range fields {
			if fields[i] != requeststatus.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := rsq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := rsq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := rsq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := rsq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (rsq *RequestStatusQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(rsq.driver.Dialect())
	t1 := builder.Table(requeststatus.Table)
	selector := builder.Select(t1.Columns(requeststatus.Columns...)...).From(t1)
	if rsq.sql != nil {
		selector = rsq.sql
		selector.Select(selector.Columns(requeststatus.Columns...)...)
	}
	for _, p := range rsq.predicates {
		p(selector)
	}
	for _, p := range rsq.order {
		p(selector)
	}
	if offset := rsq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := rsq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// RequestStatusGroupBy is the group-by builder for RequestStatus entities.
type RequestStatusGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (rsgb *RequestStatusGroupBy) Aggregate(fns ...AggregateFunc) *RequestStatusGroupBy {
	rsgb.fns = append(rsgb.fns, fns...)
	return rsgb
}

// Scan applies the group-by query and scans the result into the given value.
func (rsgb *RequestStatusGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := rsgb.path(ctx)
	if err != nil {
		return err
	}
	rsgb.sql = query
	return rsgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (rsgb *RequestStatusGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := rsgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (rsgb *RequestStatusGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(rsgb.fields) > 1 {
		return nil, errors.New("ent: RequestStatusGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := rsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (rsgb *RequestStatusGroupBy) StringsX(ctx context.Context) []string {
	v, err := rsgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (rsgb *RequestStatusGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = rsgb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{requeststatus.Label}
	default:
		err = fmt.Errorf("ent: RequestStatusGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (rsgb *RequestStatusGroupBy) StringX(ctx context.Context) string {
	v, err := rsgb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (rsgb *RequestStatusGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(rsgb.fields) > 1 {
		return nil, errors.New("ent: RequestStatusGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := rsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (rsgb *RequestStatusGroupBy) IntsX(ctx context.Context) []int {
	v, err := rsgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (rsgb *RequestStatusGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = rsgb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{requeststatus.Label}
	default:
		err = fmt.Errorf("ent: RequestStatusGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (rsgb *RequestStatusGroupBy) IntX(ctx context.Context) int {
	v, err := rsgb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (rsgb *RequestStatusGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(rsgb.fields) > 1 {
		return nil, errors.New("ent: RequestStatusGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := rsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (rsgb *RequestStatusGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := rsgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (rsgb *RequestStatusGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = rsgb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{requeststatus.Label}
	default:
		err = fmt.Errorf("ent: RequestStatusGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (rsgb *RequestStatusGroupBy) Float64X(ctx context.Context) float64 {
	v, err := rsgb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (rsgb *RequestStatusGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(rsgb.fields) > 1 {
		return nil, errors.New("ent: RequestStatusGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := rsgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (rsgb *RequestStatusGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := rsgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (rsgb *RequestStatusGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = rsgb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{requeststatus.Label}
	default:
		err = fmt.Errorf("ent: RequestStatusGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (rsgb *RequestStatusGroupBy) BoolX(ctx context.Context) bool {
	v, err := rsgb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (rsgb *RequestStatusGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range rsgb.fields {
		if !requeststatus.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := rsgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rsgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (rsgb *RequestStatusGroupBy) sqlQuery() *sql.Selector {
	selector := rsgb.sql
	columns := make([]string, 0, len(rsgb.fields)+len(rsgb.fns))
	columns = append(columns, rsgb.fields...)
	for _, fn := range rsgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(rsgb.fields...)
}

// RequestStatusSelect is the builder for selecting fields of RequestStatus entities.
type RequestStatusSelect struct {
	*RequestStatusQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (rss *RequestStatusSelect) Scan(ctx context.Context, v interface{}) error {
	if err := rss.prepareQuery(ctx); err != nil {
		return err
	}
	rss.sql = rss.RequestStatusQuery.sqlQuery(ctx)
	return rss.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (rss *RequestStatusSelect) ScanX(ctx context.Context, v interface{}) {
	if err := rss.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (rss *RequestStatusSelect) Strings(ctx context.Context) ([]string, error) {
	if len(rss.fields) > 1 {
		return nil, errors.New("ent: RequestStatusSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := rss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (rss *RequestStatusSelect) StringsX(ctx context.Context) []string {
	v, err := rss.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (rss *RequestStatusSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = rss.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{requeststatus.Label}
	default:
		err = fmt.Errorf("ent: RequestStatusSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (rss *RequestStatusSelect) StringX(ctx context.Context) string {
	v, err := rss.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (rss *RequestStatusSelect) Ints(ctx context.Context) ([]int, error) {
	if len(rss.fields) > 1 {
		return nil, errors.New("ent: RequestStatusSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := rss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (rss *RequestStatusSelect) IntsX(ctx context.Context) []int {
	v, err := rss.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (rss *RequestStatusSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = rss.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{requeststatus.Label}
	default:
		err = fmt.Errorf("ent: RequestStatusSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (rss *RequestStatusSelect) IntX(ctx context.Context) int {
	v, err := rss.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (rss *RequestStatusSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(rss.fields) > 1 {
		return nil, errors.New("ent: RequestStatusSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := rss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (rss *RequestStatusSelect) Float64sX(ctx context.Context) []float64 {
	v, err := rss.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (rss *RequestStatusSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = rss.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{requeststatus.Label}
	default:
		err = fmt.Errorf("ent: RequestStatusSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (rss *RequestStatusSelect) Float64X(ctx context.Context) float64 {
	v, err := rss.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (rss *RequestStatusSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(rss.fields) > 1 {
		return nil, errors.New("ent: RequestStatusSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := rss.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (rss *RequestStatusSelect) BoolsX(ctx context.Context) []bool {
	v, err := rss.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (rss *RequestStatusSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = rss.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{requeststatus.Label}
	default:
		err = fmt.Errorf("ent: RequestStatusSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (rss *RequestStatusSelect) BoolX(ctx context.Context) bool {
	v, err := rss.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (rss *RequestStatusSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := rss.sqlQuery().Query()
	if err := rss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (rss *RequestStatusSelect) sqlQuery() sql.Querier {
	selector := rss.sql
	selector.Select(selector.Columns(rss.fields...)...)
	return selector
}
