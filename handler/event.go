package handler

import (
	"github.com/labstack/echo"
	"github.com/nskondratev/api-page-go-back/events"
	"github.com/nskondratev/api-page-go-back/ws"
	"net/http"
	"strconv"
)

func (h *Handler) GetEvent(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	event, err := h.eventStore.GetById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	if event == nil {
		return c.JSON(http.StatusNotFound, &errorResponseEnvelope{
			Error: "Not found",
		})
	}
	return c.JSON(http.StatusOK, &responseEnvelope{
		Data: event,
	})
}

func (h *Handler) ListEvents(c echo.Context) error {
	sort := c.QueryParam("sort")
	descending := c.QueryParam("descending") == "true"
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 100
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}
	eType := c.QueryParam("type")
	query := c.QueryParam("query")
	if len(sort) < 1 {
		sort = "createdAt"
		descending = true
	}
	eventsList, total, err := h.eventStore.List(offset, limit, sort, descending, eType, query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, &paginationResponseEnvelope{
		Data:  eventsList,
		Total: total,
	})
}

func (h *Handler) CreateEvent(c echo.Context) error {
	req := &eventCreateRequest{}
	event := &events.Event{}
	if err := req.bind(c, event); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	if err := h.eventStore.Create(event); err != nil {
		return c.JSON(http.StatusInternalServerError, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	wsMessage := &ws.ApEventMessage{
		EventConst: ws.EventCreated,
		Data: &ws.ApMessageEventEnvelope{
			Event: event,
		},
	}
	if err := h.wsHub.Broadcast(wsMessage); err != nil {
		h.logger.Warnf("Error while broadcasting EVENT_CREATED to ws: %s", err.Error())
	}
	return c.JSON(http.StatusOK, &responseEnvelope{
		Data: event,
	})
}

func (h *Handler) UpdateEvent(c echo.Context) error {
	req := &eventUpdateRequest{}
	event := &events.Event{}
	if err := req.bind(c, event); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	if err := h.eventStore.Update(event); err != nil {
		return c.JSON(http.StatusInternalServerError, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	wsMessage := &ws.ApEventMessage{
		EventConst: ws.EventUpdated,
		Data: &ws.ApMessageEventEnvelope{
			Event: event,
		},
	}
	if err := h.wsHub.Broadcast(wsMessage); err != nil {
		h.logger.Warnf("Error while broadcasting EVENT_UPDATED to ws: %s", err.Error())
	}
	return c.JSON(http.StatusOK, &responseEnvelope{
		Data: event,
	})
}

func (h *Handler) DeleteEvent(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	e := &events.Event{
		ID: id,
	}
	if err := h.eventStore.Delete(e); err != nil {
		return c.JSON(http.StatusInternalServerError, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	wsMessage := &ws.ApIdMessage{
		EventConst: ws.EventDeleted,
		Data: &ws.ApMessageOnlyIdEnvelope{
			ID: e.ID,
		},
	}
	if err := h.wsHub.Broadcast(wsMessage); err != nil {
		h.logger.Warnf("Error while broadcasting EVENT_DELETED to ws: %s", err.Error())
	}
	return c.NoContent(http.StatusOK)
}
