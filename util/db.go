package util

import (
	"reflect"
	"strings"
)

// EnPointer 对象转指针
func EnPointer[T any](o T) *T {
	return &o
}

// DePointer 指针转对象
func DePointer[T any](o *T) T {
	return *o
}

// DeReference reflect指针转对象
func DeReference(o interface{}) interface{} {
	value := reflect.ValueOf(o)
	if value.Kind() == reflect.Ptr {
		return value.Elem().Interface()
	}
	return o
}

// ExtractDBColumn 提取结构体中所有字段的db标签为数组
func ExtractDBColumn(o interface{}) (columns []string) {
	o = DeReference(o)
	t := reflect.TypeOf(o)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagDB := field.Tag.Get("db")
		if tagDB != "" {
			columns = append(columns, field.Tag.Get("db"))
		}
	}
	return columns
}

// ExtractDBColumnStr 提取结构体中所有字段的db标签为sql字符串，以逗号分割
func ExtractDBColumnStr(o interface{}) string {
	columns := ExtractDBColumn(o)
	return strings.Join(columns, ",")
}

// ExtractDBColumnStrWithPrefix 提取结构体中所有字段的db标签为sql字符串，添加前缀并以逗号分割
func ExtractDBColumnStrWithPrefix(o interface{}, prefix string) string {
	columns := ExtractDBColumn(o)
	for i := 0; i < len(columns); i++ {
		columns[i] = prefix + columns[i]
	}
	return strings.Join(columns, ",")
}

// ExtractDBColumnStrWithAlias 提取结构体中所有字段的db标签为sql alias字符串，添加前缀并使用别名
func ExtractDBColumnStrWithAlias(o interface{}, prefix string) string {
	columns := ExtractDBColumn(o)
	for i := 0; i < len(columns); i++ {
		columns[i] = prefix + columns[i]
		columns[i] = columns[i] + " AS " + "'" + columns[i] + "'"
	}
	return strings.Join(columns, ",")
}

// ExtractDBNotZeroColumn 提取结构体中非空字段的db标签为数组
func ExtractDBNotZeroColumn(o interface{}) (columns []string) {
	o = DeReference(o)
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tagDB := field.Tag.Get("db")
		if tagDB != "" && !value.IsZero() {
			columns = append(columns, field.Tag.Get("db"))
		}
	}
	return columns
}

// ExtractDBNotZeroColumnStr 提取结构体中非空字段的db标签为sql字符串，以逗号分割
func ExtractDBNotZeroColumnStr(o interface{}) string {
	columns := ExtractDBNotZeroColumn(o)
	return strings.Join(columns, ",")
}

// ExtractDBNotZeroColumnStrWithPrefix 提取结构体中非空字段的db标签为sql字符串，添加前缀并以逗号分割
func ExtractDBNotZeroColumnStrWithPrefix(o interface{}, prefix string) string {
	columns := ExtractDBNotZeroColumn(o)
	for i := 0; i < len(columns); i++ {
		columns[i] = prefix + columns[i]
	}
	return strings.Join(columns, ",")
}

// ExtractDBNotZeroColumnSet 提取结构体中非空字段的db标签为sql set字符串，以逗号分割
func ExtractDBNotZeroColumnSet(o interface{}) string {
	columns := ExtractDBNotZeroColumn(o)
	sets := make([]string, 0)
	for _, column := range columns {
		sets = append(sets, column+"=:"+column)
	}

	return strings.Join(sets, ",")
}

// ExtractDBNotZeroColumnSetWithPrefix 提取结构体中非空字段的db标签为sql set字符串，添加前缀并以逗号分割
func ExtractDBNotZeroColumnSetWithPrefix(o interface{}, prefix string) string {
	columns := ExtractDBNotZeroColumn(o)
	sets := make([]string, 0)
	for _, column := range columns {
		sets = append(sets, prefix+column+"=:"+column)
	}

	return strings.Join(sets, ",")
}
