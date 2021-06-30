// Code generated by entc, DO NOT EDIT.

package groupbudget

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Amount applies equality check predicate on the "amount" field. It's identical to AmountEQ.
func Amount(v int) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldAmount), v))
	})
}

// Comment applies equality check predicate on the "comment" field. It's identical to CommentEQ.
func Comment(v string) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldComment), v))
	})
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// AmountEQ applies the EQ predicate on the "amount" field.
func AmountEQ(v int) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldAmount), v))
	})
}

// AmountNEQ applies the NEQ predicate on the "amount" field.
func AmountNEQ(v int) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldAmount), v))
	})
}

// AmountIn applies the In predicate on the "amount" field.
func AmountIn(vs ...int) predicate.GroupBudget {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.GroupBudget(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldAmount), v...))
	})
}

// AmountNotIn applies the NotIn predicate on the "amount" field.
func AmountNotIn(vs ...int) predicate.GroupBudget {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.GroupBudget(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldAmount), v...))
	})
}

// AmountGT applies the GT predicate on the "amount" field.
func AmountGT(v int) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldAmount), v))
	})
}

// AmountGTE applies the GTE predicate on the "amount" field.
func AmountGTE(v int) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldAmount), v))
	})
}

// AmountLT applies the LT predicate on the "amount" field.
func AmountLT(v int) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldAmount), v))
	})
}

// AmountLTE applies the LTE predicate on the "amount" field.
func AmountLTE(v int) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldAmount), v))
	})
}

// CommentEQ applies the EQ predicate on the "comment" field.
func CommentEQ(v string) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldComment), v))
	})
}

// CommentNEQ applies the NEQ predicate on the "comment" field.
func CommentNEQ(v string) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldComment), v))
	})
}

// CommentIn applies the In predicate on the "comment" field.
func CommentIn(vs ...string) predicate.GroupBudget {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.GroupBudget(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldComment), v...))
	})
}

// CommentNotIn applies the NotIn predicate on the "comment" field.
func CommentNotIn(vs ...string) predicate.GroupBudget {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.GroupBudget(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldComment), v...))
	})
}

// CommentGT applies the GT predicate on the "comment" field.
func CommentGT(v string) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldComment), v))
	})
}

// CommentGTE applies the GTE predicate on the "comment" field.
func CommentGTE(v string) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldComment), v))
	})
}

// CommentLT applies the LT predicate on the "comment" field.
func CommentLT(v string) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldComment), v))
	})
}

// CommentLTE applies the LTE predicate on the "comment" field.
func CommentLTE(v string) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldComment), v))
	})
}

// CommentContains applies the Contains predicate on the "comment" field.
func CommentContains(v string) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldComment), v))
	})
}

// CommentHasPrefix applies the HasPrefix predicate on the "comment" field.
func CommentHasPrefix(v string) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldComment), v))
	})
}

// CommentHasSuffix applies the HasSuffix predicate on the "comment" field.
func CommentHasSuffix(v string) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldComment), v))
	})
}

// CommentIsNil applies the IsNil predicate on the "comment" field.
func CommentIsNil() predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldComment)))
	})
}

// CommentNotNil applies the NotNil predicate on the "comment" field.
func CommentNotNil() predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldComment)))
	})
}

// CommentEqualFold applies the EqualFold predicate on the "comment" field.
func CommentEqualFold(v string) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldComment), v))
	})
}

// CommentContainsFold applies the ContainsFold predicate on the "comment" field.
func CommentContainsFold(v string) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldComment), v))
	})
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.GroupBudget {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.GroupBudget(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.GroupBudget {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.GroupBudget(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreatedAt), v))
	})
}

// HasGroup applies the HasEdge predicate on the "group" edge.
func HasGroup() predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(GroupTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, GroupTable, GroupColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasGroupWith applies the HasEdge predicate on the "group" edge with a given conditions (other predicates).
func HasGroupWith(preds ...predicate.Group) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(GroupInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, GroupTable, GroupColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasTransaction applies the HasEdge predicate on the "transaction" edge.
func HasTransaction() predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TransactionTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, TransactionTable, TransactionColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTransactionWith applies the HasEdge predicate on the "transaction" edge with a given conditions (other predicates).
func HasTransactionWith(preds ...predicate.Transaction) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TransactionInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, TransactionTable, TransactionColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.GroupBudget) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.GroupBudget) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.GroupBudget) predicate.GroupBudget {
	return predicate.GroupBudget(func(s *sql.Selector) {
		p(s.Not())
	})
}