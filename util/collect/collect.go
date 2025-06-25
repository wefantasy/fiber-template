package collect

import "reflect"

// Filter 根据给定的过滤函数过滤切片中的元素
// 参数:
//
//	slice []T - 要过滤的切片
//	f func(T) bool - 过滤函数，返回 true 表示保留元素，返回 false 表示过滤掉元素
//
// 返回值:
//
//	[]T - 过滤后的切片
func Filter[T any](slice []T, f func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func Contains[T any](slice []T, f func(T) bool) bool {
	var result []T
	result = Filter(slice, f)
	return len(result) > 0
}

// FilterByStruct 根据给定的结构体过滤切片中的元素
// 参数:
//
//	slice []T - 要过滤的切片
//	filter F - 过滤条件的结构体，结构体中的非零值字段将作为过滤条件
//
// 返回值:
//
//	[]T - 过滤后的切片
func FilterByStruct[T, F any](slice []T, filter F) []T {
	filtered := slice
	tVal := reflect.ValueOf(filter)
	tType := reflect.TypeOf(filter)

	// 遍历 t 的所有字段
	for i := 0; i < tVal.NumField(); i++ {
		fieldVal := tVal.Field(i)
		fieldType := tType.Field(i)

		// 如果字段值是零值，则跳过不作为过滤条件
		if fieldVal.IsZero() {
			continue
		}

		// 过滤当前切片，保留字段值与 t 中该字段值相等的元素
		var temp []T
		for _, elem := range filtered {
			elemVal := reflect.ValueOf(elem)
			elemFieldVal := elemVal.FieldByName(fieldType.Name)
			if !elemFieldVal.IsValid() {
				// 如果元素中没有该字段，跳过
				continue
			}
			// 比较字段值是否相等
			if reflect.DeepEqual(elemFieldVal.Interface(), fieldVal.Interface()) {
				temp = append(temp, elem)
			}
		}
		filtered = temp
	}
	return filtered
}

// Map 是一个通用的、支持泛型的函数
// T 是输入类型，U 是输出类型
// 它接收一个 T 类型的切片和一个转换函数 f(T) -> U
// 返回一个 U 类型的切片
func Map[T, U any](collection []T, iteratee func(T) U) []U {
	result := make([]U, 0, len(collection))
	for _, item := range collection {
		result = append(result, iteratee(item))
	}
	return result
}
