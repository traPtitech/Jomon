// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent/predicate"
	"github.com/traPtitech/Jomon/ent/request"
	"github.com/traPtitech/Jomon/ent/requesttarget"
	"github.com/traPtitech/Jomon/ent/user"
)

// RequestTargetQuery is the builder for querying RequestTarget entities.
type RequestTargetQuery struct {
	config
	ctx         *QueryContext
	order       []OrderFunc
	inters      []Interceptor
	predicates  []predicate.RequestTarget
	withRequest *RequestQuery
	withUser    *UserQuery
	withFKs     bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the RequestTargetQuery builder.
func (rtq *RequestTargetQuery) Where(ps ...predicate.RequestTarget) *RequestTargetQuery {
	rtq.predicates = append(rtq.predicates, ps...)
	return rtq
}

// Limit the number of records to be returned by this query.
func (rtq *RequestTargetQuery) Limit(limit int) *RequestTargetQuery {
	rtq.ctx.Limit = &limit
	return rtq
}

// Offset to start from.
func (rtq *RequestTargetQuery) Offset(offset int) *RequestTargetQuery {
	rtq.ctx.Offset = &offset
	return rtq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (rtq *RequestTargetQuery) Unique(unique bool) *RequestTargetQuery {
	rtq.ctx.Unique = &unique
	return rtq
}

// Order specifies how the records should be ordered.
func (rtq *RequestTargetQuery) Order(o ...OrderFunc) *RequestTargetQuery {
	rtq.order = append(rtq.order, o...)
	return rtq
}

// QueryRequest chains the current query on the "request" edge.
func (rtq *RequestTargetQuery) QueryRequest() *RequestQuery {
	query := (&RequestClient{config: rtq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rtq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rtq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(requesttarget.Table, requesttarget.FieldID, selector),
			sqlgraph.To(request.Table, request.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, requesttarget.RequestTable, requesttarget.RequestColumn),
		)
		fromU = sqlgraph.SetNeighbors(rtq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryUser chains the current query on the "user" edge.
func (rtq *RequestTargetQuery) QueryUser() *UserQuery {
	query := (&UserClient{config: rtq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rtq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rtq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(requesttarget.Table, requesttarget.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, requesttarget.UserTable, requesttarget.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(rtq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first RequestTarget entity from the query.
// Returns a *NotFoundError when no RequestTarget was found.
func (rtq *RequestTargetQuery) First(ctx context.Context) (*RequestTarget, error) {
	nodes, err := rtq.Limit(1).All(setContextOp(ctx, rtq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{requesttarget.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (rtq *RequestTargetQuery) FirstX(ctx context.Context) *RequestTarget {
	node, err := rtq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first RequestTarget ID from the query.
// Returns a *NotFoundError when no RequestTarget ID was found.
func (rtq *RequestTargetQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = rtq.Limit(1).IDs(setContextOp(ctx, rtq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{requesttarget.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (rtq *RequestTargetQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := rtq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single RequestTarget entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one RequestTarget entity is found.
// Returns a *NotFoundError when no RequestTarget entities are found.
func (rtq *RequestTargetQuery) Only(ctx context.Context) (*RequestTarget, error) {
	nodes, err := rtq.Limit(2).All(setContextOp(ctx, rtq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{requesttarget.Label}
	default:
		return nil, &NotSingularError{requesttarget.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (rtq *RequestTargetQuery) OnlyX(ctx context.Context) *RequestTarget {
	node, err := rtq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only RequestTarget ID in the query.
// Returns a *NotSingularError when more than one RequestTarget ID is found.
// Returns a *NotFoundError when no entities are found.
func (rtq *RequestTargetQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = rtq.Limit(2).IDs(setContextOp(ctx, rtq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{requesttarget.Label}
	default:
		err = &NotSingularError{requesttarget.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (rtq *RequestTargetQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := rtq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of RequestTargets.
func (rtq *RequestTargetQuery) All(ctx context.Context) ([]*RequestTarget, error) {
	ctx = setContextOp(ctx, rtq.ctx, "All")
	if err := rtq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*RequestTarget, *RequestTargetQuery]()
	return withInterceptors[[]*RequestTarget](ctx, rtq, qr, rtq.inters)
}

// AllX is like All, but panics if an error occurs.
func (rtq *RequestTargetQuery) AllX(ctx context.Context) []*RequestTarget {
	nodes, err := rtq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of RequestTarget IDs.
func (rtq *RequestTargetQuery) IDs(ctx context.Context) (ids []uuid.UUID, err error) {
	if rtq.ctx.Unique == nil && rtq.path != nil {
		rtq.Unique(true)
	}
	ctx = setContextOp(ctx, rtq.ctx, "IDs")
	if err = rtq.Select(requesttarget.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (rtq *RequestTargetQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := rtq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (rtq *RequestTargetQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, rtq.ctx, "Count")
	if err := rtq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, rtq, querierCount[*RequestTargetQuery](), rtq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (rtq *RequestTargetQuery) CountX(ctx context.Context) int {
	count, err := rtq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (rtq *RequestTargetQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, rtq.ctx, "Exist")
	switch _, err := rtq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (rtq *RequestTargetQuery) ExistX(ctx context.Context) bool {
	exist, err := rtq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the RequestTargetQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (rtq *RequestTargetQuery) Clone() *RequestTargetQuery {
	if rtq == nil {
		return nil
	}
	return &RequestTargetQuery{
		config:      rtq.config,
		ctx:         rtq.ctx.Clone(),
		order:       append([]OrderFunc{}, rtq.order...),
		inters:      append([]Interceptor{}, rtq.inters...),
		predicates:  append([]predicate.RequestTarget{}, rtq.predicates...),
		withRequest: rtq.withRequest.Clone(),
		withUser:    rtq.withUser.Clone(),
		// clone intermediate query.
		sql:  rtq.sql.Clone(),
		path: rtq.path,
	}
}

// WithRequest tells the query-builder to eager-load the nodes that are connected to
// the "request" edge. The optional arguments are used to configure the query builder of the edge.
func (rtq *RequestTargetQuery) WithRequest(opts ...func(*RequestQuery)) *RequestTargetQuery {
	query := (&RequestClient{config: rtq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	rtq.withRequest = query
	return rtq
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (rtq *RequestTargetQuery) WithUser(opts ...func(*UserQuery)) *RequestTargetQuery {
	query := (&UserClient{config: rtq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	rtq.withUser = query
	return rtq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Amount int `json:"amount,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.RequestTarget.Query().
//		GroupBy(requesttarget.FieldAmount).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (rtq *RequestTargetQuery) GroupBy(field string, fields ...string) *RequestTargetGroupBy {
	rtq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &RequestTargetGroupBy{build: rtq}
	grbuild.flds = &rtq.ctx.Fields
	grbuild.label = requesttarget.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Amount int `json:"amount,omitempty"`
//	}
//
//	client.RequestTarget.Query().
//		Select(requesttarget.FieldAmount).
//		Scan(ctx, &v)
func (rtq *RequestTargetQuery) Select(fields ...string) *RequestTargetSelect {
	rtq.ctx.Fields = append(rtq.ctx.Fields, fields...)
	sbuild := &RequestTargetSelect{RequestTargetQuery: rtq}
	sbuild.label = requesttarget.Label
	sbuild.flds, sbuild.scan = &rtq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a RequestTargetSelect configured with the given aggregations.
func (rtq *RequestTargetQuery) Aggregate(fns ...AggregateFunc) *RequestTargetSelect {
	return rtq.Select().Aggregate(fns...)
}

func (rtq *RequestTargetQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range rtq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, rtq); err != nil {
				return err
			}
		}
	}
	for _, f := range rtq.ctx.Fields {
		if !requesttarget.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if rtq.path != nil {
		prev, err := rtq.path(ctx)
		if err != nil {
			return err
		}
		rtq.sql = prev
	}
	return nil
}

func (rtq *RequestTargetQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*RequestTarget, error) {
	var (
		nodes       = []*RequestTarget{}
		withFKs     = rtq.withFKs
		_spec       = rtq.querySpec()
		loadedTypes = [2]bool{
			rtq.withRequest != nil,
			rtq.withUser != nil,
		}
	)
	if rtq.withRequest != nil || rtq.withUser != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, requesttarget.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*RequestTarget).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &RequestTarget{config: rtq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, rtq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := rtq.withRequest; query != nil {
		if err := rtq.loadRequest(ctx, query, nodes, nil,
			func(n *RequestTarget, e *Request) { n.Edges.Request = e }); err != nil {
			return nil, err
		}
	}
	if query := rtq.withUser; query != nil {
		if err := rtq.loadUser(ctx, query, nodes, nil,
			func(n *RequestTarget, e *User) { n.Edges.User = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (rtq *RequestTargetQuery) loadRequest(ctx context.Context, query *RequestQuery, nodes []*RequestTarget, init func(*RequestTarget), assign func(*RequestTarget, *Request)) error {
	ids := make([]uuid.UUID, 0, len(nodes))
	nodeids := make(map[uuid.UUID][]*RequestTarget)
	for i := range nodes {
		if nodes[i].request_target == nil {
			continue
		}
		fk := *nodes[i].request_target
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(request.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "request_target" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (rtq *RequestTargetQuery) loadUser(ctx context.Context, query *UserQuery, nodes []*RequestTarget, init func(*RequestTarget), assign func(*RequestTarget, *User)) error {
	ids := make([]uuid.UUID, 0, len(nodes))
	nodeids := make(map[uuid.UUID][]*RequestTarget)
	for i := range nodes {
		if nodes[i].request_target_user == nil {
			continue
		}
		fk := *nodes[i].request_target_user
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(user.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "request_target_user" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (rtq *RequestTargetQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := rtq.querySpec()
	_spec.Node.Columns = rtq.ctx.Fields
	if len(rtq.ctx.Fields) > 0 {
		_spec.Unique = rtq.ctx.Unique != nil && *rtq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, rtq.driver, _spec)
}

func (rtq *RequestTargetQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(requesttarget.Table, requesttarget.Columns, sqlgraph.NewFieldSpec(requesttarget.FieldID, field.TypeUUID))
	_spec.From = rtq.sql
	if unique := rtq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if rtq.path != nil {
		_spec.Unique = true
	}
	if fields := rtq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, requesttarget.FieldID)
		for i := range fields {
			if fields[i] != requesttarget.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := rtq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := rtq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := rtq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := rtq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (rtq *RequestTargetQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(rtq.driver.Dialect())
	t1 := builder.Table(requesttarget.Table)
	columns := rtq.ctx.Fields
	if len(columns) == 0 {
		columns = requesttarget.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if rtq.sql != nil {
		selector = rtq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if rtq.ctx.Unique != nil && *rtq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range rtq.predicates {
		p(selector)
	}
	for _, p := range rtq.order {
		p(selector)
	}
	if offset := rtq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := rtq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// RequestTargetGroupBy is the group-by builder for RequestTarget entities.
type RequestTargetGroupBy struct {
	selector
	build *RequestTargetQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (rtgb *RequestTargetGroupBy) Aggregate(fns ...AggregateFunc) *RequestTargetGroupBy {
	rtgb.fns = append(rtgb.fns, fns...)
	return rtgb
}

// Scan applies the selector query and scans the result into the given value.
func (rtgb *RequestTargetGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, rtgb.build.ctx, "GroupBy")
	if err := rtgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*RequestTargetQuery, *RequestTargetGroupBy](ctx, rtgb.build, rtgb, rtgb.build.inters, v)
}

func (rtgb *RequestTargetGroupBy) sqlScan(ctx context.Context, root *RequestTargetQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(rtgb.fns))
	for _, fn := range rtgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*rtgb.flds)+len(rtgb.fns))
		for _, f := range *rtgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*rtgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rtgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// RequestTargetSelect is the builder for selecting fields of RequestTarget entities.
type RequestTargetSelect struct {
	*RequestTargetQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (rts *RequestTargetSelect) Aggregate(fns ...AggregateFunc) *RequestTargetSelect {
	rts.fns = append(rts.fns, fns...)
	return rts
}

// Scan applies the selector query and scans the result into the given value.
func (rts *RequestTargetSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, rts.ctx, "Select")
	if err := rts.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*RequestTargetQuery, *RequestTargetSelect](ctx, rts.RequestTargetQuery, rts, rts.inters, v)
}

func (rts *RequestTargetSelect) sqlScan(ctx context.Context, root *RequestTargetQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(rts.fns))
	for _, fn := range rts.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*rts.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rts.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
