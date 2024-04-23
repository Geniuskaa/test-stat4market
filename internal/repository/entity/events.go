package entity

import "time"

type Event struct {
	EventType string    `json:"eventType"`
	UserID    int       `json:"userID"`
	EventTime time.Time `json:"eventTime"`
	Payload   string    `json:"payload"`
}

type EventsFilter struct {
	EventType  string
	EventsFrom time.Time
	EventsTo   time.Time
}
