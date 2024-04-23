package repository

import (
	"context"
	"database/sql"

	"github.com/uptrace/go-clickhouse/ch"
	"github.com/uptrace/go-clickhouse/ch/chschema"
)

type BaseQuery struct {
	ctx    context.Context
	runner RunnerWrapper
}

func (q *BaseQuery) Context() context.Context {
	return q.ctx
}

func (q *BaseQuery) Runner() RunnerWrapper {
	return q.runner
}

type RunnerWrapper interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*ch.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*ch.Rows, error)
	QueryRow(query string, args ...any) *ch.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *ch.Row
	NewSelect() *ch.SelectQuery
	NewRaw(query string, args ...any) *ch.RawQuery
	NewInsert() *ch.InsertQuery
	NewCreateTable() *ch.CreateTableQuery
	NewDropTable() *ch.DropTableQuery
	NewTruncateTable() *ch.TruncateTableQuery
	NewCreateView() *ch.CreateViewQuery
	NewDropView() *ch.DropViewQuery
	ResetModel(ctx context.Context, models ...any) error
	Formatter() chschema.Formatter
	WithFormatter(fmter chschema.Formatter) *ch.DB
	FormatQuery(query string, args ...any) string
}
