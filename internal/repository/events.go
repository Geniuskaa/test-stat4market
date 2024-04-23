package repository

import (
	"test-stat4market/internal/repository/entity"
)

type eventsQuery struct {
	BaseQuery
}

type EventsQuery interface {
	Insert(data *entity.Event) error
	ListByEventTypeAndDate(filter entity.EventsFilter) ([]entity.Event, error)
}

func (q *eventsQuery) Insert(data *entity.Event) error {
	_, err := q.runner.ExecContext(q.Context(),
		`INSERT INTO test.events (eventID, eventType, userID, eventTime, payload) VALUES (now64(),?,?,?,?);`,
		data.EventType, data.UserID, data.EventTime, data.Payload)
	if err != nil {
		return err
	}

	return err
}

func (q *eventsQuery) ListByEventTypeAndDate(filter entity.EventsFilter) ([]entity.Event, error) {
	rows, err := q.runner.QueryContext(q.Context(),
		`SELECT eventType, userID, eventTime, payload FROM test.events WHERE eventType = ? AND eventTime BETWEEN ? AND ?;`,
		filter.EventType, filter.EventsFrom, filter.EventsTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]entity.Event, 0)
	for rows.Next() {
		data := entity.Event{}
		err = rows.Scan(&data.EventType, &data.UserID, &data.EventTime, &data.Payload)
		if err != nil {
			return nil, err
		}
		res = append(res, data)
	}

	return res, err
}
