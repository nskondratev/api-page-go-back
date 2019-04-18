package store

import (
	"fmt"
	"github.com/nskondratev/api-page-go-back/events"
	"github.com/nskondratev/api-page-go-back/logger"
	"sort"
	"strings"
	"sync"
	"time"
)

type Memory struct {
	logger  logger.Logger
	records []*events.Event
	mu      *sync.Mutex
}

type MemoryConfig struct {
	Logger logger.Logger
}

func NewMemory(c *MemoryConfig) *Memory {
	return &Memory{
		logger:  c.Logger,
		records: make([]*events.Event, 0),
		mu:      &sync.Mutex{},
	}
}

func (s *Memory) GetById(id uint64) (*events.Event, error) {
	s.logger.Debugf("events.memory store GetById method call. page id = %d", id)
	for _, p := range s.records {
		s.logger.Debugf("compare with event id %d", p.ID)
		if p.ID == id {
			return p, nil
		}
	}
	return nil, nil
}

func (s *Memory) List(offset, limit int, sort string, descending bool, eType, query string) ([]*events.EventList, int, error) {
	eventsList, total := make([]*events.EventList, 0), 0
	q := strings.ToLower(query)
	for _, el := range s.records {
		if (len(q) < 1 || (strings.Contains(strings.ToLower(el.Constant), q) || strings.Contains(strings.ToLower(el.Label.String), q) || strings.Contains(strings.ToLower(el.Value), q))) && (len(eType) < 1 || eType == el.Type) {
			eventsList = append(eventsList, EventToEventList(el))
		}
	}
	total = len(eventsList)
	bSort := strings.Builder{}
	if len(sort) > 0 {
		bSort.WriteString(sort)
	} else {
		bSort.WriteString("id")
	}
	sorter, err := getSorterByKey(bSort.String())
	if err != nil {
		return eventsList, total, err
	}
	sorter.Sort(eventsList, descending)
	l := limit
	if l > total || l < 0 {
		l = total
	}
	return eventsList[offset : offset+l], total, nil
}

func (s *Memory) Update(e *events.Event) error {
	s.logger.Debugf("events.memory store update method call. event id = %d", e.ID)
	e.UpdatedAt = time.Now()
	s.mu.Lock()
	for i, el := range s.records {
		if el.ID == e.ID {
			s.records[i] = e
			break
		}
	}
	s.mu.Unlock()
	return nil
}

func (s *Memory) Delete(e *events.Event) error {
	s.logger.Debugf("events.memory store delete method call. event id = %d", e.ID)
	s.mu.Lock()
	for i, el := range s.records {
		if el.ID == e.ID {
			copy(s.records[i:], s.records[i+1:])
			s.records[len(s.records)-1] = nil
			s.records = s.records[:len(s.records)-1]
			break
		}
	}
	s.mu.Unlock()
	return nil
}

func (s *Memory) Create(e *events.Event) error {
	s.mu.Lock()
	e.ID = uint64(len(s.records) + 1)
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	s.records = append(s.records, e)
	s.mu.Unlock()
	return nil
}

// Sorting helpers

type by func(e1, e2 *events.EventList) bool

func (by by) Sort(toSort []*events.EventList, desc bool) {
	s := &eventSorter{
		events: toSort,
		by:     by,
		desc:   desc,
	}
	sort.Sort(s)
}

type eventSorter struct {
	events []*events.EventList
	by     func(e1, e2 *events.EventList) bool
	desc   bool
}

func (s *eventSorter) Len() int {
	return len(s.events)
}

func (s *eventSorter) Swap(i, j int) {
	s.events[i], s.events[j] = s.events[j], s.events[i]
}

func (s *eventSorter) Less(i, j int) bool {
	if s.desc {
		return s.by(s.events[j], s.events[i])
	} else {
		return s.by(s.events[i], s.events[j])
	}
}

const (
	id        = "id"
	constant  = "constant"
	label     = "label"
	value     = "value"
	createdAt = "createdAt"
	updatedAt = "updatedAt"
)

func getSorterByKey(key string) (by, error) {
	switch key {
	case id:
		return by(func(e1, e2 *events.EventList) bool {
			return e1.ID < e2.ID
		}), nil
	case constant:
		return by(func(e1, e2 *events.EventList) bool {
			return strings.Compare(e1.Constant, e2.Constant) == -1
		}), nil
	case label:
		return by(func(e1, e2 *events.EventList) bool {
			return strings.Compare(e1.Label.String, e2.Label.String) == -1
		}), nil
	case value:
		return by(func(e1, e2 *events.EventList) bool {
			return strings.Compare(e1.Value, e2.Value) == -1
		}), nil
	case createdAt:
		return by(func(e1, e2 *events.EventList) bool {
			return e1.CreatedAt.Before(e2.CreatedAt)
		}), nil
	case updatedAt:
		return by(func(e1, e2 *events.EventList) bool {
			return e1.UpdatedAt.Before(e2.UpdatedAt)
		}), nil
	default:
		return nil, fmt.Errorf("unkown sort key: %s", key)
	}
}

func EventToEventList(e *events.Event) *events.EventList {
	return &events.EventList{
		ID:        e.ID,
		Constant:  e.Constant,
		Label:     e.Label,
		Value:     e.Value,
		Type:      e.Type,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
