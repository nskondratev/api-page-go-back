package events

import (
	"github.com/nskondratev/api-page-go-back/util"
	"time"
)

type Field struct {
	ID          uint64          `json:"id" gorm:"AUTO_INCREMENT;primary_key"`
	EventId     uint64          `json:"eventId" gorm:"column:eventId"`
	Type        string          `json:"type" gorm:"size:255;column:type"`
	Key         util.NullString `json:"key" gorm:"size:255;column:key"`
	Required    bool            `json:"required" gorm:"type:TINYINT(1);default:0;column:required"`
	Description string          `json:"description" gorm:"type:text;column:description"`
	CreatedAt   time.Time       `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt" gorm:"column:updatedAt"`
}

func (Field) TableName() string {
	return "event_fields"
}

type Event struct {
	ID          uint64          `json:"id" gorm:"AUTO_INCREMENT;primary_key"`
	Constant    string          `json:"constant" gorm:"size:255;column:constant"`
	Label       util.NullString `json:"label" gorm:"size:255;column:label"`
	Value       string          `json:"value" gorm:"size:255;column:value"`
	Description string          `json:"description" gorm:"type:text;column:description"`
	Type        string          `json:"type" gorm:"type:ENUM('frontend','client');default:'frontend'"`
	Fields      []Field         `json:"fields" gorm:"foreignKey:eventId;"`
	CreatedAt   time.Time       `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt" gorm:"column:updatedAt"`
}

type EventList struct {
	ID        uint64          `json:"id" gorm:"AUTO_INCREMENT;primary_key"`
	Constant  string          `json:"constant" gorm:"size:255;column:constant"`
	Label     util.NullString `json:"label" gorm:"size:255;column:label"`
	Value     string          `json:"value" gorm:"size:255;column:value"`
	Type      string          `json:"type" gorm:"type:ENUM('frontend','client');default:'frontend'"`
	CreatedAt time.Time       `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt time.Time       `json:"updatedAt" gorm:"column:updatedAt"`
}

func (Event) TableName() string {
	return "events"
}

func (EventList) TableName() string {
	return "events"
}
