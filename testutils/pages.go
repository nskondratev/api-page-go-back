package testutils

import (
	"github.com/jinzhu/gorm"
	"github.com/nskondratev/api-page-go-back/pages"
	"strings"
)

func CreatePagesTable(db *gorm.DB) {
	db.AutoMigrate(&pages.Page{})
}

func DropPagesTable(db *gorm.DB) {
	db.DropTable(&pages.Page{})
}

func ComparePagesPart(p1, p2 *pages.Page) bool {
	return p1.ID == p2.ID && strings.Compare(p1.Title, p2.Title) == 0 && strings.Compare(p1.Text, p2.Text) == 0
}

func ComparePagesListPart(p1, p2 *pages.PageList) bool {
	return p1.ID == p2.ID && strings.Compare(p1.Title, p2.Title) == 0
}
