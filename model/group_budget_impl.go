package model

import "github.com/traPtitech/Jomon/ent"

func convertEntGroupBudgetToModelGroupBudget(entgb *ent.GroupBudget) *GroupBudget {
	if entgb == nil {
		return nil
	}
	return &GroupBudget{
		ID:        entgb.ID,
		Amount:    entgb.Amount,
		Comment:   entgb.Comment,
		CreatedAt: entgb.CreatedAt,
	}
}
