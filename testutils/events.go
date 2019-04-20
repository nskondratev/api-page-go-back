package testutils

import (
	"github.com/jinzhu/gorm"
	"github.com/nskondratev/api-page-go-back/events"
	"strings"
)

func CreateEventsTable(db *gorm.DB) {
	db.AutoMigrate(&events.Event{}).AutoMigrate(&events.Field{})
}

func DropEventsTable(db *gorm.DB) {
	db.DropTable(&events.Field{}).DropTable(&events.Event{})
}

func CompareEventsListPart(e1, e2 *events.EventList) bool {
	return e1.ID == e2.ID &&
		strings.Compare(e1.Constant, e2.Constant) == 0 &&
		strings.Compare(e1.Value, e2.Value) == 0 &&
		strings.Compare(e1.Label.String, e2.Label.String) == 0
}

func CompareEventsPart(e1, e2 *events.Event) bool {
	return e1.ID == e2.ID &&
		strings.Compare(e1.Constant, e2.Constant) == 0 &&
		strings.Compare(e1.Value, e2.Value) == 0 &&
		strings.Compare(e1.Label.String, e2.Label.String) == 0 &&
		strings.Compare(e1.Description, e2.Description) == 0 &&
		compareFields(e1.Fields, e2.Fields)
}

func CompareFieldsPart(f1, f2 *events.Field) bool {
	return f1.ID == f2.ID &&
		strings.Compare(f1.Key.String, f2.Key.String) == 0 &&
		f1.Required == f2.Required &&
		strings.Compare(f1.Type, f2.Type) == 0 &&
		strings.Compare(f1.Description, f2.Description) == 0
}

func compareFields(fa1, fa2 []events.Field) bool {
	if len(fa1) != len(fa2) {
		return false
	}

	for i, fa1item := range fa2 {
		fa2item := fa2[i]
		if !CompareFieldsPart(&fa1item, &fa2item) {
			return false
		}
	}

	return true
}
