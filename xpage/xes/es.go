package xes

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/webee/x/xpage"
)

// Paginate 对查询query进行分页
func Paginate(query *gorm.DB, page, perPage int, items interface{}) (*xpage.Pagination, error) {
	var total int
	done := make(chan error, 1)
	go getTotal(query, &total, done)

	if err := query.Offset((page - 1) * perPage).Limit(perPage).Find(items).Error; err != nil {
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
		Items:    items,
		Page:     page,
		PerPage:  perPage,
		Total:    total}, nil
}

func getTotal(query *gorm.DB, total *int, done chan error) {
	done <- query.Count(total).Error
}
