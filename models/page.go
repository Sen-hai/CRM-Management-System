package models

type Pagination struct {
	TotalCount  int64         `json:"total_count"`  // 总记录数
	TotalPages  int           `json:"total_pages"`  // 总页数
	CurrentPage int           `json:"current_page"` // 当前页数
	PageSize    int           `json:"page_size"`    // 每页记录数
	Data        []interface{} `json:"data"`         // 当前页的数据
}
