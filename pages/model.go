package pages

import (
	"time"
)

//go:generate reform

//reform:page
type Page struct {
	ID        uint64    `json:"id" gorm:"AUTO_INCREMENT;primary_key" reform:"id,pk"`
	Title     string    `json:"title" gorm:"size:255;column:title" reform:"title"`
	Text      string    `json:"text" gorm:"type:text;column:text" reform:"text"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt" reform:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt" reform:"updatedAt"`
}

type PageList struct {
	ID        uint64    `json:"id" form:"id" gorm:"column:id;AUTO_INCREMENT;primary_key"`
	Title     string    `json:"title" gorm:"size:255;column:title"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt"`
}

func (Page) TableName() string {
	return "page"
}

func (PageList) TableName() string {
	return "page"
}
