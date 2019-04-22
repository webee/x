package xes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic/v7"
	"github.com/webee/x/xpage"
)

type esPaginator struct {
	search *elastic.SearchService
}

// NewPaginator 创建分页器
func NewPaginator(search *elastic.SearchService) xpage.Paginator {
	return &esPaginator{search: search}
}

// Paginate 实现Paginator接口
func (p *esPaginator) Paginate(page, perPage int) (*xpage.Pagination, error) {
	var (
		err error
		ctx = context.Background()
	)

	res, err := p.search.From(xpage.CalcOffset(page, perPage)).Size(perPage).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("pagination error: %v", err)
	}

	total := res.TotalHits()
	items := make([]map[string]interface{}, 0)
	for _, hit := range res.Hits.Hits {
		d := make(map[string]interface{})
		err := json.Unmarshal(hit.Source, &d)
		if err != nil {
			// FIXME: log here.
			continue
		}
		d["_highlight"] = hit.Highlight
		items = append(items, d)
	}

	return xpage.NewPagination(page, perPage, total, items), nil
}
