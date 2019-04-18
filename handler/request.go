package handler

import (
	"github.com/labstack/echo"
	"github.com/nskondratev/api-page-go-back/events"
	"github.com/nskondratev/api-page-go-back/pages"
	"github.com/nskondratev/api-page-go-back/util"
	"strconv"
)

type pageUpdateRequest struct {
	Page struct {
		ID    uint64 `json:"id" validate:"required"`
		Title string `json:"title" validate:"required"`
		Text  string `json:"text" validate:"required"`
	} `json:"page"`
}

func (r *pageUpdateRequest) bind(c echo.Context, p *pages.Page) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	r.Page.ID = id
	if err := c.Validate(r); err != nil {
		return err
	}
	p.ID = r.Page.ID
	p.Title = r.Page.Title
	p.Text = r.Page.Text
	return nil
}

type pageCreateRequest struct {
	Page struct {
		Title string `json:"title" validate:"required"`
		Text  string `json:"text" validate:"required"`
	} `json:"page"`
}

func (r *pageCreateRequest) bind(c echo.Context, p *pages.Page) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	p.Title = r.Page.Title
	p.Text = r.Page.Text
	return nil
}

type eventCreateRequest struct {
	Event struct {
		Label       util.NullString `json:"label"`
		Constant    string          `json:"constant" validate:"required"`
		Value       string          `json:"value" validate:"required"`
		Description string          `json:"description" validate:"required"`
		Type        string          `json:"type" validate:"required"`
		Fields      []struct {
			Key         util.NullString `json:"key" validate:"required"`
			Type        string          `json:"type" validate:"required"`
			Required    bool            `json:"required" validate:"required"`
			Description string          `json:"description" validate:"required"`
		} `json:"fields"`
	} `json:"event"`
}

func (r *eventCreateRequest) bind(c echo.Context, e *events.Event) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	e.Label = r.Event.Label
	e.Constant = r.Event.Constant
	e.Value = r.Event.Value
	e.Description = r.Event.Description
	e.Type = r.Event.Type
	e.Fields = make([]events.Field, len(r.Event.Fields), len(r.Event.Fields))
	for index, element := range r.Event.Fields {
		f := events.Field{
			Key:         element.Key,
			Type:        element.Type,
			Required:    element.Required,
			Description: element.Description,
		}
		e.Fields[index] = f
	}
	return nil
}

type eventUpdateRequest struct {
	Event struct {
		ID          uint64          `json:"id" validate:"required"`
		Label       util.NullString `json:"label"`
		Constant    string          `json:"constant" validate:"required"`
		Value       string          `json:"value" validate:"required"`
		Description string          `json:"description" validate:"required"`
		Type        string          `json:"type" validate:"required"`
		Fields      []struct {
			Key         util.NullString `json:"key" validate:"required"`
			Type        string          `json:"type" validate:"required"`
			Required    bool            `json:"required" validate:"required"`
			Description string          `json:"description" validate:"required"`
		} `json:"fields"`
	} `json:"event"`
}

func (r *eventUpdateRequest) bind(c echo.Context, e *events.Event) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	r.Event.ID = id
	if err := c.Validate(r); err != nil {
		return err
	}
	e.ID = r.Event.ID
	e.Label = r.Event.Label
	e.Constant = r.Event.Constant
	e.Value = r.Event.Value
	e.Description = r.Event.Description
	e.Type = r.Event.Type
	e.Fields = make([]events.Field, len(r.Event.Fields), len(r.Event.Fields))
	for index, element := range r.Event.Fields {
		f := events.Field{
			Key:         element.Key,
			Type:        element.Type,
			Required:    element.Required,
			Description: element.Description,
		}
		e.Fields[index] = f
	}
	return nil
}
