package model

import "embed"

type Pagination struct {
	Size   int         `json:"size" uri:"size" db:"size"`
	Page   int         `json:"page" uri:"page"`                 // 页码，从1开始
	Offset int         `json:"offset" uri:"offset" db:"offset"` // 偏移量，从0开始
	Data   interface{} `json:"data"`                            // 数据
	Total  int         `json:"total"`                           // 总记录数
	Pages  int         `json:"pages"`                           // 总页数
}

func (o *Pagination) Format() {
	if o.Page < 1 {
		o.Page = 1
	}
	if o.Size < 0 {
		o.Size = 0
	}
	if o.Total < 0 {
		o.Total = 0
	}
	if o.Size != 0 {
		o.Pages = (o.Total + o.Size - 1) / o.Size
	}
	o.Offset = (o.Page - 1) * o.Size
}

//go:embed *.sql
var SqlFS embed.FS
