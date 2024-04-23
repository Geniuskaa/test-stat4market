package repository

import (
	"context"

	"github.com/uptrace/go-clickhouse/ch"
)

type dao struct {
	RunnerWrapper
}

type DAO interface {
	NewEventsQuery(ctx context.Context) EventsQuery
}

func NewDAO(db *ch.DB) DAO {
	return &dao{db}
}

func (d *dao) NewEventsQuery(ctx context.Context) EventsQuery {
	return &eventsQuery{BaseQuery{
		ctx:    ctx,
		runner: d.RunnerWrapper,
	}}
}
