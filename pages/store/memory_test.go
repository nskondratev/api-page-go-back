package store

import (
	"github.com/nskondratev/api-page-go-back/pages"
	"reflect"
	"strings"
	"testing"
	"time"
)

type createTestCase struct {
	page       *pages.Page
	totalRows  int
	lastPageID uint64
}

func TestMemory_Create(t *testing.T) {
	s := NewMemory(&MemoryConfig{})
	cases := []createTestCase{
		{&pages.Page{Title: "Page 1", Text: "Page 1 text"}, 1, 1},
		{&pages.Page{Title: "Page 2", Text: "Page 2 text"}, 2, 2},
	}
	for caseNum, item := range cases {
		err := s.Create(item.page)
		if err != nil {
			t.Errorf("[%d] error while created page %+v", caseNum, item.page)
		}

		if len(s.records) != item.totalRows {
			t.Errorf("[%d] total rows mismatch. Want %d, received %d", caseNum, item.totalRows, len(s.records))
		}

		if s.records[len(s.records)-1].ID != item.lastPageID {
			t.Errorf("[%d] last page id mismatch. Want %d, received %d", caseNum, item.lastPageID, s.records[len(s.records)-1].ID)
		}
	}
}

type deleteTestCase struct {
	page      *pages.Page
	totalRows int
}

func TestMemory_Delete(t *testing.T) {
	s := NewMemory(&MemoryConfig{})
	_ = s.Create(&pages.Page{Title: "Page 1", Text: "Page 1 text"})
	_ = s.Create(&pages.Page{Title: "Page 2", Text: "Page 2 text"})
	cases := []deleteTestCase{
		{&pages.Page{ID: 1}, 1},
		{&pages.Page{ID: 10}, 1},
	}
	for caseNum, item := range cases {
		err := s.Delete(item.page)
		if err != nil {
			t.Errorf("[%d] error while deleting page %+v", caseNum, item.page)
		}

		if len(s.records) != item.totalRows {
			t.Errorf("[%d] total rows mismatch. Want %d, received %d", caseNum, item.totalRows, len(s.records))
		}

		if deletedRecord, err := s.GetById(item.page.ID); deletedRecord != nil || err != nil {
			t.Errorf("[%d] Record was not deleted", caseNum)
		}
	}
}

type getByIdTestCase struct {
	id   uint64
	page *pages.Page
}

func TestMemory_GetById(t *testing.T) {
	s := NewMemory(&MemoryConfig{})
	p1 := &pages.Page{Title: "Page 1", Text: "Page 1 text"}
	p2 := &pages.Page{Title: "Page 2", Text: "Page 2 text"}
	_ = s.Create(p1)
	_ = s.Create(p2)

	cases := []getByIdTestCase{
		{1, p1},
		{2, p2},
		{3, nil},
	}

	for caseNum, item := range cases {
		page, err := s.GetById(item.id)

		if err != nil {
			t.Errorf("[%d] failed to get page by id %+v", caseNum, item.page)
		}

		if (page != nil && (strings.Compare(page.Title, item.page.Title) != 0 || strings.Compare(page.Text, item.page.Text) != 0 || item.id != page.ID)) || (page == nil && item.page != nil) {
			t.Errorf("[%d] failed to fetch existent page, want %+v, received %+v", caseNum, item.page, page)
		}
	}
}

type updateTestCase struct {
	page *pages.Page
}

func TestMemory_Update(t *testing.T) {
	s := NewMemory(&MemoryConfig{})
	p1 := &pages.Page{Title: "Page 1", Text: "Page 1 text"}
	p2 := &pages.Page{Title: "Page 2", Text: "Page 2 text"}
	_ = s.Create(p1)
	_ = s.Create(p2)
	up1 := &pages.Page{ID: 1, Title: "Page 1 updated", Text: "Page 1 updated"}
	up2 := &pages.Page{ID: 1, Title: "Page 2 updated", Text: "Page 2 updated"}

	cases := []updateTestCase{
		{up1},
		{up2},
	}

	for caseNum, item := range cases {
		if err := s.Update(item.page); err != nil {
			t.Errorf("[%d] error while updating page %+v", caseNum, item.page)
		}

		page, err := s.GetById(item.page.ID)

		if err != nil || (strings.Compare(page.Title, item.page.Title) != 0 || strings.Compare(page.Text, item.page.Text) != 0) {
			t.Errorf("[%d] page was not updated. Want: %+v, received: %+v", caseNum, item.page, page)
		}
	}
}

type pageToPageListTestCase struct {
	page     *pages.Page
	pageList *pages.PageList
}

func TestPageToPageList(t *testing.T) {
	s := NewMemory(&MemoryConfig{})
	p1 := &pages.Page{Title: "Page 1", Text: "Page 1 text"}
	p2 := &pages.Page{Title: "Page 2", Text: "Page 2 text"}
	_ = s.Create(p1)
	_ = s.Create(p2)

	p1s, _ := s.GetById(1)
	p2s, _ := s.GetById(2)

	pl1 := &pages.PageList{
		ID:        p1s.ID,
		Title:     p1s.Title,
		CreatedAt: p1s.CreatedAt,
		UpdatedAt: p1s.UpdatedAt,
	}
	pl2 := &pages.PageList{
		ID:        p2s.ID,
		Title:     p2s.Title,
		CreatedAt: p2s.CreatedAt,
		UpdatedAt: p2s.UpdatedAt,
	}

	cases := []pageToPageListTestCase{
		{p1s, pl1},
		{p2s, pl2},
	}

	for caseNum, item := range cases {
		receivedPageList := PageToPageList(item.page)

		if !reflect.DeepEqual(receivedPageList, item.pageList) {
			t.Errorf("[%d] failed to convert Page to PageList. want: %+v, received: %+v", caseNum, receivedPageList, item.pageList)
		}
	}
}

type sortTestCase struct {
	sorter   by
	toSort   []*pages.PageList
	expected []*pages.PageList
}

var timeLayout = "2006-01-02 15:04:05"

func reversePagesList(pages []*pages.PageList) []*pages.PageList {
	for i := 0; i < len(pages)/2; i++ {
		j := len(pages) - i - 1
		pages[i], pages[j] = pages[j], pages[i]
	}
	return pages
}

func TestBy_Sort(t *testing.T) {
	t1, _ := time.Parse(timeLayout, "2019-04-01 00:00:00")
	t2, _ := time.Parse(timeLayout, "2019-04-02 00:00:00")
	t3, _ := time.Parse(timeLayout, "2019-04-02 15:00:00")
	t4, _ := time.Parse(timeLayout, "2019-04-03 15:00:00")

	toSort := []*pages.PageList{
		{ID: 1, Title: "A", CreatedAt: t3, UpdatedAt: t3},
		{ID: 2, Title: "B", CreatedAt: t2, UpdatedAt: t4},
		{ID: 3, Title: "AB", CreatedAt: t1, UpdatedAt: t2},
	}

	sorterById, _ := getSorterByKey("id")
	sorterByTitle, _ := getSorterByKey("title")
	sorterByCreatedAt, _ := getSorterByKey("createdAt")
	sorterByUpdatedAt, _ := getSorterByKey("updatedAt")

	if sorter, err := getSorterByKey("unknownKey"); sorter != nil && err != nil {
		t.Errorf("fail to throw error on unkown sorting key")
	}

	cases := []sortTestCase{
		{sorterById, toSort, toSort},
		{sorterByTitle, toSort, []*pages.PageList{toSort[0], toSort[2], toSort[1]}},
		{sorterByCreatedAt, toSort, []*pages.PageList{toSort[2], toSort[1], toSort[0]}},
		{sorterByUpdatedAt, toSort, []*pages.PageList{toSort[2], toSort[0], toSort[1]}},
	}

	for caseNum, item := range cases {
		item.sorter.Sort(item.toSort, false)

		if !reflect.DeepEqual(item.toSort, item.expected) {
			t.Errorf("[%d] failed to sort. want: %+v, received: %+v", caseNum, item.expected, item.toSort)
		}

		item.sorter.Sort(item.toSort, true)

		if !reflect.DeepEqual(item.toSort, reversePagesList(item.expected)) {
			t.Errorf("[%d] failed to sort. want: %+v, received: %+v", caseNum, item.expected, item.toSort)
		}
	}
}

type listTestCase struct {
	offset         int
	limit          int
	sort           string
	descending     bool
	query          string
	expectedResult []*pages.PageList
	expectedTotal  int
	isErrorNil     bool
}

func TestMemory_List(t *testing.T) {
	s := NewMemory(&MemoryConfig{})

	p1 := &pages.Page{Title: "Page 1", Text: "Page 1 text"}
	p2 := &pages.Page{Title: "Page query 2", Text: "Page 2 text"}
	p3 := &pages.Page{Title: "Test page 3", Text: "Page 3 text"}
	_ = s.Create(p1)
	_ = s.Create(p2)
	_ = s.Create(p3)

	p1s, _ := s.GetById(1)
	p2s, _ := s.GetById(2)
	p3s, _ := s.GetById(3)

	pl1 := &pages.PageList{
		ID:        p1s.ID,
		Title:     p1s.Title,
		CreatedAt: p1s.CreatedAt,
		UpdatedAt: p1s.UpdatedAt,
	}
	pl2 := &pages.PageList{
		ID:        p2s.ID,
		Title:     p2s.Title,
		CreatedAt: p2s.CreatedAt,
		UpdatedAt: p2s.UpdatedAt,
	}

	pl3 := &pages.PageList{
		ID:        p3s.ID,
		Title:     p3s.Title,
		CreatedAt: p3s.CreatedAt,
		UpdatedAt: p3s.UpdatedAt,
	}

	cases := []listTestCase{
		{0, 5, "", false, "", []*pages.PageList{pl1, pl2, pl3}, 3, true},
		{0, 5, "", true, "", []*pages.PageList{pl3, pl2, pl1}, 3, true},
		{0, 5, "title", false, "", []*pages.PageList{pl1, pl2, pl3}, 3, true},
		{0, 5, "title", true, "", []*pages.PageList{pl3, pl2, pl1}, 3, true},
		{0, -1, "", false, "", []*pages.PageList{pl1, pl2, pl3}, 3, true},
		{1, 1, "", false, "", []*pages.PageList{pl2}, 3, true},
		{1, 2, "", false, "", []*pages.PageList{pl2, pl3}, 3, true},
		{0, 5, "", false, "page", []*pages.PageList{pl1, pl2, pl3}, 3, true},
		{0, 5, "", false, "query", []*pages.PageList{pl2}, 1, true},
		{0, 5, "unknownSortKey", false, "query", []*pages.PageList{pl2}, 1, false},
	}

	for caseNum, item := range cases {
		receivedList, receivedTotal, err := s.List(item.offset, item.limit, item.sort, item.descending, item.query)
		if item.isErrorNil && err != nil {
			t.Errorf("[%d] error while fetching list: %s", caseNum, err.Error())
		} else if !item.isErrorNil && err == nil {
			t.Errorf("[%d] error should be not nil", caseNum)
		}

		if item.isErrorNil {
			if receivedTotal != item.expectedTotal {
				t.Errorf("[%d] total mismatch. want: %d, received: %d", caseNum, item.expectedTotal, receivedTotal)
			}

			if !reflect.DeepEqual(receivedList, item.expectedResult) {
				t.Errorf("[%d] list mismatch. want: %+v, received: %+v", caseNum, item.expectedResult, receivedList)
			}
		}
	}
}
