package xgorm

import (
	"fmt"

	"github.com/webee/x/xpage"
	"gorm.io/gorm"
)

type gormPaginator struct {
	tx    *gorm.DB
	query *gorm.DB
	items interface{}
}

// NewPaginator 创建分页器
func NewPaginator(tx, query *gorm.DB, items interface{}) xpage.Paginator {
	return &gormPaginator{
		tx:    tx,
		query: query,
		items: items,
	}
}

// Paginate 实现Paginator接口
func (p *gormPaginator) Paginate(page, perPage int) (*xpage.Pagination, error) {

	if err := p.query.Clauses(Hints{Clauses: []string{"SELECT"}, Content: "SQL_CALC_FOUND_ROWS"}).Offset(xpage.CalcOffset(page, perPage)).Limit(perPage).Find(p.items).Error; err != nil {
		return nil, fmt.Errorf("pagination error: %v", err)
	}

	var total int64
	if err := p.tx.Raw("SELECT FOUND_ROWS()").Scan(&total).Error; err != nil {
		return nil, fmt.Errorf("pagination count error: %v", err)
	}

	return xpage.NewPagination(page, perPage, total, p.items), nil
}
