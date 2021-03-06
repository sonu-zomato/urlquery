package urlquery

import (
	"fmt"
	"testing"
)

type BuilderChild struct {
	Description string `query:"desc"`
	Long        uint16 `query:",vip"`
	Height      int    `query:"-"`
}

type BuilderInfo struct {
	Id       int
	Name     string         `query:"name"`
	Child    BuilderChild   `query:"child"`
	ChildPtr *BuilderChild  `query:"childPtr"`
	Children []BuilderChild `query:"children"`
	Params   map[string]rune
	status   bool
	UintPtr  uintptr
}

func TestMarshal(t *testing.T) {
	data := getMockData()

	_, err := Marshal(data)

	if err != nil {
		t.Error(err)
	}
}

func TestMarshal_Struct(t *testing.T) {
	data := getMockData2()

	bytes, err := Marshal(data)

	if err != nil {
		t.Error(err)
		return
	}

	v := BuilderInfo{}
	err = Unmarshal(bytes, &v)
	if err != nil {
		t.Error(err)
		return
	}

	if v.Name != "child" || v.status != false || v.Child.Height != 0 || len(v.Children) != 5 {
		fmt.Println(v.Name, v.status, v.Child.Height, v.Children)
		t.Error("Marshal Unmarshal is not equal")
		return
	}
}

func TestMarshal_NilPtr_Struct(t *testing.T) {
	data := getMockData3()

	bytes, err := Marshal(data)

	if err != nil {
		t.Error(err)
		return
	}

	v := BuilderInfo{}
	err = Unmarshal(bytes, &v)
	if err != nil {
		t.Error(err)
		return
	}

	if v.Name != "child3" || v.status != false || v.Child.Height != 0 || len(v.Children) != 5 {
		fmt.Println(v.Name, v.status, v.Child.Height, v.Children)
		t.Error("Marshal Unmarshal is not equal")
		return
	}

	if v.ChildPtr != nil {
		t.Error("The child pointer should be nil not ", v.ChildPtr)
		return
	}
	if v.Params != nil {
		t.Error("The params map should be nil")
		return
	}
}

func TestMarshal_Slice(t *testing.T) {
	data := []string{"a", "b"}

	bytes, err := Marshal(data)

	if err != nil {
		t.Error(err)
	}

	if string(bytes) != "0=a&1=b" {
		t.Error("failed to Marchal slice")
	}
}

func TestMarshal_Array(t *testing.T) {
	data := [3]int32{10, 200, 50}

	bytes, err := Marshal(data)

	if err != nil {
		t.Error(err)
	}

	if string(bytes) != "0=10&1=200&2=50" {
		t.Error("failed to Marchal slice")
	}
}

type TestPoint struct {
	X, Y int
}

type TestCircle struct {
	TestPoint
	R int
}

func TestMarshal_AnonymousFields(t *testing.T) {
	data := &TestCircle{R: 1}
	data.TestPoint.X = 12
	data.TestPoint.Y = 13

	bytes, err := Marshal(data)
	if err != nil {
		t.Error(err)
	}

	if string(bytes) != "X=12&Y=13&R=1" {
		t.Error("failed to Marshal anonymous fields")
	}
}

func TestMarshal_DuplicateCall(t *testing.T) {
	d1 := BuilderChild{
		Description: "a",
		Long:        10,
	}

	encoder := NewEncoder()
	encoder.Marshal(d1)

	d2 := BuilderChild{
		Description: "bb",
		Long:        200,
	}
	bytes2, err := encoder.Marshal(d2)
	if err != nil {
		t.Error(err)
	}

	if string(bytes2) != "desc=bb&Long=200" {
		t.Error("failed to Marshal duplicate call")
	}
}

//BenchmarkMarshal-4     	  295726	     11902 ns/op
func BenchmarkMarshal(b *testing.B) {
	data := getMockData2()

	for i := 0; i < b.N; i++ {
		_, err := Marshal(data)
		if err != nil {
			b.Error(err)
		}
	}
}

func getMockData() map[string]interface{} {
	var (
		f32 = float32(1.2)
		f64 = float64(13.4343453535343242342)
		i8  = int8(3)
		i64 = int64(9999999 * 9999999)
		u64 = uint16(567)
	)
	return map[string]interface{}{
		"id":     1,
		"fit":    true,
		"vip":    false,
		"desc":   "测试",
		"f32":    f32,
		"f64":    f64,
		"int8":   i8,
		"int64":  i64,
		"uint16": u64,
		"map": map[interface{}]interface{}{
			"caption": "test",
			5:         []int{11, 22},
			"child":   getMockData2(),
		},
		"struct": BuilderInfo{
			Id:   222,
			Name: "test",
		},
	}
}

func getMockData2() BuilderInfo {
	return BuilderInfo{
		Name: "child",
		Children: []BuilderChild{
			{Description: "d1", Height: 180},
			{Description: "d2", Long: 140},
			{Description: "d4"},
			{Description: "d5", Long: 1, Height: 20},
			{Description: "d6"},
		},
		Child:    BuilderChild{Description: "c1", Height: 20},
		ChildPtr: &BuilderChild{Description: "cptr", Long: 14, Height: 220},
		Params: map[string]rune{
			"abc":      111,
			"bbb":      222,
			"whoIsWho": 344340,
		},
		status:  true,
		UintPtr: uintptr(222),
	}
}

func getMockData3() BuilderInfo {
	return BuilderInfo{
		Name: "child3",
		Children: []BuilderChild{
			{Description: "d31", Height: 180},
			{Description: "d32", Long: 140},
			{Description: "d34"},
			{Description: "d35", Long: 1, Height: 20},
			{Description: "d36"},
		},
		Child:    BuilderChild{},
		ChildPtr: nil,
		Params:   nil,
		status:   true,
		UintPtr:  uintptr(2222),
	}
}

func TestEmptyStruct(t *testing.T) {
	data := &TestCircle{}
	bytes, err := Marshal(data)
	if err != nil {
		t.Error(err)
	}

	if string(bytes) != "" {
		t.Error("failed to Marshal anonymous fields")
	}
}
