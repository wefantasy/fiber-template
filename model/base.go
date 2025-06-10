package model

type Pagination struct {
	Size   int         `json:"size" uri:"size" db:"size"`
	Page   int         `json:"page" uri:"page" db:"page"`
	Offset int         `json:"offset" uri:"offset" db:"offset"`
	Data   interface{} `json:"data"`  // 数据
	Total  int         `json:"total"` // 总记录数
	Pages  int         `json:"pages"` // 总页数
}
