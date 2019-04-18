package handler

import (
	"github.com/labstack/echo"
	"net/http"
)

type graphqlRequest struct {
	Query string `json:"query" validate:"required"`
}

func (r *graphqlRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	return nil
}

func (h *Handler) handleGraphQLQuery(c echo.Context) error {
	gq := &graphqlRequest{}
	if err := gq.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	result, err := h.gqlHub.Execute(gq.Query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &errorResponseEnvelope{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, result)
}
