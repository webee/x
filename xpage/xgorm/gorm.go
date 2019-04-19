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
	var total int64
	done := make(chan error, 1)
	go getTotal(p.query, &total, done)

	if err := p.query.Offset(xpage.CalcOffset(page, perPage)).Limit(perPage).Find(p.items).Error; err != nil {
		return nil, fmt.Errorf("pagination error: %v", err)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("pagination error: %v", err)
	}

	return xpage.NewPagination(page, perPage, total, p.items), nil
}

func getTotal(query *gorm.DB, total *int64, done chan error) {
	done <- query.Count(total).Error
}
