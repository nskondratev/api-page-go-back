package ws

import "net/http"

type HubMock struct{}

func NewHubMock() IHub {
	return &HubMock{}
}

func (h *HubMock) Run() {}

func (h *HubMock) Broadcast(message ApMessage) error { return nil }

func (h *HubMock) ServeWs(w http.ResponseWriter, r *http.Request) {}
