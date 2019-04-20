package store

import (
	"github.com/jinzhu/gorm"
	"github.com/nskondratev/api-page-go-back/pages"
	"github.com/nskondratev/api-page-go-back/testutils"
	"strings"
	"testing"
	"time"
)

type gormCreateTestCase struct {
	pageToCreate *pages.Page
	idToFetch    uint64
	ok           bool
}

func TestGorm_Create(t *testing.T) {
	d, ps := setup(t)
	testutils.CreatePagesTable(d)
	defer testutils.DropPagesTable(d)

	cases := []gormCreateTestCase{
		{&pages.Page{Title: "Page 1", Text: "Page 1 text"}, 1, true},
		{&pages.Page{Title: "Page 2", Text: "Page 2 text"}, 2, true},
		{&pages.Page{Title: strings.Repeat("a", 300), Text: "Too long title"}, 0, false},
	}

	for caseNum, item := range cases {
		err := ps.Create(item.pageToCreate)

		if item.ok && err != nil {
			t.Errorf("[%d] should create without error, but failed: %s", caseNum, err.Error())
		} else if !item.ok && err == nil {
			t.Errorf("[%d] should return error", caseNum)
		}

		if item.idToFetch > 0 {
			p := &pages.Page{}

			if err := d.First(p, item.idToFetch).Error; err != nil {
				t.Errorf("[%d] Can not fetch created page: %s", caseNum, err.Error())
			}

			if !testutils.ComparePagesPart(item.pageToCreate, p) {
				t.Errorf("[%d] Created and fetched pages mismatch. Wanted: %+v, received: %+v", caseNum, item.pageToCreate, p)
			}
		}
	}

}

type gormGetByIdTestCase struct {
	id           uint64
	ok           bool
	expectedPage *pages.Page
}

func TestGorm_GetById(t *testing.T) {
	d, ps := setup(t)
	testutils.CreatePagesTable(d)
	defer testutils.DropPagesTable(d)

	p := &pages.Page{
		Title: "Page 1",
		Text:  "Page 1 text",
	}

	if err := d.Create(p).Error; err != nil {
		t.Errorf("Can not create page for testing: %s", err.Error())
	}

	cases := []gormGetByIdTestCase{
		{1, true, p},
		{2, true, nil},
	}

	for caseNum, item := range cases {
		receivedPage, err := ps.GetById(item.id)

		if item.ok && err != nil {
			t.Errorf("[%d] should fetch without error, but failed: %s", caseNum, err.Error())
		} else if !item.ok && err == nil {
			t.Errorf("[%d] should return error", caseNum)
		}

		if item.expectedPage != nil && !testutils.ComparePagesPart(item.expectedPage, receivedPage) {
			t.Errorf("[%d] Fetched pages mismatch. Wanted: %+v, received: %+v", caseNum, item.expectedPage, receivedPage)
		}
	}
}

type gormDeleteTestCase struct {
	pageToDelete *pages.Page
	idToFetch    uint64
	ok           bool
}

func TestGorm_Delete(t *testing.T) {
	d, ps := setup(t)
	testutils.CreatePagesTable(d)
	defer testutils.DropPagesTable(d)

	p := &pages.Page{
		Title: "Page 1",
		Text:  "Page 1 text",
	}

	if err := d.Create(p).Error; err != nil {
		t.Errorf("Can not create page for testing: %s", err.Error())
	}

	cases := []gormDeleteTestCase{
		{p, 1, true},
		{&pages.Page{ID: 2}, 0, false},
	}

	for caseNum, item := range cases {
		err := ps.Delete(item.pageToDelete)

		if item.ok && err != nil {
			t.Errorf("[%d] should delete without error, but failed: %s", caseNum, err.Error())
		} else if !item.ok && err == nil {
			t.Errorf("[%d] should return error", caseNum)
		}

		if item.idToFetch > 0 {
			if err := d.First(p, item.idToFetch).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
				t.Errorf("[%d] Can not fetch created page: %s", caseNum, err.Error())
			}
		}
	}
}

type gormUpdateTestCase struct {
	pageToUpdate *pages.Page
	expectedPage *pages.Page
	ok           bool
}

func TestGorm_Update(t *testing.T) {
	d, ps := setup(t)
	testutils.CreatePagesTable(d)
	defer testutils.DropPagesTable(d)

	p := &pages.Page{
		Title: "Page 1",
		Text:  "Page 1 text",
	}

	if err := d.Create(p).Error; err != nil {
		t.Errorf("Can not create page for testing: %s", err.Error())
	}

	p.Title = "Page 1 updated"
	p.Text = "Page 1 text updated"

	cases := []gormUpdateTestCase{
		{p, &pages.Page{ID: 1, Title: p.Title, Text: p.Text}, true},
		{&pages.Page{ID: 2}, nil, false},
	}

	for caseNum, item := range cases {
		err := ps.Update(item.pageToUpdate)

		if item.ok && err != nil {
			t.Errorf("[%d] should update without error, but failed: %s", caseNum, err.Error())
		} else if !item.ok && err == nil {
			t.Errorf("[%d] should return error", caseNum)
		}

		if item.expectedPage != nil {
			p := &pages.Page{}

			if err := d.First(p, item.expectedPage.ID).Error; err != nil {
				t.Errorf("[%d] Can not fetch updated page: %s", caseNum, err.Error())
			}

			if !testutils.ComparePagesPart(p, item.expectedPage) {
				t.Errorf("[%d] Fetched pages mismatch. Wanted: %+v, received: %+v", caseNum, item.expectedPage, p)
			}
		}
	}
}

type gormListTestCase struct {
	offset         int
	limit          int
	sort           string
	descending     bool
	query          string
	expectedResult []*pages.PageList
	expectedTotal  int
	isErrorNil     bool
}

func TestGorm_List(t *testing.T) {
	d, ps := setup(t)
	testutils.CreatePagesTable(d)
	defer testutils.DropPagesTable(d)

	pagesToCreate := []*pages.Page{
		{Title: "Page 1", Text: "Page 1 text", CreatedAt: time.Now().Local().Add(time.Second * -3), UpdatedAt: time.Now().Local().Add(time.Second * -3)},
		{Title: "Page query 2", Text: "Page 2 text", CreatedAt: time.Now().Local().Add(time.Second * -2), UpdatedAt: time.Now().Local().Add(time.Second * -2)},
		{Title: "Test page 3", Text: "Page 3 text", CreatedAt: time.Now().Local().Add(time.Second * -1), UpdatedAt: time.Now().Local().Add(time.Second * -1)},
	}

	for _, p := range pagesToCreate {
		if err := ps.Create(p); err != nil {
			t.Fatalf("Failed to create test page: %s", err.Error())
		}
	}

	pl := make([]*pages.PageList, len(pagesToCreate), len(pagesToCreate))

	for i, p := range pagesToCreate {
		pl[i] = &pages.PageList{
			ID:        p.ID,
			Title:     p.Title,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}

	cases := []gormListTestCase{
		{0, 5, "", false, "", []*pages.PageList{pl[0], pl[1], pl[2]}, 3, true},
		{0, 5, "", true, "", []*pages.PageList{pl[2], pl[1], pl[0]}, 3, true},
		{0, 5, "title", false, "", []*pages.PageList{pl[0], pl[1], pl[2]}, 3, true},
		{0, 5, "title", true, "", []*pages.PageList{pl[2], pl[1], pl[0]}, 3, true},
		{0, -1, "", false, "", []*pages.PageList{pl[0], pl[1], pl[2]}, 3, true},
		{1, 1, "", false, "", []*pages.PageList{pl[1]}, 3, true},
		{1, 2, "", false, "", []*pages.PageList{pl[1], pl[2]}, 3, true},
		{0, 5, "", false, "page", []*pages.PageList{pl[0], pl[1], pl[2]}, 3, true},
		{0, 5, "", false, "query", []*pages.PageList{pl[1]}, 1, true},
		{0, 5, "unknownSortKey", false, "query", []*pages.PageList{}, 0, false},
	}

	for caseNum, item := range cases {
		receivedList, receivedTotal, err := ps.List(item.offset, item.limit, item.sort, item.descending, item.query)
		if item.isErrorNil && err != nil {
			t.Errorf("[%d] error while fetching list: %s", caseNum, err.Error())
		} else if !item.isErrorNil && err == nil {
			t.Errorf("[%d] error should be not nil", caseNum)
		}

		if item.isErrorNil {
			if receivedTotal != item.expectedTotal {
				t.Errorf("[%d] total mismatch. want: %d, received: %d", caseNum, item.expectedTotal, receivedTotal)
			}

			if !pageListsEqual(receivedList, item.expectedResult) {
				t.Errorf("[%d] list mismatch. want: %+v, received: %+v", caseNum, item.expectedResult, receivedList)
			}
		}
	}
}

// Utility functions

func setup(t *testing.T) (*gorm.DB, pages.Store) {
	d, err := testutils.NewGormTestDB()

	if err != nil {
		t.Fatalf("Error while establishing connection")
	}

	ps := NewGorm(&GormConfig{
		DB: d,
	})

	return d, ps
}

func pageListsEqual(pl1, pl2 []*pages.PageList) bool {
	if len(pl1) != len(pl2) {
		return false
	}

	for i, pl1item := range pl1 {
		pl2item := pl2[i]
		if !testutils.ComparePagesListPart(pl1item, pl2item) {
			return false
		}
	}

	return true
}
