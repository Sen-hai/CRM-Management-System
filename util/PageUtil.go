package util

import "math"

// Pagination 分页结构体
type Pagination struct {
	TotalCount  int64         `json:"total_count"`  //总记录数
	TotalPages  int64         `json:"total_pages"`  //总页数
	CurrentPage int64         `json:"current_page"` // 当前页数
	PageSize    int64         `json:"page_size"`    //每页记录数
	Data        []interface{} `json:"data"`         // 当前页的数据
}

// NewPagination 创建分页实例
func NewPagination(totalCount, currentPage, pageSize int64) *Pagination {
	totalPages := int64(math.Ceil(float64(totalCount) / float64(pageSize)))
	if currentPage < 1 {
		currentPage = 1
	} else if currentPage > totalPages {
		currentPage = totalPages
	}
	return &Pagination{
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    pageSize,
	}
}
