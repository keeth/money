// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package money

import (
	"database/sql"
)

type Acc struct {
	ID        int64
	CreatedAt string
	UpdatedAt string
	Name      string
	Xid       string
	Kind      string
	IsActive  int64
}

type Cat struct {
	ID        int64
	CreatedAt string
	UpdatedAt string
	Name      string
	Kind      string
	IsActive  int64
}

type Plan struct {
	ID         int64
	CreatedAt  string
	UpdatedAt  string
	StartDate  sql.NullString
	EndDate    sql.NullString
	CatID      int64
	AmountExpr string
	Period     string
}

type PlanPeriod struct {
	CreatedAt   string
	UpdatedAt   string
	PlanID      int64
	PeriodStart string
	PeriodEnd   string
	Amount      float64
}

type Rule struct {
	ID         int64
	CreatedAt  string
	UpdatedAt  string
	StartDate  sql.NullString
	EndDate    sql.NullString
	TestExpr   string
	CatID      sql.NullInt64
	AmountExpr sql.NullString
	DescExpr   sql.NullString
	DateExpr   sql.NullString
	Ord        int64
}

type Tx struct {
	ID         int64
	CreatedAt  string
	UpdatedAt  string
	Xid        string
	Date       string
	Desc       string
	Amount     float64
	OrigDate   sql.NullString
	OrigDesc   sql.NullString
	OrigAmount sql.NullFloat64
	AccID      int64
	CatID      sql.NullInt64
	Ord        string
}
