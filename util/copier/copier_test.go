package copier

import (
	"app/util"
	"testing"
)

// --- Example Usage ---
type SourceStruct struct {
	ID          int
	Name        string
	CountStr    *string // *string "abc"
	ValidStrInt *string // *string "123"
	FloatStrInt *string // *string "123.0"
	BadFloatStr *string // *string "123.45"
	SalaryFloat float32 // float32 789.5
	AgeInt      int     // int 42
	IntToStr    *int    // int 123
	Int16ToInt8 *int16  // int8 123
}

type DestStruct struct {
	ID          int64
	Name        *string
	CountStr    *int     // Expect nil from "abc"
	ValidStrInt *int     // Expect *int(123)
	FloatStrInt *int     // Expect *int(123) from "123.0"
	BadFloatStr *int     // Expect nil from "123.45"
	SalaryFloat *int     // Expect nil from float32 789.5 (due to fractional part)
	AgeInt      *float64 // Expect *float64(42.0)
	NonExistent *string  // Should remain nil
	IntToStr    *string  // Expect *string(123)
	Int16ToInt8 *int8    // int8 123
}

func Test_CopyProperties_Struct(t *testing.T) {
	strAbc := "abc"
	str123 := "123"
	str123Dot0 := "123.0"
	str123Dot45 := "123.45"
	int123 := 123
	var int16ToInt8 int16 = 123

	src := SourceStruct{
		ID:          1,
		Name:        "SourceItem",
		CountStr:    &strAbc,
		ValidStrInt: &str123,
		FloatStrInt: &str123Dot0,
		BadFloatStr: &str123Dot45,
		SalaryFloat: 789.5,
		AgeInt:      42,
		IntToStr:    &int123,
		Int16ToInt8: &int16ToInt8,
	}

	var dest DestStruct

	t.Logf("Before copy (dest): %s\n", util.StructToJson(&dest))
	t.Logf("Before copy (src): %s\n", util.StructToJson(&src))
	err := CopyProperties(&src, &dest)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("After copy (dest): %s\n", util.StructToJson(&dest))
	if dest.ID != 1 || *dest.Name != "SourceItem" || dest.CountStr != nil || *dest.ValidStrInt != 123 || *dest.FloatStrInt != 123 || dest.BadFloatStr != nil || dest.SalaryFloat != nil || *dest.AgeInt != 42.0 || dest.NonExistent != nil || *dest.IntToStr != "123" || *dest.Int16ToInt8 != 123 {
		t.Error("CopyProperties failed")
	}
}

func Test_CopyProperties_Map(t *testing.T) {
	mapSrc := map[string]any{
		"Name":        "MapSource",
		"CountStr":    "xyz",
		"ValidStrInt": "777",
		"FloatStrInt": "777.0",
		"BadFloatStr": "777.89",
		"SalaryFloat": float64(100.99),
		"AgeInt":      uint(25),
	}
	var destFromMap DestStruct
	t.Logf("Before copy (destFromMap): %+s\n", util.StructToJson(&destFromMap))
	t.Logf("Before copy (mapSrc): %+s\n", util.StructToJson(&mapSrc))
	CopyProperties(mapSrc, &destFromMap)
	t.Logf("After copy (destFromMap): %+s\n", util.StructToJson(&destFromMap))
	if destFromMap.ID != 0 || *destFromMap.Name != "MapSource" || destFromMap.CountStr != nil || *destFromMap.ValidStrInt != 777 || *destFromMap.FloatStrInt != 777 || destFromMap.BadFloatStr != nil || destFromMap.SalaryFloat != nil || *destFromMap.AgeInt != 25.0 {
		t.Error("CopyProperties failed")
	}
}

func TestCopyPropertiesSubMap(t *testing.T) {
	src := SourceStruct{
		ID:          123,
		Name:        "SourceItem",
		CountStr:    util.EnPointer("abc"),
		ValidStrInt: util.EnPointer("123"),
		SalaryFloat: 789.5,
	}
	dest := SourceStruct{
		ID:       456,
		Name:     "SourceItem",
		CountStr: util.EnPointer("zcv"),
	}
	err := CopyProperties(dest, &src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(util.StructToJson(&src))
}
