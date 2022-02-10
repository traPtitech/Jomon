// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent/predicate"
	"github.com/traPtitech/Jomon/ent/transaction"
	"github.com/traPtitech/Jomon/ent/transactiondetail"
)

// TransactionDetailQuery is the builder for querying TransactionDetail entities.
type TransactionDetailQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.TransactionDetail
	// eager-loading edges.
	withTransaction *TransactionQuery
	withFKs         bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the TransactionDetailQuery builder.
func (tdq *TransactionDetailQuery) Where(ps ...predicate.TransactionDetail) *TransactionDetailQuery {
	tdq.predicates = append(tdq.predicates, ps...)
	return tdq
}

// Limit adds a limit step to the query.
func (tdq *TransactionDetailQuery) Limit(limit int) *TransactionDetailQuery {
	tdq.limit = &limit
	return tdq
}

// Offset adds an offset step to the query.
func (tdq *TransactionDetailQuery) Offset(offset int) *TransactionDetailQuery {
	tdq.offset = &offset
	return tdq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (tdq *TransactionDetailQuery) Unique(unique bool) *TransactionDetailQuery {
	tdq.unique = &unique
	return tdq
}

// Order adds an order step to the query.
func (tdq *TransactionDetailQuery) Order(o ...OrderFunc) *TransactionDetailQuery {
	tdq.order = append(tdq.order, o...)
	return tdq
}

// QueryTransaction chains the current query on the "transaction" edge.
func (tdq *TransactionDetailQuery) QueryTransaction() *TransactionQuery {
	query := &TransactionQuery{config: tdq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := tdq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := tdq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(transactiondetail.Table, transactiondetail.FieldID, selector),
			sqlgraph.To(transaction.Table, transaction.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, transactiondetail.TransactionTable, transactiondetail.TransactionColumn),
		)
		fromU = sqlgraph.SetNeighbors(tdq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first TransactionDetail entity from the query.
// Returns a *NotFoundError when no TransactionDetail was found.
func (tdq *TransactionDetailQuery) First(ctx context.Context) (*TransactionDetail, error) {
	nodes, err := tdq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{transactiondetail.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (tdq *TransactionDetailQuery) FirstX(ctx context.Context) *TransactionDetail {
	node, err := tdq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first TransactionDetail ID from the query.
// Returns a *NotFoundError when no TransactionDetail ID was found.
func (tdq *TransactionDetailQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = tdq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{transactiondetail.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (tdq *TransactionDetailQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := tdq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single TransactionDetail entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when exactly one TransactionDetail entity is not found.
// Returns a *NotFoundError when no TransactionDetail entities are found.
func (tdq *TransactionDetailQuery) Only(ctx context.Context) (*TransactionDetail, error) {
	nodes, err := tdq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{transactiondetail.Label}
	default:
		return nil, &NotSingularError{transactiondetail.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (tdq *TransactionDetailQuery) OnlyX(ctx context.Context) *TransactionDetail {
	node, err := tdq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only TransactionDetail ID in the query.
// Returns a *NotSingularError when exactly one TransactionDetail ID is not found.
// Returns a *NotFoundError when no entities are found.
func (tdq *TransactionDetailQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = tdq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{transactiondetail.Label}
	default:
		err = &NotSingularError{transactiondetail.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (tdq *TransactionDetailQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := tdq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of TransactionDetails.
func (tdq *TransactionDetailQuery) All(ctx context.Context) ([]*TransactionDetail, error) {
	if err := tdq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return tdq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (tdq *TransactionDetailQuery) AllX(ctx context.Context) []*TransactionDetail {
	nodes, err := tdq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of TransactionDetail IDs.
func (tdq *TransactionDetailQuery) IDs(ctx context.Context) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	if err := tdq.Select(transactiondetail.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (tdq *TransactionDetailQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := tdq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (tdq *TransactionDetailQuery) Count(ctx context.Context) (int, error) {
	if err := tdq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return tdq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (tdq *TransactionDetailQuery) CountX(ctx context.Context) int {
	count, err := tdq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (tdq *TransactionDetailQuery) Exist(ctx context.Context) (bool, error) {
	if err := tdq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return tdq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (tdq *TransactionDetailQuery) ExistX(ctx context.Context) bool {
	exist, err := tdq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the TransactionDetailQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (tdq *TransactionDetailQuery) Clone() *TransactionDetailQuery {
	if tdq == nil {
		return nil
	}
	return &TransactionDetailQuery{
		config:          tdq.config,
		limit:           tdq.limit,
		offset:          tdq.offset,
		order:           append([]OrderFunc{}, tdq.order...),
		predicates:      append([]predicate.TransactionDetail{}, tdq.predicates...),
		withTransaction: tdq.withTransaction.Clone(),
		// clone intermediate query.
		sql:  tdq.sql.Clone(),
		path: tdq.path,
	}
}

// WithTransaction tells the query-builder to eager-load the nodes that are connected to
// the "transaction" edge. The optional arguments are used to configure the query builder of the edge.
func (tdq *TransactionDetailQuery) WithTransaction(opts ...func(*TransactionQuery)) *TransactionDetailQuery {
	query := &TransactionQuery{config: tdq.config}
	for _, opt := range opts {
		opt(query)
	}
	tdq.withTransaction = query
	return tdq
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
//	client.TransactionDetail.Query().
//		GroupBy(transactiondetail.FieldAmount).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (tdq *TransactionDetailQuery) GroupBy(field string, fields ...string) *TransactionDetailGroupBy {
	group := &TransactionDetailGroupBy{config: tdq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := tdq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return tdq.sqlQuery(ctx), nil
	}
	return group
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
//	client.TransactionDetail.Query().
//		Select(transactiondetail.FieldAmount).
//		Scan(ctx, &v)
//
func (tdq *TransactionDetailQuery) Select(fields ...string) *TransactionDetailSelect {
	tdq.fields = append(tdq.fields, fields...)
	return &TransactionDetailSelect{TransactionDetailQuery: tdq}
}

func (tdq *TransactionDetailQuery) prepareQuery(ctx context.Context) error {
	for _, f := range tdq.fields {
		if !transactiondetail.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if tdq.path != nil {
		prev, err := tdq.path(ctx)
		if err != nil {
			return err
		}
		tdq.sql = prev
	}
	return nil
}

func (tdq *TransactionDetailQuery) sqlAll(ctx context.Context) ([]*TransactionDetail, error) {
	var (
		nodes       = []*TransactionDetail{}
		withFKs     = tdq.withFKs
		_spec       = tdq.querySpec()
		loadedTypes = [1]bool{
			tdq.withTransaction != nil,
		}
	)
	if tdq.withTransaction != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, transactiondetail.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &TransactionDetail{config: tdq.config}
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
	if err := sqlgraph.QueryNodes(ctx, tdq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := tdq.withTransaction; query != nil {
		ids := make([]uuid.UUID, 0, len(nodes))
		nodeids := make(map[uuid.UUID][]*TransactionDetail)
		for i := range nodes {
			if nodes[i].transaction_detail == nil {
				continue
			}
			fk := *nodes[i].transaction_detail
			if _, ok := nodeids[fk]; !ok {
				ids = append(ids, fk)
			}
			nodeids[fk] = append(nodeids[fk], nodes[i])
		}
		query.Where(transaction.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "transaction_detail" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Transaction = n
			}
		}
	}

	return nodes, nil
}

func (tdq *TransactionDetailQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := tdq.querySpec()
	_spec.Node.Columns = tdq.fields
	if len(tdq.fields) > 0 {
		_spec.Unique = tdq.unique != nil && *tdq.unique
	}
	return sqlgraph.CountNodes(ctx, tdq.driver, _spec)
}

func (tdq *TransactionDetailQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := tdq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (tdq *TransactionDetailQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   transactiondetail.Table,
			Columns: transactiondetail.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: transactiondetail.FieldID,
			},
		},
		From:   tdq.sql,
		Unique: true,
	}
	if unique := tdq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := tdq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, transactiondetail.FieldID)
		for i := range fields {
			if fields[i] != transactiondetail.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := tdq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := tdq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := tdq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := tdq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (tdq *TransactionDetailQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(tdq.driver.Dialect())
	t1 := builder.Table(transactiondetail.Table)
	columns := tdq.fields
	if len(columns) == 0 {
		columns = transactiondetail.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if tdq.sql != nil {
		selector = tdq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if tdq.unique != nil && *tdq.unique {
		selector.Distinct()
	}
	for _, p := range tdq.predicates {
		p(selector)
	}
	for _, p := range tdq.order {
		p(selector)
	}
	if offset := tdq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := tdq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// TransactionDetailGroupBy is the group-by builder for TransactionDetail entities.
type TransactionDetailGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (tdgb *TransactionDetailGroupBy) Aggregate(fns ...AggregateFunc) *TransactionDetailGroupBy {
	tdgb.fns = append(tdgb.fns, fns...)
	return tdgb
}

// Scan applies the group-by query and scans the result into the given value.
func (tdgb *TransactionDetailGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := tdgb.path(ctx)
	if err != nil {
		return err
	}
	tdgb.sql = query
	return tdgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (tdgb *TransactionDetailGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := tdgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (tdgb *TransactionDetailGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(tdgb.fields) > 1 {
		return nil, errors.New("ent: TransactionDetailGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := tdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (tdgb *TransactionDetailGroupBy) StringsX(ctx context.Context) []string {
	v, err := tdgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (tdgb *TransactionDetailGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = tdgb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{transactiondetail.Label}
	default:
		err = fmt.Errorf("ent: TransactionDetailGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (tdgb *TransactionDetailGroupBy) StringX(ctx context.Context) string {
	v, err := tdgb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (tdgb *TransactionDetailGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(tdgb.fields) > 1 {
		return nil, errors.New("ent: TransactionDetailGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := tdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (tdgb *TransactionDetailGroupBy) IntsX(ctx context.Context) []int {
	v, err := tdgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (tdgb *TransactionDetailGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = tdgb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{transactiondetail.Label}
	default:
		err = fmt.Errorf("ent: TransactionDetailGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (tdgb *TransactionDetailGroupBy) IntX(ctx context.Context) int {
	v, err := tdgb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (tdgb *TransactionDetailGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(tdgb.fields) > 1 {
		return nil, errors.New("ent: TransactionDetailGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := tdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (tdgb *TransactionDetailGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := tdgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (tdgb *TransactionDetailGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = tdgb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{transactiondetail.Label}
	default:
		err = fmt.Errorf("ent: TransactionDetailGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (tdgb *TransactionDetailGroupBy) Float64X(ctx context.Context) float64 {
	v, err := tdgb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (tdgb *TransactionDetailGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(tdgb.fields) > 1 {
		return nil, errors.New("ent: TransactionDetailGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := tdgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (tdgb *TransactionDetailGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := tdgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (tdgb *TransactionDetailGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = tdgb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{transactiondetail.Label}
	default:
		err = fmt.Errorf("ent: TransactionDetailGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (tdgb *TransactionDetailGroupBy) BoolX(ctx context.Context) bool {
	v, err := tdgb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (tdgb *TransactionDetailGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range tdgb.fields {
		if !transactiondetail.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := tdgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := tdgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (tdgb *TransactionDetailGroupBy) sqlQuery() *sql.Selector {
	selector := tdgb.sql.Select()
	aggregation := make([]string, 0, len(tdgb.fns))
	for _, fn := range tdgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(tdgb.fields)+len(tdgb.fns))
		for _, f := range tdgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(tdgb.fields...)...)
}

// TransactionDetailSelect is the builder for selecting fields of TransactionDetail entities.
type TransactionDetailSelect struct {
	*TransactionDetailQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (tds *TransactionDetailSelect) Scan(ctx context.Context, v interface{}) error {
	if err := tds.prepareQuery(ctx); err != nil {
		return err
	}
	tds.sql = tds.TransactionDetailQuery.sqlQuery(ctx)
	return tds.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (tds *TransactionDetailSelect) ScanX(ctx context.Context, v interface{}) {
	if err := tds.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (tds *TransactionDetailSelect) Strings(ctx context.Context) ([]string, error) {
	if len(tds.fields) > 1 {
		return nil, errors.New("ent: TransactionDetailSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := tds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (tds *TransactionDetailSelect) StringsX(ctx context.Context) []string {
	v, err := tds.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (tds *TransactionDetailSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = tds.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{transactiondetail.Label}
	default:
		err = fmt.Errorf("ent: TransactionDetailSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (tds *TransactionDetailSelect) StringX(ctx context.Context) string {
	v, err := tds.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (tds *TransactionDetailSelect) Ints(ctx context.Context) ([]int, error) {
	if len(tds.fields) > 1 {
		return nil, errors.New("ent: TransactionDetailSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := tds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (tds *TransactionDetailSelect) IntsX(ctx context.Context) []int {
	v, err := tds.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (tds *TransactionDetailSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = tds.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{transactiondetail.Label}
	default:
		err = fmt.Errorf("ent: TransactionDetailSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (tds *TransactionDetailSelect) IntX(ctx context.Context) int {
	v, err := tds.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (tds *TransactionDetailSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(tds.fields) > 1 {
		return nil, errors.New("ent: TransactionDetailSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := tds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (tds *TransactionDetailSelect) Float64sX(ctx context.Context) []float64 {
	v, err := tds.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (tds *TransactionDetailSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = tds.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{transactiondetail.Label}
	default:
		err = fmt.Errorf("ent: TransactionDetailSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (tds *TransactionDetailSelect) Float64X(ctx context.Context) float64 {
	v, err := tds.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (tds *TransactionDetailSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(tds.fields) > 1 {
		return nil, errors.New("ent: TransactionDetailSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := tds.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (tds *TransactionDetailSelect) BoolsX(ctx context.Context) []bool {
	v, err := tds.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (tds *TransactionDetailSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = tds.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{transactiondetail.Label}
	default:
		err = fmt.Errorf("ent: TransactionDetailSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (tds *TransactionDetailSelect) BoolX(ctx context.Context) bool {
	v, err := tds.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (tds *TransactionDetailSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := tds.sql.Query()
	if err := tds.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
