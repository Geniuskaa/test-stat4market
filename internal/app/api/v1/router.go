package v1

import (
	"github.com/gofiber/fiber/v2"
	"test-stat4market/internal/app/dto"
	"test-stat4market/internal/service"
)

type Router struct {
	fiber.Router
	eventService service.EventsService
}

func NewRoute(router fiber.Router, eventService service.EventsService) *Router {
	r := Router{
		Router:       router,
		eventService: eventService,
	}

	return &r
}

func (r *Router) Routes() {
	r.Post("/event", r.saveEvent)
	r.Post("/events", r.listEvents)

}

// Сохраняет полученное событие
func (r *Router) saveEvent(ctx *fiber.Ctx) error {

	req := dto.Event{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	err = r.eventService.SaveEvent(ctx.Context(), &req)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusAccepted)
}

// Возвращает события полученные поле применения фильтра
func (r *Router) listEvents(ctx *fiber.Ctx) error {

	req := dto.EventsFilter{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	resp, err := r.eventService.ListEvents(ctx.Context(), req)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	ctx.Status(fiber.StatusOK)
	return ctx.JSON(resp)
}
