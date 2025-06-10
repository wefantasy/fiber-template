package util

import (
	"app/log"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
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

// JsonIndex 从JSON字符串中获取指定索引的值
func JsonIndex(data string, index string) (string, error) {
	if len(index) == 0 {
		return data, nil
	}

	var result interface{}
	var d map[string]interface{}
	err := json.Unmarshal([]byte(data), &d)
	if err != nil {
		log.Warnf("json unmarshal failed: %v", err)
		return "", err
	}

	indexes := strings.Split(index, ".")
	for i, v := range indexes {
		result = d[v]
		if i < len(indexes)-1 {
			if m, ok := d[v].(map[string]interface{}); ok {
				d = m
			} else {
				return "", fmt.Errorf("invalid index '%s'", indexes[i+1])
			}
		}
	}
	return result.(string), nil
}

func JsonToStructWithIndex[T any](data string, index string, t *T) error {
	if len(data) == 0 {
		return fmt.Errorf("data is empty")
	}

	if len(index) == 0 {
		err := json.Unmarshal([]byte(data), t)
		if err != nil {
			log.Warnf("json unmarshal failed: %v", err)
			return err
		}
		return nil
	}

	var result interface{}
	var d map[string]interface{}
	err := json.Unmarshal([]byte(data), &d)
	if err != nil {
		log.Warnf("json unmarshal failed: %v", err)
		return err
	}

	indexes := strings.Split(index, ".")
	for i, v := range indexes {
		result = d[v]
		if i < len(indexes)-1 {
			if m, ok := d[v].(map[string]interface{}); ok {
				d = m
			} else {
				return fmt.Errorf("invalid index '%s'", indexes[i+1])
			}
		}
	}

	target := StructToJson(result)
	err = json.Unmarshal([]byte(target), t)
	if err != nil {
		log.Warnf("json unmarshal failed: %v", err)
		return err
	}
	return nil
}

// StructToJson 将任意结构体或对象转换为JSON字符串
// 参数:
//
//	o interface{} - 需要转换的对象，可以是任意类型
//
// 返回值:
//
//	string - 转换后的JSON字符串，如果转换失败则返回空字符串
func StructToJson(o interface{}) string {
	resultJson, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(resultJson)
}

// StructToMap 将结构体转换为 map[string]interface{}
// 参数:
//
//	obj interface{} - 需要转换的结构体
//
// 返回值:
//
//	map[string]interface{} - 转换后的 map，键为结构体字段的 JSON 标签或字段名，值为字段的值
func StructToMap(obj interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	v := reflect.ValueOf(obj)
	t := v.Type()

	// 处理指针类型
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		// 忽略不可导出字段
		if !value.CanInterface() {
			continue
		}
		// 优先使用 JSON 标签作为键名
		tag := field.Tag.Get("json")
		if tag == "" {
			tag = field.Name
		} else {
			tag = strings.Split(tag, ",")[0] // 处理 omitempty 等选项
		}
		m[tag] = value.Interface()
	}
	return m
}

// IsDigits 检查字符串是否只包含数字字符
func IsDigits(str string) bool {
	for _, c := range str {
		if c < '0' || c > '9' {
			return false
		}
	}
	return str != "" // 空字符串返回 false
}

// IsNumeric 检查字符串是否为数字（整数或浮点数）
func IsNumeric(str string) bool {
	// 先尝试转换为浮点数
	if _, err := strconv.ParseFloat(str, 64); err == nil {
		return true
	}
	// 再检查是否为纯整数（根据需求可选）
	return IsDigits(str)
}

func RandString(n int) string {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
