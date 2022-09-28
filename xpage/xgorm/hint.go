package xgorm

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Exprs []clause.Expression

func (exprs Exprs) Build(builder clause.Builder) {
	for idx, expr := range exprs {
		if idx > 0 {
			builder.WriteByte(' ')
		}
		expr.Build(builder)
	}
}

type Hints struct {
	Prefix  string
	Suffix  string
	Content string

	Clauses []string
	Before  bool
	After   bool
}

func (hints Hints) ModifyStatement(stmt *gorm.Statement) {
	for _, name := range hints.Clauses {
		name = strings.ToUpper(name)
		clause := stmt.Clauses[name]
		switch {
		case hints.Before:
			if clause.BeforeExpression == nil {
				clause.BeforeExpression = hints
			} else if old, ok := clause.BeforeExpression.(Hints); ok {
				old.Merge(hints)
				clause.BeforeExpression = old
			} else {
				clause.BeforeExpression = Exprs{clause.BeforeExpression, hints}
			}
		case hints.After:
			if clause.AfterExpression == nil {
				clause.AfterExpression = hints
			} else if old, ok := clause.AfterExpression.(Hints); ok {
				old.Merge(hints)
				clause.AfterExpression = old
			} else {
				clause.AfterExpression = Exprs{clause.AfterExpression, hints}
			}
		default:
			if clause.AfterNameExpression == nil {
				clause.AfterNameExpression = hints
			} else if old, ok := clause.AfterNameExpression.(Hints); ok {
				old.Merge(hints)
				clause.AfterNameExpression = old
			} else {
				clause.AfterNameExpression = Exprs{clause.AfterNameExpression, hints}
			}
		}

		stmt.Clauses[name] = clause
	}
}

func (hints Hints) Build(builder clause.Builder) {
	builder.WriteString(hints.Prefix)
	builder.WriteString(hints.Content)
	builder.WriteString(hints.Suffix)
}

func (hints Hints) Merge(h Hints) {
	hints.Content += " " + h.Content
}
