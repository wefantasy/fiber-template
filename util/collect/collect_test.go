package collect

import (
	"app/util"
	"fmt"
	"testing"
)

func TestFilter(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		filterFn func(int) bool
		expected []int
	}{
		{
			name:  "empty slice",
			input: []int{},
			filterFn: func(i int) bool {
				return true
			},
			expected: []int{},
		},
		{
			name:  "filter even numbers",
			input: []int{1, 2, 3, 4, 5},
			filterFn: func(i int) bool {
				return i%2 == 0
			},
			expected: []int{2, 4},
		},
		{
			name:  "filter all elements",
			input: []int{1, 2, 3},
			filterFn: func(i int) bool {
				return true
			},
			expected: []int{1, 2, 3},
		},
		{
			name:  "filter no elements",
			input: []int{1, 2, 3},
			filterFn: func(i int) bool {
				return false
			},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(tt.input, tt.filterFn)
			if len(got) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(got))
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("at index %d: expected %v, got %v", i, tt.expected[i], got[i])
				}
			}
		})
	}
}

// 测试结构体
type User struct {
	Name       string
	Occupation string
	Country    string
}
type User2 struct {
	Name string
	Sex  string
}

func TestFilterByStruct(t *testing.T) {
	users := []User{
		{"John Doe", "gardener", "USA"},
		{"Paul Smith", "programmer", "Canada"},
		{"Lucia Mala", "teacher", "Slovakia"},
		{"Tomas Smutny", "programmer", "Slovakia"},
	}

	filter := User2{Name: "Lucia Mala", Sex: "Male"}

	filteredUsers := FilterByStruct(users, filter)
	fmt.Println(util.StructToJson(filteredUsers))
}
