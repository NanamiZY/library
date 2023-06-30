package tools

import (
	"strconv"
)

// 前端常用的select * from ... where id>10 order by ... limit 10 查id为10-20的数据
// 给book表建立几个索引
// 瀑布流分页是怎么实现的
// Google是怎么实现分页的
// 后端常用
type Page[T any] struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	Total       int `json:"total"` //总数
	Pages       int `json:"pages"` //总页数
	Result      []T `json:"result"`
}

func Pages[T any](res []T, currentPageString string) Page[T] {
	currentPage, _ := strconv.Atoi(currentPageString)
	pageSize := 10
	var pages int
	if len(res)%pageSize == 0 {
		pages = len(res) / pageSize
	} else {
		pages = len(res)/pageSize + 1
	}
	if !(currentPage >= 1 && currentPage <= pages) {
		return Page[T]{
			Result: nil,
		}
	}
	offset := (currentPage - 1) * pageSize
	limit := pageSize
	result := res[offset : offset+limit]
	if len(result) == 0 {
		return Page[T]{}
	}
	page := Page[T]{
		CurrentPage: currentPage,
		PageSize:    pageSize,
		Total:       len(res),
		Pages:       pages,
		Result:      result,
	}
	return page
}
