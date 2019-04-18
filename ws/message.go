package ws

import (
	"encoding/json"
	"github.com/nskondratev/api-page-go-back/events"
	"github.com/nskondratev/api-page-go-back/pages"
)

type ApMessage interface {
	BuildWsMessage() ([]byte, error)
}

type ApMessageEventEnvelope struct {
	Event *events.Event `json:"event"`
}

type ApMessagePageEnvelope struct {
	Page *pages.Page `json:"page"`
}

type ApMessageOnlyIdEnvelope struct {
	ID uint64 `json:"id"`
}

type ApEventMessage struct {
	EventConst string                  `json:"event"`
	Data       *ApMessageEventEnvelope `json:"data"`
}

type ApPageMessage struct {
	EventConst string                 `json:"event"`
	Data       *ApMessagePageEnvelope `json:"data"`
}

type ApIdMessage struct {
	EventConst string                   `json:"event"`
	Data       *ApMessageOnlyIdEnvelope `json:"data"`
}

func (aep *ApEventMessage) BuildWsMessage() ([]byte, error) {
	return json.Marshal(aep)
}

func (aip *ApIdMessage) BuildWsMessage() ([]byte, error) {
	return json.Marshal(aip)
}

func (app *ApPageMessage) BuildWsMessage() ([]byte, error) {
	return json.Marshal(app)
}
