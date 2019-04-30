package store

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/nskondratev/api-page-go-back/logger"
	"github.com/nskondratev/api-page-go-back/pages"
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

func NewGorm(c *GormConfig) pages.Store {
	return &Gorm{
		db:     c.DB,
		logger: c.Logger,
	}
}

func (ps *Gorm) GetById(id uint64) (*pages.Page, error) {
	var page pages.Page
	if err := ps.db.First(&page, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &page, nil
}

func (ps *Gorm) List(offset, limit int, sort string, descending bool, query string) ([]*pages.PageList, int, error) {
	pagesList, total := []*pages.PageList{nil}, 0
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
	qb := ps.db.Model(&pagesList)
	if len(query) > 0 {
		qb = qb.Where("`title` LIKE ?", "%"+query+"%")
	}
	if err := qb.Count(&total).Error; err != nil {
		return pagesList, total, err
	}
	err := qb.Offset(offset).Limit(limit).Order(bSort.String()).Find(&pagesList).Error
	return pagesList, total, err
}

func (ps *Gorm) Update(p *pages.Page) error {
	res := ps.db.First(&pages.Page{}, p.ID)

	if res.Error != nil {
		if gorm.IsRecordNotFoundError(res.Error) {
			return fmt.Errorf("[pages.store.gorm] page with id = %d does not exist", p.ID)
		}
		return res.Error
	}

	res = ps.db.Save(&p)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected < 1 {
		return fmt.Errorf("[pages.store.gorm] page with id = %d was not updated", p.ID)
	}

	return nil
}

func (ps *Gorm) Delete(p *pages.Page) error {
	res := ps.db.Delete(p)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected < 1 {
		return fmt.Errorf("[pages.store.gorm] page with id = %d was not deleted", p.ID)
	}
	return nil
}

func (ps *Gorm) Create(p *pages.Page) error {
	return ps.db.Create(p).Error
}
