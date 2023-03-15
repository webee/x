package xgorm

import (
	"fmt"
	"reflect"

	"github.com/webee/x/xpage"
	"gorm.io/gorm"
)

type gormPaginator struct {
	tx     *gorm.DB
	query  *gorm.DB
	items  interface{}
	option *PaginationOption
}

type PaginationOption struct {
	UseSecond   bool
	IgnoreTotal bool
}

// NewPaginator 创建分页器
func NewPaginator(tx, query *gorm.DB, items interface{}, option *PaginationOption) xpage.Paginator {
	return &gormPaginator{
		tx:     tx,
		query:  query,
		items:  items,
		option: option,
	}
}

// Paginate 实现Paginator接口
func (p *gormPaginator) Paginate(page, perPage int) (*xpage.Pagination, error) {
	ignoreTotal := p.option != nil && p.option.IgnoreTotal
	if !ignoreTotal && p.option != nil && p.option.UseSecond {
		return p.Paginate2(page, perPage)
	}

	var total int64
	if ignoreTotal {
		if err := p.query.Offset(xpage.CalcOffset(page, perPage)).Limit(perPage).Find(p.items).Error; err != nil {
			return nil, fmt.Errorf("pagination error: %v", err)
		}
	} else {
		if err := p.query.Clauses(Hints{Clauses: []string{"SELECT"}, Content: "SQL_CALC_FOUND_ROWS"}).Offset(xpage.CalcOffset(page, perPage)).Limit(perPage).Find(p.items).Error; err != nil {
			return nil, fmt.Errorf("pagination error: %v", err)
		}

		if err := p.tx.Raw("SELECT FOUND_ROWS()").Scan(&total).Error; err != nil {
			return nil, fmt.Errorf("pagination count error: %v", err)
		}
	}

	return xpage.NewPagination(page, perPage, total, p.items), nil
}

// Paginate2 实现Paginator接口
func (p *gormPaginator) Paginate2(page, perPage int) (*xpage.Pagination, error) {
	model := reflect.New(reflect.Indirect(reflect.ValueOf(p.items)).Type().Elem()).Interface()
	var total int64
	if err := p.query.Model(model).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("pagination count error: %v", err)
	}

	if err := p.query.Offset(xpage.CalcOffset(page, perPage)).Limit(perPage).Find(p.items).Error; err != nil {
		return nil, fmt.Errorf("pagination error: %v", err)
	}

	return xpage.NewPagination(page, perPage, total, p.items), nil
}
