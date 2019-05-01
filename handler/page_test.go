package handler

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/nskondratev/api-page-go-back/pages"
	"github.com/nskondratev/api-page-go-back/pages/store"
	"github.com/nskondratev/api-page-go-back/router"
	"github.com/nskondratev/api-page-go-back/ws"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandler_CreatePage(t *testing.T) {
	// Setup
	e, h, _ := setupPageHandlerTest()

	cases := []handlerCreateTestCase{
		{`{"title":"page1","text":"Page 1 text"}`, http.StatusOK, `"id":1,"title":"page1","text":"Page 1 text"`},
		{`{"title":"page1","text":"Page 1 text}`, http.StatusUnprocessableEntity, emptyStr},
		{`{"text":"Page 1 text"}`, http.StatusUnprocessableEntity, emptyStr},
	}

	for caseNum, item := range cases {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(item.inputData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreatePage(c)

		if err != nil {
			t.Errorf("[%d] Fail to create page. Error: %s, input data: %s", caseNum, err.Error(), item.inputData)
		}

		if rec.Code != item.responseCode {
			t.Errorf("[%d] Unexpected response code. Wanted: %d, received: %d, response body: %s", caseNum, item.responseCode, rec.Code, rec.Body.String())
		}

		if len(item.responseBodyShouldContain) > 0 && !strings.Contains(rec.Body.String(), item.responseBodyShouldContain) {
			t.Errorf("[%d] Response body doesn't contain needed info. Wanted: %s, received: %s", caseNum, item.responseBodyShouldContain, rec.Body.String())
		}
	}
}

func TestHandler_GetPage(t *testing.T) {
	e, h, ps := setupPageHandlerTest()

	_ = ps.Create(&pages.Page{
		Title: "Page 1",
		Text:  "Page 1 text",
	})

	cases := []handlerGetTestCase{
		{"1", http.StatusOK, `"id":1,"title":"Page 1","text":"Page 1 text"`},
		{"badparam", http.StatusUnprocessableEntity, emptyStr},
		{"45", http.StatusNotFound, `"error":"Not found"`},
	}

	for caseNum, item := range cases {
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(emptyStr))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/pages/:id")
		c.SetParamNames("id")
		c.SetParamValues(item.id)

		err := h.GetPage(c)

		if err != nil {
			t.Errorf("[%d] Fail to get page. Error: %s, id: %s", caseNum, err.Error(), item.id)
		}

		if rec.Code != item.responseCode {
			t.Errorf("[%d] Unexpected response code. Wanted: %d, received: %d, response body: %s", caseNum, item.responseCode, rec.Code, rec.Body.String())
		}

		if len(item.responseBodyShouldContain) > 0 && !strings.Contains(rec.Body.String(), item.responseBodyShouldContain) {
			t.Errorf("[%d] Response body doesn't contain needed info. Wanted: %s, received: %s", caseNum, item.responseBodyShouldContain, rec.Body.String())
		}
	}
}

func TestHandler_DeletePage(t *testing.T) {
	e, h, ps := setupPageHandlerTest()

	_ = ps.Create(&pages.Page{
		Title: "Page 1",
		Text:  "Page 1 text",
	})

	cases := []handlerDeleteTestCase{
		{"1", http.StatusOK, emptyStr},
		{"badparam", http.StatusUnprocessableEntity, emptyStr},
	}

	for caseNum, item := range cases {
		req := httptest.NewRequest(http.MethodDelete, "/", strings.NewReader(emptyStr))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/pages/:id")
		c.SetParamNames("id")
		c.SetParamValues(item.id)

		err := h.DeletePage(c)

		if err != nil {
			t.Errorf("[%d] Fail to delete page. Error: %s, id: %s", caseNum, err.Error(), item.id)
		}

		if rec.Code != item.responseCode {
			t.Errorf("[%d] Unexpected response code. Wanted: %d, received: %d, response body: %s", caseNum, item.responseCode, rec.Code, rec.Body.String())
		}

		if len(item.responseBodyShouldContain) > 0 && !strings.Contains(rec.Body.String(), item.responseBodyShouldContain) {
			t.Errorf("[%d] Response body doesn't contain needed info. Wanted: %s, received: %s", caseNum, item.responseBodyShouldContain, rec.Body.String())
		}
	}
}

func TestHandler_UpdatePage(t *testing.T) {
	e, h, ps := setupPageHandlerTest()

	_ = ps.Create(&pages.Page{
		Title: "Page 1",
		Text:  "Page 1 text",
	})

	cases := []handlerUpdateTestCase{
		{"1", `{"title":"Page 1 updated","text":"Page 1 updated text"}`, http.StatusOK, `"title":"Page 1 updated","text":"Page 1 updated text"`},
		{"badparam", `{"title":"Page 1 updated","text":"Page 1 updated text"}`, http.StatusUnprocessableEntity, emptyStr},
		{"1", `{"text":"Page 1 updated text"}`, http.StatusUnprocessableEntity, emptyStr},
		{"1", `{"title":"Page 1 updated","text":"Page 1 updated text}`, http.StatusUnprocessableEntity, emptyStr},
	}

	for caseNum, item := range cases {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(item.inputData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/pages/:id")
		c.SetParamNames("id")
		c.SetParamValues(item.id)

		err := h.UpdatePage(c)

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

type listPagesResponse struct {
	Total int           `json:"total"`
	Data  []*pages.Page `json:"data"`
}

type handlerListPagesTestCase struct {
	queryParams             map[string]string
	responseCode            int
	responseBodyShouldEqual *listPagesResponse
}

func TestHandler_ListPages(t *testing.T) {
	e, h, ps := setupPageHandlerTest()

	pagesToCreate := []*pages.Page{
		{Title: "Page 1", Text: "Page 1 text"},
		{Title: "Test page 2", Text: "Page 2 text"},
		{Title: "Page 3 query", Text: "Page 3 text"},
	}

	for _, page := range pagesToCreate {
		_ = ps.Create(page)
	}

	cases := []handlerListPagesTestCase{
		{emptyQueryParamsMap, http.StatusOK, &listPagesResponse{
			Total: 3,
			Data:  []*pages.Page{pagesToCreate[2], pagesToCreate[1], pagesToCreate[0]},
		}},
		{map[string]string{
			"limit":  "1",
			"offset": "0",
		}, http.StatusOK, &listPagesResponse{
			Total: 3,
			Data:  []*pages.Page{pagesToCreate[2]},
		}},
		{map[string]string{
			"query": "Query",
		}, http.StatusOK, &listPagesResponse{
			Total: 1,
			Data:  []*pages.Page{pagesToCreate[2]},
		}},
		{map[string]string{
			"limit":      "-1",
			"offset":     "0",
			"sort":       "id",
			"descending": "true",
		}, http.StatusOK, &listPagesResponse{
			Total: 3,
			Data:  []*pages.Page{pagesToCreate[2], pagesToCreate[1], pagesToCreate[0]},
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

		c.SetPath("/pages")

		err := h.ListPages(c)

		if err != nil {
			t.Errorf("[%d] Fail to list pages. Error: %s", caseNum, err.Error())
		}

		if rec.Code != item.responseCode {
			t.Errorf("[%d] Unexpected response code. Wanted: %d, received: %d, response body: %s", caseNum, item.responseCode, rec.Code, rec.Body.String())
		}

		if item.responseBodyShouldEqual != nil {
			parsedResBody := &listPagesResponse{}
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
				wantedPage := item.responseBodyShouldEqual.Data[i]
				receivedPage := parsedResBody.Data[i]
				if wantedPage.ID != receivedPage.ID {
					t.Errorf("[%d] page mismatch. Wanted: %+v, received: %+v", caseNum, wantedPage, receivedPage)
				}
			}
		}
	}
}

func setupPageHandlerTest() (*echo.Echo, *Handler, *store.Memory) {
	e := router.New()

	ps := store.NewMemory(&store.MemoryConfig{
		Logger: e.Logger,
	})

	h := New(&Config{
		Logger:    e.Logger,
		PageStore: ps,
		WsHub:     ws.NewHubMock(),
	})

	return e, h, ps
}
