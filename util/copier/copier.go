package copier

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// TransferListType 转化列表类型
//
//	src any - 源对象，可以是结构体或map[string]any
//	dest any - 目标对象，必须是指向结构体的指针
func TransferListType[S, T any](objects []S, targetObjects *[]T) error {
	for _, object := range objects {
		var targetObject T
		err := CopyProperties(object, &targetObject)
		if err != nil {
			return err
		}
		*targetObjects = append(*targetObjects, targetObject)
	}
	return nil
}

// CopyProperties 将源对象的属性值复制到目标对象中
//
//	src any - 源对象，可以是结构体或map[string]any或基本类型
//	dest any - 目标对象，必须是指针
func CopyProperties(src, dest any) error {
	destVal := reflect.ValueOf(dest)
	// 检查目标对象是否是指针类型
	if destVal.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer")
	}
	destElem := destVal.Elem()

	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() == reflect.Ptr {
		// 检查指针是否为nil
		if srcVal.IsNil() {
			return errors.New("src pointer is nil")
		}
		// 解引用指针获取实际值
		srcVal = srcVal.Elem()
	}

	switch srcVal.Kind() {
	case reflect.Struct:
		// 如果源对象是结构体，调用copyFromStruct函数
		return copyFromStruct(srcVal, destElem)
	case reflect.Map:
		// 如果源对象是map，检查key是否为string类型
		if srcVal.Type().Key().Kind() != reflect.String {
			return errors.New("src map keys must be string")
		}
		return copyFromMap(srcVal, destElem)
	default:
		return copyValue(srcVal, destElem)
	}
}

// copyFromStruct 从源结构体复制字段值到目标结构体
//
//	srcStructVal reflect.Value - 源结构体的反射值
//	destStructElem reflect.Value - 目标结构体元素的反射值(必须可设置)
func copyFromStruct(srcStructVal, destStructElem reflect.Value) error {
	srcType := srcStructVal.Type()
	// 遍历源结构体的所有字段
	for i := 0; i < srcStructVal.NumField(); i++ {
		srcFieldVal := srcStructVal.Field(i)
		srcFieldType := srcType.Field(i)

		// 跳过未导出的字段(包路径不为空表示未导出)
		if srcFieldType.PkgPath != "" {
			continue
		}

		// 在目标结构体中查找同名字段
		destField := destStructElem.FieldByName(srcFieldType.Name)
		// 如果目标字段不存在或不可设置，则跳过
		if !destField.IsValid() || !destField.CanSet() {
			continue
		}
		copyValue(srcFieldVal, destField)
	}
	return nil
}

// copyFromMap 从map类型源对象复制值到目标结构体
//
//	srcMapVal reflect.Value - 源map的反射值
//	destStructElem reflect.Value - 目标结构体元素的反射值(必须可设置)
func copyFromMap(srcMapVal, destStructElem reflect.Value) error {
	for _, keyVal := range srcMapVal.MapKeys() {
		// 获取map key对应的字段名
		fieldName := keyVal.String()
		// 获取map中对应key的值
		srcFieldValFromMap := srcMapVal.MapIndex(keyVal)

		// 在目标结构体中查找同名字段
		destField := destStructElem.FieldByName(fieldName)
		// 如果目标字段不存在或不可设置，则跳过
		if !destField.IsValid() || !destField.CanSet() {
			continue
		}
		copyValue(srcFieldValFromMap, destField)
	}
	return nil
}

// copyValue 执行实际的属性值复制和类型转换
// 参数:
//
//	srcValInput reflect.Value - 源值的反射值
//	destVal reflect.Value - 目标值的反射值
//
// 功能:
//  1. 处理接口类型解包
//  2. 处理指针类型解引用
//  3. 处理各种类型间的转换(数字、字符串、布尔等)
//  4. 处理目标值为指针类型的情况
func copyValue(srcValInput, destVal reflect.Value) error {
	// 创建可修改的副本
	srcVal := srcValInput

	if !srcVal.IsValid() {
		return errors.New("src is invalid")
	}

	// 处理接口类型解包
	if srcVal.Kind() == reflect.Interface {
		if srcVal.IsNil() {
			if destVal.Kind() == reflect.Ptr && destVal.CanSet() && !destVal.IsNil() {
				destVal.Set(reflect.Zero(destVal.Type()))
			}
			return errors.New("src interface is nil")
		}
		srcVal = srcVal.Elem()
	}

	// 处理指针类型解引用
	if srcVal.Kind() == reflect.Ptr {
		if srcVal.IsNil() {
			if destVal.Kind() == reflect.Ptr && destVal.CanSet() && !destVal.IsNil() {
				destVal.Set(reflect.Zero(destVal.Type()))
			}
			return errors.New("src pointer is nil")
		}
		srcVal = srcVal.Elem()
	}

	// 获取目标类型信息
	var targetType reflect.Type
	isDestPtr := destVal.Kind() == reflect.Ptr
	if isDestPtr {
		targetType = destVal.Type().Elem()
	} else {
		targetType = destVal.Type()
	}

	// 储存转换后的值
	var convertedVal reflect.Value
	// 是否成功转换
	conversionSuccessful := false

	// 处理的是需要手动实现的复杂类型转换
	{
		// 创建目标类型的临时值
		tempTargetVal := reflect.New(targetType).Elem()
		srcKind := srcVal.Kind()
		targetKind := tempTargetVal.Kind()

		// 根据目标类型确定位数
		bits := 0
		switch targetKind {
		case reflect.Int8, reflect.Uint8:
			bits = 8
		case reflect.Int16, reflect.Uint16:
			bits = 16
		case reflect.Int32, reflect.Uint32:
			bits = 32
		case reflect.Int64, reflect.Uint64:
			bits = 64
		case reflect.Int, reflect.Uint:
			bits = strconv.IntSize
		}

		// 根据目标类型进行特定转换
		switch targetKind {
		case reflect.String:
			// 从数值或布尔转字符串
			if isNumeric(srcKind) || srcKind == reflect.Bool {
				if srcVal.CanInterface() {
					tempTargetVal.SetString(fmt.Sprintf("%v", srcVal.Interface()))
					conversionSuccessful = true
				}
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// 从字符串转整型
			if srcKind == reflect.String {
				valStr := srcVal.String()
				if i, err := strconv.ParseInt(valStr, 10, bits); err == nil { // 尝试直接转整型
					tempTargetVal.SetInt(i)
					conversionSuccessful = true
				} else { // 尝试先转浮点数再转整型
					if f, errFloat := strconv.ParseFloat(valStr, 64); errFloat == nil {
						if f == math.Trunc(f) { // 检查是否为整数
							wholeNumberStr := fmt.Sprintf("%.0f", f)
							if iWhole, errWhole := strconv.ParseInt(wholeNumberStr, 10, bits); errWhole == nil {
								tempTargetVal.SetInt(iWhole)
								conversionSuccessful = true
							}
						}
					}
				}
			} else if isNumeric(srcKind) { // 从数值转整型
				f64Val := toFloat64(srcVal)
				if f64Val == math.Trunc(f64Val) {
					valToSet := int64(f64Val)
					if !tempTargetVal.OverflowInt(valToSet) {
						tempTargetVal.SetInt(valToSet)
						conversionSuccessful = true
					}
				}
			} else if srcKind == reflect.Bool { // 从布尔转整型
				if srcVal.Bool() {
					tempTargetVal.SetInt(1)
				} else {
					tempTargetVal.SetInt(0)
				}
				conversionSuccessful = true
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if srcKind == reflect.String {
				valStr := srcVal.String()
				if u, err := strconv.ParseUint(valStr, 10, bits); err == nil {
					tempTargetVal.SetUint(u)
					conversionSuccessful = true
				} else {
					if f, errFloat := strconv.ParseFloat(valStr, 64); errFloat == nil {
						if f >= 0 && f == math.Trunc(f) {
							wholeNumberStr := fmt.Sprintf("%.0f", f)
							if uWhole, errWhole := strconv.ParseUint(wholeNumberStr, 10, bits); errWhole == nil {
								tempTargetVal.SetUint(uWhole)
								conversionSuccessful = true
							}
						}
					}
				}
			} else if isNumeric(srcKind) {
				f64Val := toFloat64(srcVal)
				if f64Val >= 0 && f64Val == math.Trunc(f64Val) {
					valToSet := uint64(f64Val)
					if !tempTargetVal.OverflowUint(valToSet) {
						tempTargetVal.SetUint(valToSet)
						conversionSuccessful = true
					}
				}
			} else if srcKind == reflect.Bool {
				if srcVal.Bool() {
					tempTargetVal.SetUint(1)
				} else {
					tempTargetVal.SetUint(0)
				}
				conversionSuccessful = true
			}
		case reflect.Float32, reflect.Float64:
			targetBits := 32
			if targetKind == reflect.Float64 {
				targetBits = 64
			}

			if srcKind == reflect.String {
				if f, err := strconv.ParseFloat(srcVal.String(), targetBits); err == nil {
					tempTargetVal.SetFloat(f)
					conversionSuccessful = true
				}
			} else if isNumeric(srcKind) {
				f64Val := toFloat64(srcVal)
				if targetKind == reflect.Float32 {
					if math.Abs(f64Val) > math.MaxFloat32 && f64Val != 0 { // Check for overflow
						// conversion fails
					} else {
						tempTargetVal.SetFloat(f64Val)
						conversionSuccessful = true
					}
				} else { // Float64
					tempTargetVal.SetFloat(f64Val)
					conversionSuccessful = true
				}
			}
		case reflect.Bool:
			if srcKind == reflect.String {
				if b, err := strconv.ParseBool(srcVal.String()); err == nil {
					tempTargetVal.SetBool(b)
					conversionSuccessful = true
				}
			} else if isNumeric(srcKind) {
				f64Val := toFloat64(srcVal)
				if f64Val != 0 {
					tempTargetVal.SetBool(true)
				} else {
					tempTargetVal.SetBool(false)
				}
				conversionSuccessful = true
			}
		}
		if conversionSuccessful {
			convertedVal = tempTargetVal
		}
	}

	// 处理的是Go语言内置的类型系统可以直接处理的简单情况：如int32转int64
	if !conversionSuccessful {

		if srcVal.Type().AssignableTo(targetType) { // 检查源值是否可以直接赋值目标类型
			// 可以直接赋值，无需转换
			convertedVal = srcVal
			conversionSuccessful = true
		} else if srcVal.Type().ConvertibleTo(targetType) { // 检查源值是否可以直接转换为目标类型
			// 检查是否是数值类型转换到整型/无符号整型
			isSrcNumeric := isNumeric(srcVal.Kind())
			isTargetIntOrUint := targetType.Kind() == reflect.Int || targetType.Kind() == reflect.Int8 ||
				targetType.Kind() == reflect.Int16 || targetType.Kind() == reflect.Int32 ||
				targetType.Kind() == reflect.Int64 || targetType.Kind() == reflect.Uint ||
				targetType.Kind() == reflect.Uint8 || targetType.Kind() == reflect.Uint16 ||
				targetType.Kind() == reflect.Uint32 || targetType.Kind() == reflect.Uint64

			// 如果不是数值转整型的情况，使用标准转换（此处绕过浮点转整形，避免小数点截断）
			if !(isSrcNumeric && isTargetIntOrUint) {
				convertedVal = srcVal.Convert(targetType)
				conversionSuccessful = true
			}
		}
	}

	// 如果成功赋值，则将值设置到目标字段
	if conversionSuccessful {
		if isDestPtr {
			if destVal.IsNil() {
				destVal.Set(reflect.New(targetType))
			}
			destVal.Elem().Set(convertedVal)
		} else {
			if destVal.CanSet() {
				destVal.Set(convertedVal)
			}
		}
	}
	return nil
}

// isNumeric 检查给定的reflect.Kind是否是数值类型
// 参数:
//
//	k reflect.Kind - 要检查的类型种类
//
// 返回值:
//
//	bool - 如果是数值类型返回true，否则返回false
//
// 支持的数值类型包括:
//   - 有符号整数: Int, Int8, Int16, Int32, Int64
//   - 无符号整数: Uint, Uint8, Uint16, Uint32, Uint64
//   - 浮点数: Float32, Float64
func isNumeric(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// toFloat64 将任意数值类型的反射值转换为float64
// 参数:
//
//	v reflect.Value - 需要转换的反射值
//
// 返回值:
//
//	float64 - 转换后的浮点数值，如果类型不支持则返回NaN
//
// 功能:
//  1. 处理所有整数类型(有符号/无符号)到float64的转换
//  2. 处理float32/float64类型的直接转换
//  3. 对于非数值类型返回math.NaN()
func toFloat64(v reflect.Value) float64 {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return v.Float()
	default:
		return math.NaN()
	}
}
