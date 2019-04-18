package handler

import (
	"github.com/labstack/echo"
	"github.com/nskondratev/api-page-go-back/pages"
	"github.com/nskondratev/api-page-go-back/ws"
	"net/http"
	"strconv"
)

func (h *Handler) GetPage(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	page, err := h.pageStore.GetById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	if page == nil {
		return c.JSON(http.StatusNotFound, &errorResponseEnvelope{
			Error: "Not found",
		})
	}
	return c.JSON(http.StatusOK, &responseEnvelope{
		Data: page,
	})
}

func (h *Handler) ListPages(c echo.Context) error {
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
	query := c.QueryParam("query")
	if len(sort) < 1 {
		sort = "createdAt"
		descending = true
	}
	pagesList, total, err := h.pageStore.List(offset, limit, sort, descending, query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, &paginationResponseEnvelope{
		Data:  pagesList,
		Total: total,
	})
}

func (h *Handler) CreatePage(c echo.Context) error {
	req := &pageCreateRequest{}
	page := &pages.Page{}
	if err := req.bind(c, page); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	if err := h.pageStore.Create(page); err != nil {
		return c.JSON(http.StatusInternalServerError, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	wsMessage := &ws.ApPageMessage{
		EventConst: ws.PageCreated,
		Data: &ws.ApMessagePageEnvelope{
			Page: page,
		},
	}
	if err := h.wsHub.Broadcast(wsMessage); err != nil {
		h.logger.Warnf("Error while broadcasting PAGE_CREATED to ws: %s", err.Error())
	}
	return c.JSON(http.StatusOK, &responseEnvelope{
		Data: page,
	})
}

func (h *Handler) UpdatePage(c echo.Context) error {
	req := &pageUpdateRequest{}
	page := &pages.Page{}
	if err := req.bind(c, page); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	if err := h.pageStore.Update(page); err != nil {
		return c.JSON(http.StatusInternalServerError, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	wsMessage := &ws.ApPageMessage{
		EventConst: ws.PageUpdated,
		Data: &ws.ApMessagePageEnvelope{
			Page: page,
		},
	}
	if err := h.wsHub.Broadcast(wsMessage); err != nil {
		h.logger.Warnf("Error while broadcasting PAGE_UPDATED to ws: %s", err.Error())
	}
	return c.JSON(http.StatusOK, &responseEnvelope{
		Data: page,
	})
}

func (h *Handler) DeletePage(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	p := &pages.Page{
		ID: id,
	}
	if err := h.pageStore.Delete(p); err != nil {
		return c.JSON(http.StatusInternalServerError, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	wsMessage := &ws.ApIdMessage{
		EventConst: ws.PageDeleted,
		Data: &ws.ApMessageOnlyIdEnvelope{
			ID: p.ID,
		},
	}
	if err := h.wsHub.Broadcast(wsMessage); err != nil {
		h.logger.Warnf("Error while broadcasting PAGE_DELETED to ws: %s", err.Error())
	}
	return c.NoContent(http.StatusOK)
}
