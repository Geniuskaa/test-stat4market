package dto

type Event struct {
	EventType string `json:"eventType"`
	UserID    int    `json:"userID"`
	EventTime string `json:"eventTime"`
	Payload   string `json:"payload"`
}

type EventsFilter struct {
	EventType  string `json:"event_type"`
	EventsFrom string `json:"events_from"`
	EventsTo   string `json:"events_to"`
}
