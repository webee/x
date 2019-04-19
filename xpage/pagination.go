package xpage

// Pagination 分页结果
type Pagination struct {
	HasPrev  bool        `json:"has_prev"`
	HasNext  bool        `json:"has_next"`
	Pages    int         `json:"pages"`
	PrevPage int         `json:"prev_page"`
	NextPage int         `json:"next_page"`
	Items    interface{} `json:"items"`
	Page     int         `json:"page"`
	PerPage  int         `json:"per_page"`
	Total    int64       `json:"total"`
}

// Paginator 分页器
type Paginator interface {
	Paginate(page, perPage int) (*Pagination, error)
}

// CalcOffset cal offset
func CalcOffset(page, perPage int) int {
	return (page - 1) * perPage
}

// NewPagination 创建分页数据
func NewPagination(page, perPage int, total int64, items interface{}) *Pagination {
	pages := int((total + int64(perPage-1)) / int64(perPage))
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
	return &Pagination{
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		Pages:    pages,
		PrevPage: prevPage,
		NextPage: nextPage,
		Items:    items,
		Page:     page,
		PerPage:  perPage,
		Total:    total,
	}
}
