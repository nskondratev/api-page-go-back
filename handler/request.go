package handler

import (
	"github.com/labstack/echo"
	"github.com/nskondratev/api-page-go-back/events"
	"github.com/nskondratev/api-page-go-back/pages"
	"github.com/nskondratev/api-page-go-back/util"
	"strconv"
)

type pageUpdateRequest struct {
	ID    uint64 `json:"id" validate:"required"`
	Title string `json:"title" validate:"required"`
	Text  string `json:"text" validate:"required"`
}

func (r *pageUpdateRequest) bind(c echo.Context, p *pages.Page) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	r.ID = id
	if err := c.Validate(r); err != nil {
		return err
	}
	p.ID = r.ID
	p.Title = r.Title
	p.Text = r.Text
	return nil
}

type pageCreateRequest struct {
	Title string `json:"title" validate:"required"`
	Text  string `json:"text" validate:"required"`
}

func (r *pageCreateRequest) bind(c echo.Context, p *pages.Page) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	p.Title = r.Title
	p.Text = r.Text
	return nil
}

type fieldsRequest struct {
	Key         util.NullString `json:"key" validate:"required"`
	Type        string          `json:"type" validate:"required"`
	Required    bool            `json:"required" validate:"required"`
	Description string          `json:"description" validate:"required"`
}

type eventCreateRequest struct {
	Label       util.NullString `json:"label"`
	Constant    string          `json:"constant" validate:"required"`
	Value       string          `json:"value" validate:"required"`
	Description string          `json:"description" validate:"required"`
	Type        string          `json:"type" validate:"required"`
	Fields      []fieldsRequest `json:"fields"`
}

func (r *eventCreateRequest) bind(c echo.Context, e *events.Event) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	e.Label = r.Label
	e.Constant = r.Constant
	e.Value = r.Value
	e.Description = r.Description
	e.Type = r.Type
	e.Fields = make([]events.Field, len(r.Fields), len(r.Fields))
	for index, element := range r.Fields {
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
	ID          uint64          `json:"id" validate:"required"`
	Label       util.NullString `json:"label"`
	Constant    string          `json:"constant" validate:"required"`
	Value       string          `json:"value" validate:"required"`
	Description string          `json:"description" validate:"required"`
	Type        string          `json:"type" validate:"required"`
	Fields      []fieldsRequest `json:"fields"`
}

func (r *eventUpdateRequest) bind(c echo.Context, e *events.Event) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	r.ID = id
	if err := c.Validate(r); err != nil {
		return err
	}
	e.ID = r.ID
	e.Label = r.Label
	e.Constant = r.Constant
	e.Value = r.Value
	e.Description = r.Description
	e.Type = r.Type
	e.Fields = make([]events.Field, len(r.Fields), len(r.Fields))
	for index, element := range r.Fields {
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
