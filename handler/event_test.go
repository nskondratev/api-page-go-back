package handler

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/nskondratev/api-page-go-back/events"
	"github.com/nskondratev/api-page-go-back/events/store"
	"github.com/nskondratev/api-page-go-back/router"
	"github.com/nskondratev/api-page-go-back/testutils"
	"github.com/nskondratev/api-page-go-back/util"
	"github.com/nskondratev/api-page-go-back/ws"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandler_CreateEvent(t *testing.T) {
	e, h, _ := setupEventHandlerTest()

	cases := []handlerCreateTestCase{
		{`{"constant":"Constant 1","value":"Value 1","label":"Label 1","description":"Description 1","type":"frontend"}`, http.StatusOK, `{"id":1,"constant":"Constant 1","label":"Label 1","value":"Value 1","description":"Description 1","type":"frontend","fields":[],"createdAt":`},
		{`{"constant":"Constant 1","value":"Value 1","label":"Label 1,"description":"Description 1}`, http.StatusUnprocessableEntity, emptyStr},
		{`{"value":"Value 1"}`, http.StatusUnprocessableEntity, emptyStr},
	}

	for caseNum, item := range cases {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(item.inputData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateEvent(c)

		if err != nil {
			t.Errorf("[%d] Fail to create event. Error: %s, input data: %s", caseNum, err.Error(), item.inputData)
		}

		if rec.Code != item.responseCode {
			t.Errorf("[%d] Unexpected response code. Wanted: %d, received: %d, response body: %s", caseNum, item.responseCode, rec.Code, rec.Body.String())
		}

		if len(item.responseBodyShouldContain) > 0 && !strings.Contains(rec.Body.String(), item.responseBodyShouldContain) {
			t.Errorf("[%d] Response body doesn't contain needed info. Wanted: %s, received: %s", caseNum, item.responseBodyShouldContain, rec.Body.String())
		}
	}
}

func TestHandler_GetEvent(t *testing.T) {
	e, h, es := setupEventHandlerTest()

	l, err := util.NewNullStringFromString("Label 1")

	if err != nil {
		t.Fatalf("Can not create label from string: %s", err.Error())
	}

	if err := es.Create(&events.Event{
		Constant:    "Constant 1",
		Label:       l,
		Value:       "Value 1",
		Type:        "frontend",
		Description: "Description 1",
	}); err != nil {
		t.Fatalf("Can not create test event: %s", err.Error())
	}

	cases := []handlerGetTestCase{
		{"1", http.StatusOK, `"id":1,"constant":"Constant 1","label":"Label 1","value":"Value 1","description":"Description 1","type":"frontend","fields":null,"createdAt"`},
		{"badparam", http.StatusUnprocessableEntity, emptyStr},
		{"45", http.StatusNotFound, `"error":"Not found"`},
	}

	for caseNum, item := range cases {
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(emptyStr))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/events/:id")
		c.SetParamNames("id")
		c.SetParamValues(item.id)

		err := h.GetEvent(c)

		if err != nil {
			t.Errorf("[%d] Fail to get event. Error: %s, id: %s", caseNum, err.Error(), item.id)
		}

		if rec.Code != item.responseCode {
			t.Errorf("[%d] Unexpected response code. Wanted: %d, received: %d, response body: %s", caseNum, item.responseCode, rec.Code, rec.Body.String())
		}

		if len(item.responseBodyShouldContain) > 0 && !strings.Contains(rec.Body.String(), item.responseBodyShouldContain) {
			t.Errorf("[%d] Response body doesn't contain needed info. Wanted: %s, received: %s", caseNum, item.responseBodyShouldContain, rec.Body.String())
		}
	}
}

func TestHandler_DeleteEvent(t *testing.T) {
	e, h, es := setupEventHandlerTest()

	l, err := util.NewNullStringFromString("Label 1")

	if err != nil {
		t.Fatalf("Can not create label from string: %s", err.Error())
	}

	if err := es.Create(&events.Event{
		Constant:    "Constant 1",
		Label:       l,
		Value:       "Value 1",
		Type:        "frontend",
		Description: "Description 1",
	}); err != nil {
		t.Fatalf("Can not create test event: %s", err.Error())
	}

	cases := []handlerDeleteTestCase{
		{"1", http.StatusOK, emptyStr},
		{"badparam", http.StatusUnprocessableEntity, emptyStr},
	}

	for caseNum, item := range cases {
		req := httptest.NewRequest(http.MethodDelete, "/", strings.NewReader(emptyStr))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/events/:id")
		c.SetParamNames("id")
		c.SetParamValues(item.id)

		err := h.DeleteEvent(c)

		if err != nil {
			t.Errorf("[%d] Fail to delete event. Error: %s, id: %s", caseNum, err.Error(), item.id)
		}

		if rec.Code != item.responseCode {
			t.Errorf("[%d] Unexpected response code. Wanted: %d, received: %d, response body: %s", caseNum, item.responseCode, rec.Code, rec.Body.String())
		}

		if len(item.responseBodyShouldContain) > 0 && !strings.Contains(rec.Body.String(), item.responseBodyShouldContain) {
			t.Errorf("[%d] Response body doesn't contain needed info. Wanted: %s, received: %s", caseNum, item.responseBodyShouldContain, rec.Body.String())
		}
	}
}

func TestHandler_UpdateEvent(t *testing.T) {
	e, h, es := setupEventHandlerTest()

	l, err := util.NewNullStringFromString("Label 1")

	if err != nil {
		t.Fatalf("Can not create label from string: %s", err.Error())
	}

	if err := es.Create(&events.Event{
		Constant:    "Constant 1",
		Label:       l,
		Value:       "Value 1",
		Type:        "frontend",
		Description: "Description 1",
	}); err != nil {
		t.Fatalf("Can not create test event: %s", err.Error())
	}

	cases := []handlerUpdateTestCase{
		{"1", `{"constant":"Constant 1 updated","value":"Value 1 updated","label":"Label 1","description":"Description 1","type":"frontend"}`, http.StatusOK, `"id":1,"constant":"Constant 1 updated","label":"Label 1","value":"Value 1 updated","description":"Description 1","type":"frontend","fields":[],"createdAt"`},
		{"badparam", `{"constant":"Constant 1 updated","value":"Value 1 updated","label":"Label 1","description":"Description 1","type":"frontend"}`, http.StatusUnprocessableEntity, emptyStr},
		{"1", `{"constant":"Constant 1 updated","value":"Value 1 updated","label":"Label 1"}`, http.StatusUnprocessableEntity, emptyStr},
		{"1", `{"constant":"Constant 1 updated","value":"Value 1 updated","label":"Label 1","description":"Description 1","type":"frontend}`, http.StatusUnprocessableEntity, emptyStr},
	}

	for caseNum, item := range cases {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(item.inputData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/events/:id")
		c.SetParamNames("id")
		c.SetParamValues(item.id)

		err := h.UpdateEvent(c)

		if err != nil {
			t.Errorf("[%d] Fail to update page. Error: %s, id: %s", caseNum, err.Error(), item.id)
		}

		if rec.Code != item.responseCode {
			t.Errorf("[%d] Unexpected response code. Wanted: %d, received: %d, response body: %s", caseNum, item.responseCode, rec.Code, rec.Body.String())
		}

		if len(item.responseBodyShouldContain) > 0 && !strings.Contains(rec.Body.String(), item.responseBodyShouldContain) {
			t.Errorf("[%d] Response body doesn't contain needed info. Wanted: %s, received: %s", caseNum, item.responseBodyShouldContain, rec.Body.String())
		}
	}
}

type listEventsResponse struct {
	Total int             `json:"total"`
	Data  []*events.Event `json:"data"`
}

type handlerListEventsTestCase struct {
	queryParams             map[string]string
	responseCode            int
	responseBodyShouldEqual *listEventsResponse
}

func TestHandler_ListEvents(t *testing.T) {
	e, h, es := setupEventHandlerTest()

	l, err := testutils.NewArrayNullStringFromStrings([]string{"Label 1", "Label 2", "Label 3", "Label 4"})

	if err != nil {
		t.Fatalf("Can not create label from string: %s", err.Error())
	}

	eventsToCreate := []*events.Event{
		{Constant: "Constant 1 query", Label: l[0], Value: "Value 3", Type: "frontend"},
		{Constant: "Constant 3", Label: l[2], Value: "Value 4", Type: "frontend"},
		{Constant: "Constant 2", Label: l[3], Value: "Value 1", Type: "client"},
		{Constant: "Constant 4", Label: l[1], Value: "Value 2", Type: "client"},
	}

	for _, e := range eventsToCreate {
		if err := es.Create(e); err != nil {
			t.Fatalf("Can not create test event: %s", err.Error())
		}
	}

	cases := []handlerListEventsTestCase{
		{emptyQueryParamsMap, http.StatusOK, &listEventsResponse{
			Total: 4,
			Data:  []*events.Event{eventsToCreate[3], eventsToCreate[2], eventsToCreate[1], eventsToCreate[0]},
		}},
		{map[string]string{
			"limit":  "1",
			"offset": "0",
			"type":   "frontend",
		}, http.StatusOK, &listEventsResponse{
			Total: 2,
			Data:  []*events.Event{eventsToCreate[1]},
		}},
		{map[string]string{
			"limit":  "1",
			"offset": "0",
			"type":   "frontend",
			"query":  "Query",
		}, http.StatusOK, &listEventsResponse{
			Total: 1,
			Data:  []*events.Event{eventsToCreate[0]},
		}},
		{map[string]string{
			"limit":      "-1",
			"offset":     "0",
			"sort":       "id",
			"descending": "true",
			"type":       "client",
		}, http.StatusOK, &listEventsResponse{
			Total: 2,
			Data:  []*events.Event{eventsToCreate[3], eventsToCreate[2]},
		}},
	}

	for caseNum, item := range cases {
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(emptyStr))

		qp := &url.Values{}

		for key, val := range item.queryParams {
			qp.Add(key, val)
		}

		encodedQueryParams := qp.Encode()

		if len(encodedQueryParams) > 0 {
			req.URL.RawQuery = encodedQueryParams
		}

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/events")

		err := h.ListEvents(c)

		if err != nil {
			t.Errorf("[%d] Fail to list events. Error: %s", caseNum, err.Error())
		}

		if rec.Code != item.responseCode {
			t.Errorf("[%d] Unexpected response code. Wanted: %d, received: %d, response body: %s", caseNum, item.responseCode, rec.Code, rec.Body.String())
		}

		if item.responseBodyShouldEqual != nil {
			parsedResBody := &listEventsResponse{}
			if err := json.Unmarshal(rec.Body.Bytes(), parsedResBody); err != nil {
				t.Errorf("[%d] Can not parse response body to json. Received: %s", caseNum, rec.Body.String())
			}

			if parsedResBody.Total != item.responseBodyShouldEqual.Total {
				t.Errorf("[%d] response total mismatch. Want: %d, received: %d", caseNum, item.responseBodyShouldEqual.Total, parsedResBody.Total)
			}

			if len(parsedResBody.Data) != len(item.responseBodyShouldEqual.Data) {
				t.Errorf("[%d] returned records mismatch. Want: %d, received: %d", caseNum, len(item.responseBodyShouldEqual.Data), len(parsedResBody.Data))
			}

			for i := 0; i < len(item.responseBodyShouldEqual.Data); i++ {
				wantedEvent := item.responseBodyShouldEqual.Data[i]
				receivedEvent := parsedResBody.Data[i]
				if wantedEvent.ID != receivedEvent.ID {
					t.Errorf("[%d] event mismatch. Wanted: %+v, received: %+v", caseNum, wantedEvent, receivedEvent)
				}
			}
		}
	}
}

// Utility functions

func setupEventHandlerTest() (*echo.Echo, *Handler, *store.Memory) {
	e := router.New()

	es := store.NewMemory(&store.MemoryConfig{
		Logger: e.Logger,
	})

	h := New(&Config{
		Logger:     e.Logger,
		EventStore: es,
		WsHub:      ws.NewHubMock(),
	})

	return e, h, es
}
