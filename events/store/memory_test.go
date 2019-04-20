package store

import (
	"github.com/nskondratev/api-page-go-back/events"
	"github.com/nskondratev/api-page-go-back/testutils"
	"reflect"
	"testing"
	"time"
)

func TestMemory_Create(t *testing.T) {
	s := NewMemory(&MemoryConfig{})

	labels, err := testutils.NewArrayNullStringFromStrings([]string{
		"Label 1",
		"Label 2",
	})

	if err != nil {
		t.Fatalf("Can not create NullString labels from strings: %s", err.Error())
	}

	cases := []testutils.MemoryCreateTestCase{
		{&events.Event{Constant: "Constant 1", Label: labels[0], Value: "Value 1"}, 1, 1},
		{&events.Event{Constant: "Constant 2", Label: labels[1], Value: "Value 2"}, 2, 2},
	}

	for caseNum, item := range cases {
		eventToCreate, ok := item.ItemToCreate.(*events.Event)

		if !ok {
			t.Errorf("[%d] Can not convert test case item to create to *events.Event type", caseNum)
		}

		err := s.Create(eventToCreate)
		if err != nil {
			t.Errorf("[%d] error while created page %+v", caseNum, eventToCreate)
		}

		if len(s.records) != item.TotalRows {
			t.Errorf("[%d] total rows mismatch. Want %d, received %d", caseNum, item.TotalRows, len(s.records))
		}

		if s.records[len(s.records)-1].ID != item.LastItemID {
			t.Errorf("[%d] last page id mismatch. Want %d, received %d", caseNum, item.LastItemID, s.records[len(s.records)-1].ID)
		}
	}
}

func TestMemory_Delete(t *testing.T) {
	s := NewMemory(&MemoryConfig{})

	labels, err := testutils.NewArrayNullStringFromStrings([]string{
		"Label 1",
		"Label 2",
	})

	if err != nil {
		t.Fatalf("Can not create NullString labels from strings: %s", err.Error())
	}

	eventsToCreate := []*events.Event{
		{Constant: "Constant 1", Label: labels[0], Value: "Value 1"},
		{Constant: "Constant 2", Label: labels[1], Value: "Value 2"},
	}

	for _, eventToCreate := range eventsToCreate {
		if err := s.Create(eventToCreate); err != nil {
			t.Fatalf("Can not create test event: %s", err.Error())
		}
	}

	cases := []testutils.MemoryDeleteTestCase{
		{&events.Event{ID: 1}, 1},
		{&events.Event{ID: 10}, 1},
	}

	for caseNum, item := range cases {
		eventToDelete, ok := item.ItemToDelete.(*events.Event)

		if !ok {
			t.Errorf("[%d] Can not convert test case item to create to *events.Event type", caseNum)
		}

		err := s.Delete(eventToDelete)
		if err != nil {
			t.Errorf("[%d] error while deleting page %+v", caseNum, eventToDelete)
		}

		if len(s.records) != item.TotalRows {
			t.Errorf("[%d] total rows mismatch. Want %d, received %d", caseNum, item.TotalRows, len(s.records))
		}

		if deletedRecord, err := s.GetById(eventToDelete.ID); deletedRecord != nil || err != nil {
			t.Errorf("[%d] Record was not deleted", caseNum)
		}
	}
}

func TestMemory_GetById(t *testing.T) {
	s := NewMemory(&MemoryConfig{})

	labels, err := testutils.NewArrayNullStringFromStrings([]string{
		"Label 1",
		"Label 2",
	})

	if err != nil {
		t.Fatalf("Can not create NullString labels from strings: %s", err.Error())
	}

	eventsToCreate := []*events.Event{
		{Constant: "Constant 1", Label: labels[0], Value: "Value 1"},
		{Constant: "Constant 2", Label: labels[1], Value: "Value 2"},
	}

	for _, eventToCreate := range eventsToCreate {
		if err := s.Create(eventToCreate); err != nil {
			t.Fatalf("Can not create test event: %s", err.Error())
		}
	}

	cases := []testutils.MemoryGetByIdTestCase{
		{1, eventsToCreate[0]},
		{2, eventsToCreate[1]},
		{3, nil},
	}

	for caseNum, item := range cases {
		event, err := s.GetById(item.ID)

		if err != nil {
			t.Errorf("[%d] failed to get event by id %+v", caseNum, item.ExpectedItem)
		}

		if item.ExpectedItem != nil {
			expectedEvent, ok := item.ExpectedItem.(*events.Event)

			if !ok {
				t.Errorf("[%d] can not cast expectedItem to *events.Event: %+v", caseNum, item.ExpectedItem)
			}

			if event != nil && !testutils.CompareEventsPart(event, expectedEvent) {
				t.Errorf("[%d] failed to fetch existent event, want %+v, received %+v", caseNum, expectedEvent, event)
			}
		} else {
			if event != nil {
				t.Errorf("[%d] failed to fetch existent event, want %+v, received %+v", caseNum, item.ExpectedItem, event)
			}
		}
	}
}

func TestMemory_Update(t *testing.T) {
	s := NewMemory(&MemoryConfig{})

	labels, err := testutils.NewArrayNullStringFromStrings([]string{
		"Label 1",
		"Label 2",
	})

	if err != nil {
		t.Fatalf("Can not create NullString labels from strings: %s", err.Error())
	}

	eventsToCreate := []*events.Event{
		{Constant: "Constant 1", Label: labels[0], Value: "Value 1"},
		{Constant: "Constant 2", Label: labels[1], Value: "Value 2"},
	}

	for _, eventToCreate := range eventsToCreate {
		if err := s.Create(eventToCreate); err != nil {
			t.Fatalf("Can not create test event: %s", err.Error())
		}
	}

	cases := []testutils.MemoryUpdateTestCase{
		{eventsToCreate[0]},
		{eventsToCreate[1]},
	}

	for caseNum, item := range cases {
		eventToUpdate, ok := item.ItemToUpdate.(*events.Event)

		if !ok {
			t.Errorf("[%d] can not cast ItemToUpdate to *events.Event: %+v", caseNum, item.ItemToUpdate)
		}

		if err := s.Update(eventToUpdate); err != nil {
			t.Errorf("[%d] error while updating page %+v", caseNum, eventToUpdate)
		}

		event, err := s.GetById(eventToUpdate.ID)

		if err != nil || !testutils.CompareEventsPart(event, eventToUpdate) {
			t.Errorf("[%d] event was not updated. Want: %+v, received: %+v", caseNum, eventToUpdate, event)
		}
	}
}

type eventToEventListTestCase struct {
	event     *events.Event
	eventList *events.EventList
}

func TestEventToEventList(t *testing.T) {
	s := NewMemory(&MemoryConfig{})

	labels, err := testutils.NewArrayNullStringFromStrings([]string{
		"Label 1",
		"Label 2",
	})

	if err != nil {
		t.Fatalf("Can not create NullString labels from strings: %s", err.Error())
	}

	eventsToCreate := []*events.Event{
		{Constant: "Constant 1", Label: labels[0], Value: "Value 1"},
		{Constant: "Constant 2", Label: labels[1], Value: "Value 2"},
	}

	for _, eventToCreate := range eventsToCreate {
		if err := s.Create(eventToCreate); err != nil {
			t.Fatalf("Can not create test event: %s", err.Error())
		}
	}

	eventLists := make([]*events.EventList, len(eventsToCreate), len(eventsToCreate))
	for i, e := range eventsToCreate {
		eventLists[i] = &events.EventList{
			ID:        e.ID,
			Constant:  e.Constant,
			Label:     e.Label,
			Value:     e.Value,
			Type:      e.Type,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		}
	}

	cases := make([]eventToEventListTestCase, len(eventsToCreate), len(eventsToCreate))

	for i, e := range eventsToCreate {
		el := eventLists[i]
		cases[i] = eventToEventListTestCase{e, el}
	}

	for caseNum, item := range cases {
		receivedEventList := EventToEventList(item.event)

		if !reflect.DeepEqual(receivedEventList, item.eventList) {
			t.Errorf("[%d] failed to convert Event to EventList. want: %+v, received: %+v", caseNum, receivedEventList, item.eventList)
		}
	}
}

type sortTestCase struct {
	sorter   by
	toSort   []*events.EventList
	expected []*events.EventList
}

var timeLayout = "2006-01-02 15:04:05"

func reverseEventsList(events []*events.EventList) []*events.EventList {
	for i := 0; i < len(events)/2; i++ {
		j := len(events) - i - 1
		events[i], events[j] = events[j], events[i]
	}
	return events
}

func TestBy_Sort(t *testing.T) {
	t1, _ := time.Parse(timeLayout, "2019-04-01 00:00:00")
	t2, _ := time.Parse(timeLayout, "2019-04-02 00:00:00")
	t3, _ := time.Parse(timeLayout, "2019-04-02 15:00:00")
	t4, _ := time.Parse(timeLayout, "2019-04-03 15:00:00")

	labels, err := testutils.NewArrayNullStringFromStrings([]string{
		"Label 1",
		"Label 2",
		"Label 3",
	})

	if err != nil {
		t.Fatalf("Can not create NullString labels from strings: %s", err.Error())
	}

	toSort := []*events.EventList{
		{ID: 1, Constant: "Constant 1", Label: labels[1], Value: "Value 3", CreatedAt: t3, UpdatedAt: t3},
		{ID: 2, Constant: "Constant 3", Label: labels[0], Value: "Value 1", CreatedAt: t2, UpdatedAt: t4},
		{ID: 3, Constant: "Constant 2", Label: labels[2], Value: "Value 2", CreatedAt: t1, UpdatedAt: t2},
	}

	sorterById, _ := getSorterByKey("id")
	sorterByConstant, _ := getSorterByKey("constant")
	sorterByLabel, _ := getSorterByKey("label")
	sorterByValue, _ := getSorterByKey("value")
	sorterByCreatedAt, _ := getSorterByKey("createdAt")
	sorterByUpdatedAt, _ := getSorterByKey("updatedAt")

	if sorter, err := getSorterByKey("unknownKey"); sorter != nil && err != nil {
		t.Errorf("fail to throw error on unkown sorting key")
	}

	cases := []sortTestCase{
		{sorterById, toSort, toSort},
		{sorterByConstant, toSort, []*events.EventList{toSort[0], toSort[2], toSort[1]}},
		{sorterByCreatedAt, toSort, []*events.EventList{toSort[2], toSort[1], toSort[0]}},
		{sorterByUpdatedAt, toSort, []*events.EventList{toSort[2], toSort[0], toSort[1]}},
		{sorterByLabel, toSort, []*events.EventList{toSort[1], toSort[0], toSort[2]}},
		{sorterByValue, toSort, []*events.EventList{toSort[1], toSort[2], toSort[0]}},
	}

	for caseNum, item := range cases {
		item.sorter.Sort(item.toSort, false)

		if !reflect.DeepEqual(item.toSort, item.expected) {
			t.Errorf("[%d] failed to sort. want: %+v, received: %+v", caseNum, item.expected, item.toSort)
		}

		item.sorter.Sort(item.toSort, true)

		if !reflect.DeepEqual(item.toSort, reverseEventsList(item.expected)) {
			t.Errorf("[%d] failed to sort. want: %+v, received: %+v", caseNum, item.expected, item.toSort)
		}
	}
}

type memoryListTestCase struct {
	offset         int
	limit          int
	sort           string
	descending     bool
	query          string
	eType          string
	expectedResult []*events.EventList
	expectedTotal  int
	isErrorNil     bool
}

func TestMemory_List(t *testing.T) {
	s := NewMemory(&MemoryConfig{})

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
		if err := s.Create(eventToCreate); err != nil {
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
		receivedList, receivedTotal, err := s.List(item.offset, item.limit, item.sort, item.descending, item.eType, item.query)
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
