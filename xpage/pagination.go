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
	Total    int         `json:"total"`
}

// Paginator 分页器
type Paginator interface {
	Paginate(page, perPage int) (*Pagination, error)
}
