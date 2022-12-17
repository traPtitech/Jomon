// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent/comment"
	"github.com/traPtitech/Jomon/ent/file"
	"github.com/traPtitech/Jomon/ent/group"
	"github.com/traPtitech/Jomon/ent/predicate"
	"github.com/traPtitech/Jomon/ent/request"
	"github.com/traPtitech/Jomon/ent/requeststatus"
	"github.com/traPtitech/Jomon/ent/requesttarget"
	"github.com/traPtitech/Jomon/ent/tag"
	"github.com/traPtitech/Jomon/ent/transaction"
	"github.com/traPtitech/Jomon/ent/user"
)

// RequestQuery is the builder for querying Request entities.
type RequestQuery struct {
	config
	limit           *int
	offset          *int
	unique          *bool
	order           []OrderFunc
	fields          []string
	predicates      []predicate.Request
	withStatus      *RequestStatusQuery
	withTarget      *RequestTargetQuery
	withFile        *FileQuery
	withTag         *TagQuery
	withTransaction *TransactionQuery
	withComment     *CommentQuery
	withUser        *UserQuery
	withGroup       *GroupQuery
	withFKs         bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the RequestQuery builder.
func (rq *RequestQuery) Where(ps ...predicate.Request) *RequestQuery {
	rq.predicates = append(rq.predicates, ps...)
	return rq
}

// Limit adds a limit step to the query.
func (rq *RequestQuery) Limit(limit int) *RequestQuery {
	rq.limit = &limit
	return rq
}

// Offset adds an offset step to the query.
func (rq *RequestQuery) Offset(offset int) *RequestQuery {
	rq.offset = &offset
	return rq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (rq *RequestQuery) Unique(unique bool) *RequestQuery {
	rq.unique = &unique
	return rq
}

// Order adds an order step to the query.
func (rq *RequestQuery) Order(o ...OrderFunc) *RequestQuery {
	rq.order = append(rq.order, o...)
	return rq
}

// QueryStatus chains the current query on the "status" edge.
func (rq *RequestQuery) QueryStatus() *RequestStatusQuery {
	query := &RequestStatusQuery{config: rq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(request.Table, request.FieldID, selector),
			sqlgraph.To(requeststatus.Table, requeststatus.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, request.StatusTable, request.StatusColumn),
		)
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryTarget chains the current query on the "target" edge.
func (rq *RequestQuery) QueryTarget() *RequestTargetQuery {
	query := &RequestTargetQuery{config: rq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(request.Table, request.FieldID, selector),
			sqlgraph.To(requesttarget.Table, requesttarget.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, request.TargetTable, request.TargetColumn),
		)
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryFile chains the current query on the "file" edge.
func (rq *RequestQuery) QueryFile() *FileQuery {
	query := &FileQuery{config: rq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(request.Table, request.FieldID, selector),
			sqlgraph.To(file.Table, file.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, request.FileTable, request.FileColumn),
		)
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryTag chains the current query on the "tag" edge.
func (rq *RequestQuery) QueryTag() *TagQuery {
	query := &TagQuery{config: rq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(request.Table, request.FieldID, selector),
			sqlgraph.To(tag.Table, tag.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, request.TagTable, request.TagPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryTransaction chains the current query on the "transaction" edge.
func (rq *RequestQuery) QueryTransaction() *TransactionQuery {
	query := &TransactionQuery{config: rq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(request.Table, request.FieldID, selector),
			sqlgraph.To(transaction.Table, transaction.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, request.TransactionTable, request.TransactionColumn),
		)
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryComment chains the current query on the "comment" edge.
func (rq *RequestQuery) QueryComment() *CommentQuery {
	query := &CommentQuery{config: rq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(request.Table, request.FieldID, selector),
			sqlgraph.To(comment.Table, comment.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, request.CommentTable, request.CommentColumn),
		)
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryUser chains the current query on the "user" edge.
func (rq *RequestQuery) QueryUser() *UserQuery {
	query := &UserQuery{config: rq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(request.Table, request.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, request.UserTable, request.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryGroup chains the current query on the "group" edge.
func (rq *RequestQuery) QueryGroup() *GroupQuery {
	query := &GroupQuery{config: rq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(request.Table, request.FieldID, selector),
			sqlgraph.To(group.Table, group.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, request.GroupTable, request.GroupColumn),
		)
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Request entity from the query.
// Returns a *NotFoundError when no Request was found.
func (rq *RequestQuery) First(ctx context.Context) (*Request, error) {
	nodes, err := rq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{request.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (rq *RequestQuery) FirstX(ctx context.Context) *Request {
	node, err := rq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Request ID from the query.
// Returns a *NotFoundError when no Request ID was found.
func (rq *RequestQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = rq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{request.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (rq *RequestQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := rq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Request entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Request entity is found.
// Returns a *NotFoundError when no Request entities are found.
func (rq *RequestQuery) Only(ctx context.Context) (*Request, error) {
	nodes, err := rq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{request.Label}
	default:
		return nil, &NotSingularError{request.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (rq *RequestQuery) OnlyX(ctx context.Context) *Request {
	node, err := rq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Request ID in the query.
// Returns a *NotSingularError when more than one Request ID is found.
// Returns a *NotFoundError when no entities are found.
func (rq *RequestQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = rq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{request.Label}
	default:
		err = &NotSingularError{request.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (rq *RequestQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := rq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Requests.
func (rq *RequestQuery) All(ctx context.Context) ([]*Request, error) {
	if err := rq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return rq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (rq *RequestQuery) AllX(ctx context.Context) []*Request {
	nodes, err := rq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Request IDs.
func (rq *RequestQuery) IDs(ctx context.Context) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	if err := rq.Select(request.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (rq *RequestQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := rq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (rq *RequestQuery) Count(ctx context.Context) (int, error) {
	if err := rq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return rq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (rq *RequestQuery) CountX(ctx context.Context) int {
	count, err := rq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (rq *RequestQuery) Exist(ctx context.Context) (bool, error) {
	if err := rq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return rq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (rq *RequestQuery) ExistX(ctx context.Context) bool {
	exist, err := rq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the RequestQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (rq *RequestQuery) Clone() *RequestQuery {
	if rq == nil {
		return nil
	}
	return &RequestQuery{
		config:          rq.config,
		limit:           rq.limit,
		offset:          rq.offset,
		order:           append([]OrderFunc{}, rq.order...),
		predicates:      append([]predicate.Request{}, rq.predicates...),
		withStatus:      rq.withStatus.Clone(),
		withTarget:      rq.withTarget.Clone(),
		withFile:        rq.withFile.Clone(),
		withTag:         rq.withTag.Clone(),
		withTransaction: rq.withTransaction.Clone(),
		withComment:     rq.withComment.Clone(),
		withUser:        rq.withUser.Clone(),
		withGroup:       rq.withGroup.Clone(),
		// clone intermediate query.
		sql:    rq.sql.Clone(),
		path:   rq.path,
		unique: rq.unique,
	}
}

// WithStatus tells the query-builder to eager-load the nodes that are connected to
// the "status" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *RequestQuery) WithStatus(opts ...func(*RequestStatusQuery)) *RequestQuery {
	query := &RequestStatusQuery{config: rq.config}
	for _, opt := range opts {
		opt(query)
	}
	rq.withStatus = query
	return rq
}

// WithTarget tells the query-builder to eager-load the nodes that are connected to
// the "target" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *RequestQuery) WithTarget(opts ...func(*RequestTargetQuery)) *RequestQuery {
	query := &RequestTargetQuery{config: rq.config}
	for _, opt := range opts {
		opt(query)
	}
	rq.withTarget = query
	return rq
}

// WithFile tells the query-builder to eager-load the nodes that are connected to
// the "file" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *RequestQuery) WithFile(opts ...func(*FileQuery)) *RequestQuery {
	query := &FileQuery{config: rq.config}
	for _, opt := range opts {
		opt(query)
	}
	rq.withFile = query
	return rq
}

// WithTag tells the query-builder to eager-load the nodes that are connected to
// the "tag" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *RequestQuery) WithTag(opts ...func(*TagQuery)) *RequestQuery {
	query := &TagQuery{config: rq.config}
	for _, opt := range opts {
		opt(query)
	}
	rq.withTag = query
	return rq
}

// WithTransaction tells the query-builder to eager-load the nodes that are connected to
// the "transaction" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *RequestQuery) WithTransaction(opts ...func(*TransactionQuery)) *RequestQuery {
	query := &TransactionQuery{config: rq.config}
	for _, opt := range opts {
		opt(query)
	}
	rq.withTransaction = query
	return rq
}

// WithComment tells the query-builder to eager-load the nodes that are connected to
// the "comment" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *RequestQuery) WithComment(opts ...func(*CommentQuery)) *RequestQuery {
	query := &CommentQuery{config: rq.config}
	for _, opt := range opts {
		opt(query)
	}
	rq.withComment = query
	return rq
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *RequestQuery) WithUser(opts ...func(*UserQuery)) *RequestQuery {
	query := &UserQuery{config: rq.config}
	for _, opt := range opts {
		opt(query)
	}
	rq.withUser = query
	return rq
}

// WithGroup tells the query-builder to eager-load the nodes that are connected to
// the "group" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *RequestQuery) WithGroup(opts ...func(*GroupQuery)) *RequestQuery {
	query := &GroupQuery{config: rq.config}
	for _, opt := range opts {
		opt(query)
	}
	rq.withGroup = query
	return rq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Title string `json:"title,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Request.Query().
//		GroupBy(request.FieldTitle).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (rq *RequestQuery) GroupBy(field string, fields ...string) *RequestGroupBy {
	grbuild := &RequestGroupBy{config: rq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return rq.sqlQuery(ctx), nil
	}
	grbuild.label = request.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Title string `json:"title,omitempty"`
//	}
//
//	client.Request.Query().
//		Select(request.FieldTitle).
//		Scan(ctx, &v)
func (rq *RequestQuery) Select(fields ...string) *RequestSelect {
	rq.fields = append(rq.fields, fields...)
	selbuild := &RequestSelect{RequestQuery: rq}
	selbuild.label = request.Label
	selbuild.flds, selbuild.scan = &rq.fields, selbuild.Scan
	return selbuild
}

func (rq *RequestQuery) prepareQuery(ctx context.Context) error {
	for _, f := range rq.fields {
		if !request.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if rq.path != nil {
		prev, err := rq.path(ctx)
		if err != nil {
			return err
		}
		rq.sql = prev
	}
	return nil
}

func (rq *RequestQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Request, error) {
	var (
		nodes       = []*Request{}
		withFKs     = rq.withFKs
		_spec       = rq.querySpec()
		loadedTypes = [8]bool{
			rq.withStatus != nil,
			rq.withTarget != nil,
			rq.withFile != nil,
			rq.withTag != nil,
			rq.withTransaction != nil,
			rq.withComment != nil,
			rq.withUser != nil,
			rq.withGroup != nil,
		}
	)
	if rq.withUser != nil || rq.withGroup != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, request.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Request).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Request{config: rq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, rq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := rq.withStatus; query != nil {
		if err := rq.loadStatus(ctx, query, nodes,
			func(n *Request) { n.Edges.Status = []*RequestStatus{} },
			func(n *Request, e *RequestStatus) { n.Edges.Status = append(n.Edges.Status, e) }); err != nil {
			return nil, err
		}
	}
	if query := rq.withTarget; query != nil {
		if err := rq.loadTarget(ctx, query, nodes,
			func(n *Request) { n.Edges.Target = []*RequestTarget{} },
			func(n *Request, e *RequestTarget) { n.Edges.Target = append(n.Edges.Target, e) }); err != nil {
			return nil, err
		}
	}
	if query := rq.withFile; query != nil {
		if err := rq.loadFile(ctx, query, nodes,
			func(n *Request) { n.Edges.File = []*File{} },
			func(n *Request, e *File) { n.Edges.File = append(n.Edges.File, e) }); err != nil {
			return nil, err
		}
	}
	if query := rq.withTag; query != nil {
		if err := rq.loadTag(ctx, query, nodes,
			func(n *Request) { n.Edges.Tag = []*Tag{} },
			func(n *Request, e *Tag) { n.Edges.Tag = append(n.Edges.Tag, e) }); err != nil {
			return nil, err
		}
	}
	if query := rq.withTransaction; query != nil {
		if err := rq.loadTransaction(ctx, query, nodes,
			func(n *Request) { n.Edges.Transaction = []*Transaction{} },
			func(n *Request, e *Transaction) { n.Edges.Transaction = append(n.Edges.Transaction, e) }); err != nil {
			return nil, err
		}
	}
	if query := rq.withComment; query != nil {
		if err := rq.loadComment(ctx, query, nodes,
			func(n *Request) { n.Edges.Comment = []*Comment{} },
			func(n *Request, e *Comment) { n.Edges.Comment = append(n.Edges.Comment, e) }); err != nil {
			return nil, err
		}
	}
	if query := rq.withUser; query != nil {
		if err := rq.loadUser(ctx, query, nodes, nil,
			func(n *Request, e *User) { n.Edges.User = e }); err != nil {
			return nil, err
		}
	}
	if query := rq.withGroup; query != nil {
		if err := rq.loadGroup(ctx, query, nodes, nil,
			func(n *Request, e *Group) { n.Edges.Group = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (rq *RequestQuery) loadStatus(ctx context.Context, query *RequestStatusQuery, nodes []*Request, init func(*Request), assign func(*Request, *RequestStatus)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*Request)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.RequestStatus(func(s *sql.Selector) {
		s.Where(sql.InValues(request.StatusColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.request_status
		if fk == nil {
			return fmt.Errorf(`foreign-key "request_status" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "request_status" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (rq *RequestQuery) loadTarget(ctx context.Context, query *RequestTargetQuery, nodes []*Request, init func(*Request), assign func(*Request, *RequestTarget)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*Request)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.RequestTarget(func(s *sql.Selector) {
		s.Where(sql.InValues(request.TargetColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.request_target
		if fk == nil {
			return fmt.Errorf(`foreign-key "request_target" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "request_target" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (rq *RequestQuery) loadFile(ctx context.Context, query *FileQuery, nodes []*Request, init func(*Request), assign func(*Request, *File)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*Request)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.File(func(s *sql.Selector) {
		s.Where(sql.InValues(request.FileColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.request_file
		if fk == nil {
			return fmt.Errorf(`foreign-key "request_file" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "request_file" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (rq *RequestQuery) loadTag(ctx context.Context, query *TagQuery, nodes []*Request, init func(*Request), assign func(*Request, *Tag)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[uuid.UUID]*Request)
	nids := make(map[uuid.UUID]map[*Request]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(request.TagTable)
		s.Join(joinT).On(s.C(tag.FieldID), joinT.C(request.TagPrimaryKey[1]))
		s.Where(sql.InValues(joinT.C(request.TagPrimaryKey[0]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(request.TagPrimaryKey[0]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	neighbors, err := query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
		assign := spec.Assign
		values := spec.ScanValues
		spec.ScanValues = func(columns []string) ([]any, error) {
			values, err := values(columns[1:])
			if err != nil {
				return nil, err
			}
			return append([]any{new(uuid.UUID)}, values...), nil
		}
		spec.Assign = func(columns []string, values []any) error {
			outValue := *values[0].(*uuid.UUID)
			inValue := *values[1].(*uuid.UUID)
			if nids[inValue] == nil {
				nids[inValue] = map[*Request]struct{}{byID[outValue]: struct{}{}}
				return assign(columns[1:], values[1:])
			}
			nids[inValue][byID[outValue]] = struct{}{}
			return nil
		}
	})
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "tag" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}
func (rq *RequestQuery) loadTransaction(ctx context.Context, query *TransactionQuery, nodes []*Request, init func(*Request), assign func(*Request, *Transaction)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*Request)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Transaction(func(s *sql.Selector) {
		s.Where(sql.InValues(request.TransactionColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.request_transaction
		if fk == nil {
			return fmt.Errorf(`foreign-key "request_transaction" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "request_transaction" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (rq *RequestQuery) loadComment(ctx context.Context, query *CommentQuery, nodes []*Request, init func(*Request), assign func(*Request, *Comment)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*Request)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Comment(func(s *sql.Selector) {
		s.Where(sql.InValues(request.CommentColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.request_comment
		if fk == nil {
			return fmt.Errorf(`foreign-key "request_comment" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "request_comment" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (rq *RequestQuery) loadUser(ctx context.Context, query *UserQuery, nodes []*Request, init func(*Request), assign func(*Request, *User)) error {
	ids := make([]uuid.UUID, 0, len(nodes))
	nodeids := make(map[uuid.UUID][]*Request)
	for i := range nodes {
		if nodes[i].request_user == nil {
			continue
		}
		fk := *nodes[i].request_user
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	query.Where(user.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "request_user" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (rq *RequestQuery) loadGroup(ctx context.Context, query *GroupQuery, nodes []*Request, init func(*Request), assign func(*Request, *Group)) error {
	ids := make([]uuid.UUID, 0, len(nodes))
	nodeids := make(map[uuid.UUID][]*Request)
	for i := range nodes {
		if nodes[i].group_request == nil {
			continue
		}
		fk := *nodes[i].group_request
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	query.Where(group.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "group_request" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (rq *RequestQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := rq.querySpec()
	_spec.Node.Columns = rq.fields
	if len(rq.fields) > 0 {
		_spec.Unique = rq.unique != nil && *rq.unique
	}
	return sqlgraph.CountNodes(ctx, rq.driver, _spec)
}

func (rq *RequestQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := rq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (rq *RequestQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   request.Table,
			Columns: request.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: request.FieldID,
			},
		},
		From:   rq.sql,
		Unique: true,
	}
	if unique := rq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := rq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, request.FieldID)
		for i := range fields {
			if fields[i] != request.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := rq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := rq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := rq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := rq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (rq *RequestQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(rq.driver.Dialect())
	t1 := builder.Table(request.Table)
	columns := rq.fields
	if len(columns) == 0 {
		columns = request.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if rq.sql != nil {
		selector = rq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if rq.unique != nil && *rq.unique {
		selector.Distinct()
	}
	for _, p := range rq.predicates {
		p(selector)
	}
	for _, p := range rq.order {
		p(selector)
	}
	if offset := rq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := rq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// RequestGroupBy is the group-by builder for Request entities.
type RequestGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (rgb *RequestGroupBy) Aggregate(fns ...AggregateFunc) *RequestGroupBy {
	rgb.fns = append(rgb.fns, fns...)
	return rgb
}

// Scan applies the group-by query and scans the result into the given value.
func (rgb *RequestGroupBy) Scan(ctx context.Context, v any) error {
	query, err := rgb.path(ctx)
	if err != nil {
		return err
	}
	rgb.sql = query
	return rgb.sqlScan(ctx, v)
}

func (rgb *RequestGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range rgb.fields {
		if !request.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := rgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (rgb *RequestGroupBy) sqlQuery() *sql.Selector {
	selector := rgb.sql.Select()
	aggregation := make([]string, 0, len(rgb.fns))
	for _, fn := range rgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(rgb.fields)+len(rgb.fns))
		for _, f := range rgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(rgb.fields...)...)
}

// RequestSelect is the builder for selecting fields of Request entities.
type RequestSelect struct {
	*RequestQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (rs *RequestSelect) Scan(ctx context.Context, v any) error {
	if err := rs.prepareQuery(ctx); err != nil {
		return err
	}
	rs.sql = rs.RequestQuery.sqlQuery(ctx)
	return rs.sqlScan(ctx, v)
}

func (rs *RequestSelect) sqlScan(ctx context.Context, v any) error {
	rows := &sql.Rows{}
	query, args := rs.sql.Query()
	if err := rs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
