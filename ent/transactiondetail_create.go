// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent/transaction"
	"github.com/traPtitech/Jomon/ent/transactiondetail"
)

// TransactionDetailCreate is the builder for creating a TransactionDetail entity.
type TransactionDetailCreate struct {
	config
	mutation *TransactionDetailMutation
	hooks    []Hook
}

// SetAmount sets the "amount" field.
func (tdc *TransactionDetailCreate) SetAmount(i int) *TransactionDetailCreate {
	tdc.mutation.SetAmount(i)
	return tdc
}

// SetNillableAmount sets the "amount" field if the given value is not nil.
func (tdc *TransactionDetailCreate) SetNillableAmount(i *int) *TransactionDetailCreate {
	if i != nil {
		tdc.SetAmount(*i)
	}
	return tdc
}

// SetTarget sets the "target" field.
func (tdc *TransactionDetailCreate) SetTarget(s string) *TransactionDetailCreate {
	tdc.mutation.SetTarget(s)
	return tdc
}

// SetNillableTarget sets the "target" field if the given value is not nil.
func (tdc *TransactionDetailCreate) SetNillableTarget(s *string) *TransactionDetailCreate {
	if s != nil {
		tdc.SetTarget(*s)
	}
	return tdc
}

// SetCreatedAt sets the "created_at" field.
func (tdc *TransactionDetailCreate) SetCreatedAt(t time.Time) *TransactionDetailCreate {
	tdc.mutation.SetCreatedAt(t)
	return tdc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (tdc *TransactionDetailCreate) SetNillableCreatedAt(t *time.Time) *TransactionDetailCreate {
	if t != nil {
		tdc.SetCreatedAt(*t)
	}
	return tdc
}

// SetUpdatedAt sets the "updated_at" field.
func (tdc *TransactionDetailCreate) SetUpdatedAt(t time.Time) *TransactionDetailCreate {
	tdc.mutation.SetUpdatedAt(t)
	return tdc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (tdc *TransactionDetailCreate) SetNillableUpdatedAt(t *time.Time) *TransactionDetailCreate {
	if t != nil {
		tdc.SetUpdatedAt(*t)
	}
	return tdc
}

// SetID sets the "id" field.
func (tdc *TransactionDetailCreate) SetID(u uuid.UUID) *TransactionDetailCreate {
	tdc.mutation.SetID(u)
	return tdc
}

// SetTransactionID sets the "transaction" edge to the Transaction entity by ID.
func (tdc *TransactionDetailCreate) SetTransactionID(id uuid.UUID) *TransactionDetailCreate {
	tdc.mutation.SetTransactionID(id)
	return tdc
}

// SetTransaction sets the "transaction" edge to the Transaction entity.
func (tdc *TransactionDetailCreate) SetTransaction(t *Transaction) *TransactionDetailCreate {
	return tdc.SetTransactionID(t.ID)
}

// Mutation returns the TransactionDetailMutation object of the builder.
func (tdc *TransactionDetailCreate) Mutation() *TransactionDetailMutation {
	return tdc.mutation
}

// Save creates the TransactionDetail in the database.
func (tdc *TransactionDetailCreate) Save(ctx context.Context) (*TransactionDetail, error) {
	var (
		err  error
		node *TransactionDetail
	)
	tdc.defaults()
	if len(tdc.hooks) == 0 {
		if err = tdc.check(); err != nil {
			return nil, err
		}
		node, err = tdc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TransactionDetailMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = tdc.check(); err != nil {
				return nil, err
			}
			tdc.mutation = mutation
			if node, err = tdc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(tdc.hooks) - 1; i >= 0; i-- {
			if tdc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = tdc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tdc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (tdc *TransactionDetailCreate) SaveX(ctx context.Context) *TransactionDetail {
	v, err := tdc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tdc *TransactionDetailCreate) Exec(ctx context.Context) error {
	_, err := tdc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tdc *TransactionDetailCreate) ExecX(ctx context.Context) {
	if err := tdc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tdc *TransactionDetailCreate) defaults() {
	if _, ok := tdc.mutation.Amount(); !ok {
		v := transactiondetail.DefaultAmount
		tdc.mutation.SetAmount(v)
	}
	if _, ok := tdc.mutation.Target(); !ok {
		v := transactiondetail.DefaultTarget
		tdc.mutation.SetTarget(v)
	}
	if _, ok := tdc.mutation.CreatedAt(); !ok {
		v := transactiondetail.DefaultCreatedAt()
		tdc.mutation.SetCreatedAt(v)
	}
	if _, ok := tdc.mutation.UpdatedAt(); !ok {
		v := transactiondetail.DefaultUpdatedAt()
		tdc.mutation.SetUpdatedAt(v)
	}
	if _, ok := tdc.mutation.ID(); !ok {
		v := transactiondetail.DefaultID()
		tdc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tdc *TransactionDetailCreate) check() error {
	if _, ok := tdc.mutation.Amount(); !ok {
		return &ValidationError{Name: "amount", err: errors.New(`ent: missing required field "amount"`)}
	}
	if _, ok := tdc.mutation.Target(); !ok {
		return &ValidationError{Name: "target", err: errors.New(`ent: missing required field "target"`)}
	}
	if _, ok := tdc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "created_at"`)}
	}
	if _, ok := tdc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "updated_at"`)}
	}
	if _, ok := tdc.mutation.TransactionID(); !ok {
		return &ValidationError{Name: "transaction", err: errors.New("ent: missing required edge \"transaction\"")}
	}
	return nil
}

func (tdc *TransactionDetailCreate) sqlSave(ctx context.Context) (*TransactionDetail, error) {
	_node, _spec := tdc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tdc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		_node.ID = _spec.ID.Value.(uuid.UUID)
	}
	return _node, nil
}

func (tdc *TransactionDetailCreate) createSpec() (*TransactionDetail, *sqlgraph.CreateSpec) {
	var (
		_node = &TransactionDetail{config: tdc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: transactiondetail.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: transactiondetail.FieldID,
			},
		}
	)
	if id, ok := tdc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := tdc.mutation.Amount(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: transactiondetail.FieldAmount,
		})
		_node.Amount = value
	}
	if value, ok := tdc.mutation.Target(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: transactiondetail.FieldTarget,
		})
		_node.Target = value
	}
	if value, ok := tdc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: transactiondetail.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := tdc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: transactiondetail.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if nodes := tdc.mutation.TransactionIDs(); len(nodes) > 0 {
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
		_node.transaction_detail = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// TransactionDetailCreateBulk is the builder for creating many TransactionDetail entities in bulk.
type TransactionDetailCreateBulk struct {
	config
	builders []*TransactionDetailCreate
}

// Save creates the TransactionDetail entities in the database.
func (tdcb *TransactionDetailCreateBulk) Save(ctx context.Context) ([]*TransactionDetail, error) {
	specs := make([]*sqlgraph.CreateSpec, len(tdcb.builders))
	nodes := make([]*TransactionDetail, len(tdcb.builders))
	mutators := make([]Mutator, len(tdcb.builders))
	for i := range tdcb.builders {
		func(i int, root context.Context) {
			builder := tdcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TransactionDetailMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, tdcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tdcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, tdcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tdcb *TransactionDetailCreateBulk) SaveX(ctx context.Context) []*TransactionDetail {
	v, err := tdcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tdcb *TransactionDetailCreateBulk) Exec(ctx context.Context) error {
	_, err := tdcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tdcb *TransactionDetailCreateBulk) ExecX(ctx context.Context) {
	if err := tdcb.Exec(ctx); err != nil {
		panic(err)
	}
}
