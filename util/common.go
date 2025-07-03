package util

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"path/filepath"
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
func DeReference(o any) any {
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

	var result any
	var d map[string]any
	err := json.Unmarshal([]byte(data), &d)
	if err != nil {
		zap.S().Warnf("json unmarshal failed: %v", err)
		return "", err
	}

	indexes := strings.Split(index, ".")
	for i, v := range indexes {
		result = d[v]
		if i < len(indexes)-1 {
			if m, ok := d[v].(map[string]any); ok {
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
			zap.S().Warnf("json unmarshal failed: %v", err)
			return err
		}
		return nil
	}

	var result any
	var d map[string]any
	err := json.Unmarshal([]byte(data), &d)
	if err != nil {
		zap.S().Warnf("json unmarshal failed: %v", err)
		return err
	}

	indexes := strings.Split(index, ".")
	for i, v := range indexes {
		result = d[v]
		if i < len(indexes)-1 {
			if m, ok := d[v].(map[string]any); ok {
				d = m
			} else {
				return fmt.Errorf("invalid index '%s'", indexes[i+1])
			}
		}
	}

	target := ToJson(result)
	err = json.Unmarshal([]byte(target), t)
	if err != nil {
		zap.S().Warnf("json unmarshal failed: %v", err)
		return err
	}
	return nil
}

// ToJson 将任意结构体或对象转换为JSON字符串
// 参数:
//
//	o any - 需要转换的对象，可以是任意类型
//
// 返回值:
//
//	string - 转换后的JSON字符串，如果转换失败则返回空字符串
func ToJson(o any) string {
	resultJson, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(resultJson)
}

// ToMap 将结构体转换为 map[string]any
// 参数:
//
//	obj any - 需要转换的结构体
//
// 返回值:
//
//	map[string]any - 转换后的 map，键为结构体字段的 JSON 标签或字段名，值为字段的值
func ToMap(obj any) map[string]any {
	m := make(map[string]any)
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

// GetRootPath 通过探测 go.mod 文件来智能确定项目根目录
func GetRootPath() string {
	// 尝试从当前工作目录向上查找 go.mod
	dir, err := os.Getwd()
	if err != nil {
		return getExecutableDir()
	}

	// 无限循环，向上查找
	for {
		// 检查当前目录下是否存在 go.mod
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		// 到达文件系统根目录，仍未找到
		if dir == filepath.Dir(dir) {
			break
		}

		dir = filepath.Dir(dir)
	}

	return getExecutableDir()
}

// getExecutableDir 获取可执行文件所在的目录
func getExecutableDir() string {
	exe, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("无法获取可执行文件路径: %v", err))
	}
	return filepath.Dir(exe)
}
