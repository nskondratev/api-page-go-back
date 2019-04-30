package store

import (
	"fmt"
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

func NewGorm(c *GormConfig) events.Store {
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
	res := s.db.First(&events.Event{}, e.ID)

	if res.Error != nil {
		if gorm.IsRecordNotFoundError(res.Error) {
			return fmt.Errorf("[events.store.gorm] page with id = %d does not exist", e.ID)
		}
		return res.Error
	}

	if err := s.db.Delete(&events.Field{}, "eventId = ?", e.ID).Error; err != nil {
		return err
	}

	res = s.db.Save(&e)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected < 1 {
		return fmt.Errorf("[events.store.gorm] page with id = %d was not updated", e.ID)
	}

	return nil
}

func (s *Gorm) Delete(e *events.Event) error {
	res := s.db.Delete(e)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected < 1 {
		return fmt.Errorf("[events.store.gorm] page with id = %d was not deleted", e.ID)
	}
	return nil
}
