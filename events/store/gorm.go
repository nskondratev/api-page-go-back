package store

import (
	"github.com/jinzhu/gorm"
	"github.com/nskondratev/api-page-go-back/events"
	"github.com/nskondratev/api-page-go-back/logger"
	"strings"
)

type Gorm struct {
	db     *gorm.DB
	logger logger.Logger
}

type GormConfig struct {
	DB     *gorm.DB
	Logger logger.Logger
}

func NewGorm(c *GormConfig) *Gorm {
	return &Gorm{
		db:     c.DB,
		logger: c.Logger,
	}
}

func (s *Gorm) GetById(id uint64) (*events.Event, error) {
	var event events.Event
	if err := s.db.Preload("Fields").First(&event, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (s *Gorm) List(offset, limit int, sort string, descending bool, eType string, query string) ([]*events.EventList, int, error) {
	eventsList, total := []*events.EventList{nil}, 0
	bSort := strings.Builder{}
	if len(sort) > 0 {
		bSort.WriteString(sort)
	} else {
		bSort.WriteString("id")
	}
	orderDirection := " asc"
	if descending == true {
		orderDirection = " desc"
	}
	bSort.WriteString(orderDirection)
	qb := s.db.Model(&eventsList)
	if len(query) > 0 {
		qArg := "%" + query + "%"
		qb = qb.Where("`constant` LIKE ? OR `label` LIKE ? OR `value` LIKE ?", qArg, qArg, qArg)
	}
	if len(eType) > 0 {
		qb = qb.Where("`type` = ?", eType)
	}
	if err := qb.Count(&total).Error; err != nil {
		return eventsList, total, err
	}
	err := qb.Offset(offset).Limit(limit).Order(bSort.String()).Find(&eventsList).Error
	return eventsList, total, err
}

func (s *Gorm) Create(e *events.Event) error {
	return s.db.Create(e).Error
}

func (s *Gorm) Update(e *events.Event) error {
	if err := s.db.Delete(&events.Field{}, "eventId = ?", e.ID).Error; err != nil {
		return err
	}
	return s.db.Save(&e).Error
}

func (s *Gorm) Delete(e *events.Event) error {
	return s.db.Delete(e).Error
}
