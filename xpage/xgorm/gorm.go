package xgorm

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/webee/x/xpage"
)

type gormPaginator struct {
	query *gorm.DB
	items interface{}
}

// NewPaginator 创建分页器
func NewPaginator(query *gorm.DB, items interface{}) xpage.Paginator {
	return &gormPaginator{
		query: query,
		items: items,
	}
}

// Paginate 实现Paginator接口
func (p *gormPaginator) Paginate(page, perPage int) (*xpage.Pagination, error) {
	var total int
	done := make(chan error, 1)
	go getTotal(p.query, &total, done)

	if err := p.query.Offset((page - 1) * perPage).Limit(perPage).Find(p.items).Error; err != nil {
		return nil, fmt.Errorf("pagination error: %v", err)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("pagination error: %v", err)
	}

	pages := (total + perPage - 1) / perPage
	hasPrev := page > 1
	hasNext := page < pages
	prevPage := page
	if hasPrev {
		prevPage = page - 1
	}
	nextPage := page
	if hasNext {
		nextPage = page + 1
	}
	return &xpage.Pagination{
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		Pages:    pages,
		PrevPage: prevPage,
		NextPage: nextPage,
		Items:    p.items,
		Page:     page,
		PerPage:  perPage,
		Total:    total}, nil
}

func getTotal(query *gorm.DB, total *int, done chan error) {
	done <- query.Count(total).Error
}
