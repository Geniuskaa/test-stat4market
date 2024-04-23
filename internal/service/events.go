package service

import (
	"context"
	"time"

	"test-stat4market/internal/app/dto"
	"test-stat4market/internal/logger"
	"test-stat4market/internal/repository"
	"test-stat4market/internal/repository/entity"
)

type eventsService struct {
	dao repository.DAO
}

func NewEventsService(dao repository.DAO) EventsService {
	resp := eventsService{
		dao: dao,
	}
	return &resp
}

type EventsService interface {
	SaveEvent(ctx context.Context, req *dto.Event) error
	ListEvents(ctx context.Context, filter dto.EventsFilter) ([]dto.Event, error)
}

func (e *eventsService) SaveEvent(ctx context.Context, req *dto.Event) error {
	event := entity.Event{
		EventType: req.EventType,
		UserID:    req.UserID,
		Payload:   req.Payload,
	}

	t, err := time.Parse(time.DateTime, req.EventTime)
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{
			Msg:   "Failed to parse time",
			Error: err,
		})
	}

	event.EventTime = t

	err = e.dao.NewEventsQuery(ctx).Insert(&event)
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{
			Msg:   "Failed to save event",
			Error: err,
		})
		return err
	}

	return err
}

func (e *eventsService) ListEvents(ctx context.Context, filter dto.EventsFilter) ([]dto.Event, error) {
	from, err := time.Parse(time.DateTime, filter.EventsFrom)
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{
			Msg:   "Failed to parse time",
			Error: err,
		})
	}

	to, err := time.Parse(time.DateTime, filter.EventsTo)
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{
			Msg:   "Failed to parse time",
			Error: err,
		})
	}
	searchFilter := entity.EventsFilter{
		EventType:  filter.EventType,
		EventsFrom: from,
		EventsTo:   to,
	}

	resp, err := e.dao.NewEventsQuery(ctx).ListByEventTypeAndDate(searchFilter)
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{
			Msg:   "Failed to save event",
			Error: err,
		})
		return nil, err
	}

	convertedResp := make([]dto.Event, 0, len(resp))
	for _, val := range resp {
		convertedResp = append(convertedResp, dto.Event{
			EventType: val.EventType,
			UserID:    val.UserID,
			EventTime: val.EventTime.Format(time.RFC3339),
			Payload:   val.Payload,
		})
	}

	return convertedResp, err
}
