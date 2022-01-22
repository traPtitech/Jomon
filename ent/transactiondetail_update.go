// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent/predicate"
	"github.com/traPtitech/Jomon/ent/transaction"
	"github.com/traPtitech/Jomon/ent/transactiondetail"
)

// TransactionDetailUpdate is the builder for updating TransactionDetail entities.
type TransactionDetailUpdate struct {
	config
	hooks    []Hook
	mutation *TransactionDetailMutation
}

// Where appends a list predicates to the TransactionDetailUpdate builder.
func (tdu *TransactionDetailUpdate) Where(ps ...predicate.TransactionDetail) *TransactionDetailUpdate {
	tdu.mutation.Where(ps...)
	return tdu
}

// SetAmount sets the "amount" field.
func (tdu *TransactionDetailUpdate) SetAmount(i int) *TransactionDetailUpdate {
	tdu.mutation.ResetAmount()
	tdu.mutation.SetAmount(i)
	return tdu
}

// SetNillableAmount sets the "amount" field if the given value is not nil.
func (tdu *TransactionDetailUpdate) SetNillableAmount(i *int) *TransactionDetailUpdate {
	if i != nil {
		tdu.SetAmount(*i)
	}
	return tdu
}

// AddAmount adds i to the "amount" field.
func (tdu *TransactionDetailUpdate) AddAmount(i int) *TransactionDetailUpdate {
	tdu.mutation.AddAmount(i)
	return tdu
}

// SetTarget sets the "target" field.
func (tdu *TransactionDetailUpdate) SetTarget(s string) *TransactionDetailUpdate {
	tdu.mutation.SetTarget(s)
	return tdu
}

// SetNillableTarget sets the "target" field if the given value is not nil.
func (tdu *TransactionDetailUpdate) SetNillableTarget(s *string) *TransactionDetailUpdate {
	if s != nil {
		tdu.SetTarget(*s)
	}
	return tdu
}

// SetCreatedAt sets the "created_at" field.
func (tdu *TransactionDetailUpdate) SetCreatedAt(t time.Time) *TransactionDetailUpdate {
	tdu.mutation.SetCreatedAt(t)
	return tdu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (tdu *TransactionDetailUpdate) SetNillableCreatedAt(t *time.Time) *TransactionDetailUpdate {
	if t != nil {
		tdu.SetCreatedAt(*t)
	}
	return tdu
}

// SetUpdatedAt sets the "updated_at" field.
func (tdu *TransactionDetailUpdate) SetUpdatedAt(t time.Time) *TransactionDetailUpdate {
	tdu.mutation.SetUpdatedAt(t)
	return tdu
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (tdu *TransactionDetailUpdate) SetNillableUpdatedAt(t *time.Time) *TransactionDetailUpdate {
	if t != nil {
		tdu.SetUpdatedAt(*t)
	}
	return tdu
}

// SetTransactionID sets the "transaction" edge to the Transaction entity by ID.
func (tdu *TransactionDetailUpdate) SetTransactionID(id uuid.UUID) *TransactionDetailUpdate {
	tdu.mutation.SetTransactionID(id)
	return tdu
}

// SetTransaction sets the "transaction" edge to the Transaction entity.
func (tdu *TransactionDetailUpdate) SetTransaction(t *Transaction) *TransactionDetailUpdate {
	return tdu.SetTransactionID(t.ID)
}

// Mutation returns the TransactionDetailMutation object of the builder.
func (tdu *TransactionDetailUpdate) Mutation() *TransactionDetailMutation {
	return tdu.mutation
}

// ClearTransaction clears the "transaction" edge to the Transaction entity.
func (tdu *TransactionDetailUpdate) ClearTransaction() *TransactionDetailUpdate {
	tdu.mutation.ClearTransaction()
	return tdu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tdu *TransactionDetailUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(tdu.hooks) == 0 {
		if err = tdu.check(); err != nil {
			return 0, err
		}
		affected, err = tdu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TransactionDetailMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = tdu.check(); err != nil {
				return 0, err
			}
			tdu.mutation = mutation
			affected, err = tdu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(tdu.hooks) - 1; i >= 0; i-- {
			if tdu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = tdu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tdu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (tdu *TransactionDetailUpdate) SaveX(ctx context.Context) int {
	affected, err := tdu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tdu *TransactionDetailUpdate) Exec(ctx context.Context) error {
	_, err := tdu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tdu *TransactionDetailUpdate) ExecX(ctx context.Context) {
	if err := tdu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tdu *TransactionDetailUpdate) check() error {
	if _, ok := tdu.mutation.TransactionID(); tdu.mutation.TransactionCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "TransactionDetail.transaction"`)
	}
	return nil
}

func (tdu *TransactionDetailUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   transactiondetail.Table,
			Columns: transactiondetail.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: transactiondetail.FieldID,
			},
		},
	}
	if ps := tdu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tdu.mutation.Amount(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: transactiondetail.FieldAmount,
		})
	}
	if value, ok := tdu.mutation.AddedAmount(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: transactiondetail.FieldAmount,
		})
	}
	if value, ok := tdu.mutation.Target(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: transactiondetail.FieldTarget,
		})
	}
	if value, ok := tdu.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: transactiondetail.FieldCreatedAt,
		})
	}
	if value, ok := tdu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: transactiondetail.FieldUpdatedAt,
		})
	}
	if tdu.mutation.TransactionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   transactiondetail.TransactionTable,
			Columns: []string{transactiondetail.TransactionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: transaction.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tdu.mutation.TransactionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   transactiondetail.TransactionTable,
			Columns: []string{transactiondetail.TransactionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: transaction.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tdu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{transactiondetail.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// TransactionDetailUpdateOne is the builder for updating a single TransactionDetail entity.
type TransactionDetailUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *TransactionDetailMutation
}

// SetAmount sets the "amount" field.
func (tduo *TransactionDetailUpdateOne) SetAmount(i int) *TransactionDetailUpdateOne {
	tduo.mutation.ResetAmount()
	tduo.mutation.SetAmount(i)
	return tduo
}

// SetNillableAmount sets the "amount" field if the given value is not nil.
func (tduo *TransactionDetailUpdateOne) SetNillableAmount(i *int) *TransactionDetailUpdateOne {
	if i != nil {
		tduo.SetAmount(*i)
	}
	return tduo
}

// AddAmount adds i to the "amount" field.
func (tduo *TransactionDetailUpdateOne) AddAmount(i int) *TransactionDetailUpdateOne {
	tduo.mutation.AddAmount(i)
	return tduo
}

// SetTarget sets the "target" field.
func (tduo *TransactionDetailUpdateOne) SetTarget(s string) *TransactionDetailUpdateOne {
	tduo.mutation.SetTarget(s)
	return tduo
}

// SetNillableTarget sets the "target" field if the given value is not nil.
func (tduo *TransactionDetailUpdateOne) SetNillableTarget(s *string) *TransactionDetailUpdateOne {
	if s != nil {
		tduo.SetTarget(*s)
	}
	return tduo
}

// SetCreatedAt sets the "created_at" field.
func (tduo *TransactionDetailUpdateOne) SetCreatedAt(t time.Time) *TransactionDetailUpdateOne {
	tduo.mutation.SetCreatedAt(t)
	return tduo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (tduo *TransactionDetailUpdateOne) SetNillableCreatedAt(t *time.Time) *TransactionDetailUpdateOne {
	if t != nil {
		tduo.SetCreatedAt(*t)
	}
	return tduo
}

// SetUpdatedAt sets the "updated_at" field.
func (tduo *TransactionDetailUpdateOne) SetUpdatedAt(t time.Time) *TransactionDetailUpdateOne {
	tduo.mutation.SetUpdatedAt(t)
	return tduo
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (tduo *TransactionDetailUpdateOne) SetNillableUpdatedAt(t *time.Time) *TransactionDetailUpdateOne {
	if t != nil {
		tduo.SetUpdatedAt(*t)
	}
	return tduo
}

// SetTransactionID sets the "transaction" edge to the Transaction entity by ID.
func (tduo *TransactionDetailUpdateOne) SetTransactionID(id uuid.UUID) *TransactionDetailUpdateOne {
	tduo.mutation.SetTransactionID(id)
	return tduo
}

// SetTransaction sets the "transaction" edge to the Transaction entity.
func (tduo *TransactionDetailUpdateOne) SetTransaction(t *Transaction) *TransactionDetailUpdateOne {
	return tduo.SetTransactionID(t.ID)
}

// Mutation returns the TransactionDetailMutation object of the builder.
func (tduo *TransactionDetailUpdateOne) Mutation() *TransactionDetailMutation {
	return tduo.mutation
}

// ClearTransaction clears the "transaction" edge to the Transaction entity.
func (tduo *TransactionDetailUpdateOne) ClearTransaction() *TransactionDetailUpdateOne {
	tduo.mutation.ClearTransaction()
	return tduo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tduo *TransactionDetailUpdateOne) Select(field string, fields ...string) *TransactionDetailUpdateOne {
	tduo.fields = append([]string{field}, fields...)
	return tduo
}

// Save executes the query and returns the updated TransactionDetail entity.
func (tduo *TransactionDetailUpdateOne) Save(ctx context.Context) (*TransactionDetail, error) {
	var (
		err  error
		node *TransactionDetail
	)
	if len(tduo.hooks) == 0 {
		if err = tduo.check(); err != nil {
			return nil, err
		}
		node, err = tduo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TransactionDetailMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = tduo.check(); err != nil {
				return nil, err
			}
			tduo.mutation = mutation
			node, err = tduo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(tduo.hooks) - 1; i >= 0; i-- {
			if tduo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = tduo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tduo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (tduo *TransactionDetailUpdateOne) SaveX(ctx context.Context) *TransactionDetail {
	node, err := tduo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tduo *TransactionDetailUpdateOne) Exec(ctx context.Context) error {
	_, err := tduo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tduo *TransactionDetailUpdateOne) ExecX(ctx context.Context) {
	if err := tduo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tduo *TransactionDetailUpdateOne) check() error {
	if _, ok := tduo.mutation.TransactionID(); tduo.mutation.TransactionCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "TransactionDetail.transaction"`)
	}
	return nil
}

func (tduo *TransactionDetailUpdateOne) sqlSave(ctx context.Context) (_node *TransactionDetail, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   transactiondetail.Table,
			Columns: transactiondetail.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: transactiondetail.FieldID,
			},
		},
	}
	id, ok := tduo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "TransactionDetail.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := tduo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, transactiondetail.FieldID)
		for _, f := range fields {
			if !transactiondetail.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != transactiondetail.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tduo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tduo.mutation.Amount(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: transactiondetail.FieldAmount,
		})
	}
	if value, ok := tduo.mutation.AddedAmount(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: transactiondetail.FieldAmount,
		})
	}
	if value, ok := tduo.mutation.Target(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: transactiondetail.FieldTarget,
		})
	}
	if value, ok := tduo.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: transactiondetail.FieldCreatedAt,
		})
	}
	if value, ok := tduo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: transactiondetail.FieldUpdatedAt,
		})
	}
	if tduo.mutation.TransactionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   transactiondetail.TransactionTable,
			Columns: []string{transactiondetail.TransactionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: transaction.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tduo.mutation.TransactionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   transactiondetail.TransactionTable,
			Columns: []string{transactiondetail.TransactionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: transaction.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &TransactionDetail{config: tduo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tduo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{transactiondetail.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
