package store

import (
	"github.com/jinzhu/gorm"
	"github.com/nskondratev/api-page-go-back/events"
	"github.com/nskondratev/api-page-go-back/testutils"
	"strings"
	"testing"
)

type gormCreateTestCase struct {
	eventToCreate *events.Event
	idToFetch     uint64
	ok            bool
}

func TestGorm_Create(t *testing.T) {
	d, es := setup(t)
	testutils.CreateEventsTable(d)
	defer testutils.DropEventsTable(d)

	labels, err := testutils.NewArrayNullStringFromStrings([]string{"Label 1", "Label 2", "Label 3"})
	if err != nil {
		t.Fatalf("Can not create labels: %s", err.Error())
	}

	cases := []gormCreateTestCase{
		{&events.Event{Constant: "Constant 1", Label: labels[0], Value: "Value 1", Type: "frontend"}, 1, true},
		{&events.Event{Constant: "Constant 2", Label: labels[1], Value: "Value 2", Type: "client"}, 2, true},
		{&events.Event{Constant: strings.Repeat("a", 300), Label: labels[1], Value: "Value 2", Type: "client"}, 0, false},
	}

	for caseNum, item := range cases {
		err := es.Create(item.eventToCreate)

		if item.ok && err != nil {
			t.Errorf("[%d] should create without error, but failed: %s", caseNum, err.Error())
		} else if !item.ok && err == nil {
			t.Errorf("[%d] should return error", caseNum)
		}

		if item.idToFetch > 0 {
			e := &events.Event{}

			if err := d.First(e, item.idToFetch).Error; err != nil {
				t.Errorf("[%d] Can not fetch created event: %s", caseNum, err.Error())
			}

			if !testutils.CompareEventsPart(item.eventToCreate, e) {
				t.Errorf("[%d] Created and fetched events mismatch. Wanted: %+v, received: %+v", caseNum, item.eventToCreate, e)
			}
		}
	}
}

type gormGetByIdTestCase struct {
	id            uint64
	ok            bool
	expectedEvent *events.Event
}

func TestGorm_GetById(t *testing.T) {
	d, es := setup(t)
	testutils.CreateEventsTable(d)
	defer testutils.DropEventsTable(d)

	labels, err := testutils.NewArrayNullStringFromStrings([]string{"Label 1"})
	if err != nil {
		t.Fatalf("Can not create labels: %s", err.Error())
	}

	e := &events.Event{
		Constant: "Constant 1",
		Label:    labels[0],
		Value:    "Value 1",
		Type:     "frontend",
	}

	if err := d.Create(e).Error; err != nil {
		t.Errorf("Can not create event for testing: %s", err.Error())
	}

	cases := []gormGetByIdTestCase{
		{1, true, e},
		{2, true, nil},
	}

	for caseNum, item := range cases {
		receivedEvent, err := es.GetById(item.id)

		if item.ok && err != nil {
			t.Errorf("[%d] should fetch without error, but failed: %s", caseNum, err.Error())
		} else if !item.ok && err == nil {
			t.Errorf("[%d] should return error", caseNum)
		}

		if item.expectedEvent != nil && !testutils.CompareEventsPart(item.expectedEvent, receivedEvent) {
			t.Errorf("[%d] Fetched pages mismatch. Wanted: %+v, received: %+v", caseNum, item.expectedEvent, receivedEvent)
		}
	}
}

type gormDeleteTestCase struct {
	eventToDelete *events.Event
	idToFetch     uint64
	ok            bool
}

func TestGorm_Delete(t *testing.T) {
	d, es := setup(t)
	testutils.CreateEventsTable(d)
	defer testutils.DropEventsTable(d)

	labels, err := testutils.NewArrayNullStringFromStrings([]string{"Label 1"})
	if err != nil {
		t.Fatalf("Can not create labels: %s", err.Error())
	}

	e := &events.Event{
		Constant: "Constant 1",
		Label:    labels[0],
		Value:    "Value 1",
		Type:     "frontend",
	}

	if err := d.Create(e).Error; err != nil {
		t.Errorf("Can not create event for testing: %s", err.Error())
	}

	cases := []gormDeleteTestCase{
		{e, 1, true},
		{&events.Event{ID: 2}, 0, false},
	}

	for caseNum, item := range cases {
		err := es.Delete(item.eventToDelete)

		if item.ok && err != nil {
			t.Errorf("[%d] should delete without error, but failed: %s", caseNum, err.Error())
		} else if !item.ok && err == nil {
			t.Errorf("[%d] should return error", caseNum)
		}

		if item.idToFetch > 0 {
			if err := d.First(e, item.idToFetch).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
				t.Errorf("[%d] Can not fetch created event: %s", caseNum, err.Error())
			}
		}
	}
}

type gormUpdateTestCase struct {
	eventToUpdate *events.Event
	expectedEvent *events.Event
	ok            bool
}

func TestGorm_Update(t *testing.T) {
	d, es := setup(t)
	testutils.CreateEventsTable(d)
	defer testutils.DropEventsTable(d)

	labels, err := testutils.NewArrayNullStringFromStrings([]string{"Label 1"})
	if err != nil {
		t.Fatalf("Can not create labels: %s", err.Error())
	}

	e := &events.Event{
		Constant: "Constant 1",
		Label:    labels[0],
		Value:    "Value 1",
		Type:     "frontend",
	}

	if err := d.Create(e).Error; err != nil {
		t.Errorf("Can not create event for testing: %s", err.Error())
	}

	e.Constant = "Constant 1 updated"
	e.Value = "Value 1 updated"

	cases := []gormUpdateTestCase{
		{e, &events.Event{ID: 1, Constant: e.Constant, Label: e.Label, Value: e.Value, Type: e.Type}, true},
		{&events.Event{ID: 2}, nil, false},
	}

	for caseNum, item := range cases {
		err := es.Update(item.eventToUpdate)

		if item.ok && err != nil {
			t.Errorf("[%d] should update without error, but failed: %s", caseNum, err.Error())
		} else if !item.ok && err == nil {
			t.Errorf("[%d] should return error", caseNum)
		}

		if item.expectedEvent != nil {
			e := &events.Event{}

			if err := d.First(e, item.expectedEvent.ID).Error; err != nil {
				t.Errorf("[%d] Can not fetch updated event: %s", caseNum, err.Error())
			}

			if !testutils.CompareEventsPart(e, item.expectedEvent) {
				t.Errorf("[%d] Fetched events mismatch. Wanted: %+v, received: %+v", caseNum, item.expectedEvent, e)
			}
		}
	}
}

type gormListTestCase struct {
	offset         int
	limit          int
	sort           string
	descending     bool
	eType          string
	query          string
	expectedResult []*events.EventList
	expectedTotal  int
	isErrorNil     bool
}

func TestGorm_List(t *testing.T) {
	d, es := setup(t)
	testutils.CreateEventsTable(d)
	defer testutils.DropEventsTable(d)

	labels, err := testutils.NewArrayNullStringFromStrings([]string{
		"Label 1",
		"Label 2",
		"Label 3",
		"Label 4",
	})

	if err != nil {
		t.Fatalf("Can not create NullString labels from strings: %s", err.Error())
	}

	eventsToCreate := []*events.Event{
		{Constant: "Constant 1 query", Label: labels[0], Value: "Value 1", Type: "frontend"},
		{Constant: "AConstant 3", Label: labels[1], Value: "Value 2", Type: "frontend"},
		{Constant: "Constant 2", Label: labels[2], Value: "Value 3", Type: "client"},
		{Constant: "Constant 4", Label: labels[2], Value: "Value 4", Type: "client"},
	}

	eventsLists := make([]*events.EventList, len(eventsToCreate), len(eventsToCreate))

	for i, eventToCreate := range eventsToCreate {
		if err := es.Create(eventToCreate); err != nil {
			t.Fatalf("Can not create test event: %s", err.Error())
		}
		eventsLists[i] = EventToEventList(eventToCreate)
	}

	cases := []memoryListTestCase{
		{0, 5, "", false, "", "frontend", []*events.EventList{eventsLists[0], eventsLists[1]}, 2, true},
		{0, 5, "", false, "", "client", []*events.EventList{eventsLists[2], eventsLists[3]}, 2, true},
		{0, 5, "", true, "", "frontend", []*events.EventList{eventsLists[1], eventsLists[0]}, 2, true},
		{0, 5, "constant", false, "", "frontend", []*events.EventList{eventsLists[1], eventsLists[0]}, 2, true},
		{0, 5, "constant", true, "", "frontend", []*events.EventList{eventsLists[0], eventsLists[1]}, 2, true},
		{0, -1, "", false, "", "frontend", []*events.EventList{eventsLists[0], eventsLists[1]}, 2, true},
		{1, 1, "", false, "", "frontend", []*events.EventList{eventsLists[1]}, 2, true},
		{0, 5, "", false, "constant", "client", []*events.EventList{eventsLists[2], eventsLists[3]}, 2, true},
		{0, 5, "", false, "query", "frontend", []*events.EventList{eventsLists[0]}, 1, true},
		{0, 5, "unknownSortKey", false, "query", "client", []*events.EventList{}, 1, false},
	}

	for caseNum, item := range cases {
		receivedList, receivedTotal, err := es.List(item.offset, item.limit, item.sort, item.descending, item.eType, item.query)
		if item.isErrorNil && err != nil {
			t.Errorf("[%d] error while fetching list: %s", caseNum, err.Error())
		} else if !item.isErrorNil && err == nil {
			t.Errorf("[%d] error should be not nil", caseNum)
		}

		if item.isErrorNil {
			if receivedTotal != item.expectedTotal {
				t.Errorf("[%d] total mismatch. want: %d, received: %d", caseNum, item.expectedTotal, receivedTotal)
			}

			if !eventListsEqual(receivedList, item.expectedResult) {
				t.Errorf("[%d] list mismatch. want: %+v, received: %+v", caseNum, item.expectedResult, receivedList)
			}
		}
	}
}

// Utility functions

func setup(t *testing.T) (*gorm.DB, events.Store) {
	d, err := testutils.NewGormTestDB()

	if err != nil {
		t.Fatalf("Error while establishing connection")
	}

	es := NewGorm(&GormConfig{
		DB: d,
	})

	return d, es
}

func eventListsEqual(el1, el2 []*events.EventList) bool {
	if len(el1) != len(el2) {
		return false
	}

	for i, el1item := range el1 {
		el2item := el2[i]
		if !testutils.CompareEventsListPart(el1item, el2item) {
			return false
		}
	}

	return true
}
