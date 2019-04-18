package store

import (
	"fmt"
	"github.com/nskondratev/api-page-go-back/logger"
	"github.com/nskondratev/api-page-go-back/pages"
	"sort"
	"strings"
	"sync"
	"time"
)

type Memory struct {
	logger  logger.Logger
	records []*pages.Page
	mu      *sync.Mutex
}

type MemoryConfig struct {
	Logger logger.Logger
}

func NewMemory(c *MemoryConfig) *Memory {
	return &Memory{
		logger:  c.Logger,
		records: make([]*pages.Page, 0),
		mu:      &sync.Mutex{},
	}
}

func (s *Memory) GetById(id uint64) (*pages.Page, error) {
	for _, p := range s.records {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, nil
}

func (s *Memory) List(offset, limit int, sort string, descending bool, query string) ([]*pages.PageList, int, error) {
	pagesList, total := make([]*pages.PageList, 0), 0
	q := strings.ToLower(query)
	for _, el := range s.records {
		if len(q) < 1 || strings.Contains(strings.ToLower(el.Title), q) {
			pagesList = append(pagesList, PageToPageList(el))
		}
	}
	total = len(pagesList)
	bSort := strings.Builder{}
	if len(sort) > 0 {
		bSort.WriteString(sort)
	} else {
		bSort.WriteString("id")
	}
	sorter, err := getSorterByKey(bSort.String())
	if err != nil {
		return pagesList, total, err
	}
	sorter.Sort(pagesList, descending)
	l := limit
	if l > total || l < 0 {
		l = total
	}
	return pagesList[offset : offset+l], total, nil
}

func (s *Memory) Update(p *pages.Page) error {
	p.UpdatedAt = time.Now()
	s.mu.Lock()
	for i, el := range s.records {
		if el.ID == p.ID {
			s.records[i] = p
			break
		}
	}
	s.mu.Unlock()
	return nil
}

func (s *Memory) Delete(p *pages.Page) error {
	s.mu.Lock()
	for i, el := range s.records {
		if el.ID == p.ID {
			copy(s.records[i:], s.records[i+1:])
			s.records[len(s.records)-1] = nil
			s.records = s.records[:len(s.records)-1]
			break
		}
	}
	s.mu.Unlock()
	return nil
}

func (s *Memory) Create(p *pages.Page) error {
	s.mu.Lock()
	p.ID = uint64(len(s.records) + 1)
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	s.records = append(s.records, p)
	s.mu.Unlock()
	return nil
}

// Sorting helpers

type by func(p1, p2 *pages.PageList) bool

func (by by) Sort(toSort []*pages.PageList, desc bool) {
	ps := &pageSorter{
		pages: toSort,
		by:    by,
		desc:  desc,
	}
	sort.Sort(ps)
}

type pageSorter struct {
	pages []*pages.PageList
	by    func(p1, p2 *pages.PageList) bool
	desc  bool
}

func (s *pageSorter) Len() int {
	return len(s.pages)
}

func (s *pageSorter) Swap(i, j int) {
	s.pages[i], s.pages[j] = s.pages[j], s.pages[i]
}

func (s *pageSorter) Less(i, j int) bool {
	if s.desc {
		return s.by(s.pages[j], s.pages[i])
	} else {
		return s.by(s.pages[i], s.pages[j])
	}
}

const (
	id        = "id"
	title     = "title"
	createdAt = "createdAt"
	updatedAt = "updatedAt"
)

func getSorterByKey(key string) (by, error) {
	switch key {
	case id:
		return by(func(p1, p2 *pages.PageList) bool {
			return p1.ID < p2.ID
		}), nil
	case title:
		return by(func(p1, p2 *pages.PageList) bool {
			return strings.Compare(p1.Title, p2.Title) == -1
		}), nil
	case createdAt:
		return by(func(p1, p2 *pages.PageList) bool {
			return p1.CreatedAt.Before(p2.CreatedAt)
		}), nil
	case updatedAt:
		return by(func(p1, p2 *pages.PageList) bool {
			return p1.UpdatedAt.Before(p2.UpdatedAt)
		}), nil
	default:
		return nil, fmt.Errorf("unkown sort key: %s", key)
	}
}

func PageToPageList(p *pages.Page) *pages.PageList {
	return &pages.PageList{
		ID:        p.ID,
		Title:     p.Title,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
