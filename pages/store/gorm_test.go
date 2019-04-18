package store

import (
	"github.com/jinzhu/gorm"
	"github.com/nskondratev/api-page-go-back/db"
	"github.com/nskondratev/api-page-go-back/pages"
	"os"
	"strings"
	"testing"
)

type gormCreateTestCase struct {
	pageToCreate *pages.Page
	idToFetch    uint64
	ok           bool
}

func TestGorm_Create(t *testing.T) {
	d, ps := setup(t)
	createPagesTable(d)
	defer dropPagesTable(d)

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

			if !comparePagesPart(item.pageToCreate, p) {
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
	createPagesTable(d)
	defer dropPagesTable(d)

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

		if item.expectedPage != nil && !comparePagesPart(item.expectedPage, receivedPage) {
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
	createPagesTable(d)
	defer dropPagesTable(d)

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
	createPagesTable(d)
	defer dropPagesTable(d)

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

			if !comparePagesPart(p, item.expectedPage) {
				t.Errorf("[%d] Fetched pages mismatch. Wanted: %+v, received: %+v", caseNum, item.expectedPage, p)
			}
		}
	}
}

// Utility functions

func setup(t *testing.T) (*gorm.DB, pages.Store) {
	connStr := os.Getenv("DB_CONNECTION_STRING")
	if len(connStr) < 1 {
		connStr = "api_page:api_page@/api_page_test?charset=utf8&parseTime=True"
	}
	d, err := db.NewGorm(&db.MysqlDBConfig{
		ConnectionString: connStr,
	})

	if err != nil {
		t.Fatalf("Error while establishing connection")
	}

	ps := NewGorm(&GormConfig{
		DB: d,
	})

	return d, ps
}

func createPagesTable(d *gorm.DB) {
	d.AutoMigrate(&pages.Page{})
}

func dropPagesTable(d *gorm.DB) {
	d.DropTable(&pages.Page{})
}

func comparePagesPart(p1, p2 *pages.Page) bool {
	return p1.ID == p2.ID && strings.Compare(p1.Title, p2.Title) == 0 && strings.Compare(p1.Text, p2.Text) == 0
}
