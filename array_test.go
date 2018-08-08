package toolkit

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestVariadicToSlice(t *testing.T) {
	res := VariadicToSlice("a", "b", "c", "d", "e")
	assert.Equal(t, []interface{}{"a", "b", "c", "d", "e"}, *res)

	// t.Logf("%#v \n", *res)
	// =====> []interface {}{"a", "b", "c", "d", "e"}
}

func TestMapToSlice(t *testing.T) {
	data := map[string]interface{}{
		"Name":   "noval",
		"Age":    12,
		"IsMale": true,
	}
	res := MapToSlice(data)
	assert.Contains(t, res, "noval")
	assert.Contains(t, res, 12)
	assert.Contains(t, res, true)
	assert.Len(t, res, len(data))

	// t.Logf("%#v \n", res)
	// =====> []interface {}{"noval", 12, "true"}
}

func TestHasMemberUsingSliseOfInterface(t *testing.T) {
	data := []interface{}{"noval", "male", 12, true, []string{"ainur", "panjang"}}

	isExists1 := HasMember(data, "noval")
	assert.True(t, isExists1)
	// t.Logf("%#v \n", isExists1)
	// =====> true

	isExists2 := HasMember(data, "male")
	assert.True(t, isExists2)
	// t.Logf("%#v \n", isExists2)
	// =====> true

	isExists3 := HasMember(data, 12)
	assert.True(t, isExists3)
	// t.Logf("%#v \n", isExists3)
	// =====> true

	isExists4 := HasMember(data, []string{"ainur", "panjang"})
	assert.True(t, isExists4)
	// t.Logf("%#v \n", isExists4)
	// =====> true

	isExists5 := HasMember(data, "yoga")
	assert.False(t, isExists5)
	// t.Logf("%#v \n", isExists5)
	// =====> false
}

func TestHasMemberUsingSliceOfString(t *testing.T) {
	data := []string{"noval", "bagus", "eky"}

	isExists1 := HasMember(data, "bagus")
	assert.True(t, isExists1)
	// t.Logf("%#v \n", isExists1)
	// =====> true

	isExists2 := HasMember(data, "ainur")
	assert.False(t, isExists2)
	// t.Logf("%#v \n", isExists2)
	// =====> false
}

func TestHasMemberUsingSliseOfInterfaceWithVariousElements(t *testing.T) {
	data := []interface{}{
		map[string]string{"name": "noval"},
		[]interface{}{"noval", "bagus", "eky"},
		"aris",
	}

	isExists1 := HasMember(data, map[string]string{"name": "noval"})
	assert.True(t, isExists1)
	// t.Logf("%#v \n", isExists1)
	// =====> true

	isExists2 := HasMember(data, data[0])
	assert.True(t, isExists2)
	// t.Logf("%#v \n", isExists2)
	// =====> true

	isExists3 := HasMember(data, "panjang")
	assert.False(t, isExists3)
	// t.Logf("%#v \n", isExists3)
	// =====> false
}

func TestMemberIndexUsingSliseOfInterface(t *testing.T) {
	data := []interface{}{"noval", "male", 12, true, []string{"ainur", "panjang"}}

	isExists1, index1 := MemberIndex(data, "noval")
	assert.True(t, isExists1)
	assert.Equal(t, 0, index1)
	// t.Logf("is exists: %t, index: %d \n", isExists1, index1)
	// =====> is exists: true, index: 0

	isExists2, index2 := MemberIndex(data, "male")
	assert.True(t, isExists2)
	assert.Equal(t, 1, index2)
	// t.Logf("is exists: %t, index: %d \n", isExists2, index2)
	// =====> is exists: true, index: 1

	isExists3, index3 := MemberIndex(data, 12)
	assert.True(t, isExists3)
	assert.Equal(t, 2, index3)
	// t.Logf("is exists: %t, index: %d \n", isExists3, index3)
	// =====> is exists: true, index: 2

	isExists4, index4 := MemberIndex(data, []string{"ainur", "panjang"})
	assert.True(t, isExists4)
	assert.Equal(t, 4, index4)
	// t.Logf("is exists: %t, index: %d \n", isExists4, index4)
	// =====> is exists: true, index: 4

	isExists5, index5 := MemberIndex(data, "yoga")
	assert.False(t, isExists5)
	assert.Equal(t, -1, index5)
	// t.Logf("is exists: %t, index: %d \n", isExists5, index5)
	// =====> is exists: false, index: -1
}

func TestMemberIndexUsingSliceOfString(t *testing.T) {
	data := []string{"noval", "bagus", "eky"}

	isExists1, index1 := MemberIndex(data, "bagus")
	assert.True(t, isExists1)
	assert.Equal(t, 1, index1)
	// t.Logf("is exists: %t, index: %d \n", isExists1, index1)
	// =====> is exists: true, index: 1

	isExists2, index2 := MemberIndex(data, "ainur")
	assert.False(t, isExists2)
	assert.Equal(t, -1, index2)
	// t.Logf("is exists: %t, index: %d \n", isExists2, index2)
	// =====> is exists: false, index: -1
}

func TestMemberIndexUsingSliseOfInterfaceWithVariousElements(t *testing.T) {
	data := []interface{}{
		map[string]string{"name": "noval"},
		[]interface{}{"noval", "bagus", "eky"},
		"aris",
	}

	isExists1, index1 := MemberIndex(data, map[string]string{"name": "noval"})
	assert.True(t, isExists1)
	assert.Equal(t, 0, index1)
	// t.Logf("is exists: %t, index: %d \n", isExists1, index1)
	// =====> is exists: true, index: 0

	isExists2, index2 := MemberIndex(data, data[0])
	assert.True(t, isExists2)
	assert.Equal(t, 0, index2)
	// t.Logf("is exists: %t, index: %d \n", isExists2, index2)
	// =====> is exists: true, index: 0

	isExists3, index3 := MemberIndex(data, "panjang")
	assert.False(t, isExists3)
	assert.Equal(t, -1, index3)
	// t.Logf("is exists: %t, index: %d \n", isExists3, index3)
	// =====> is exists: false, index: -1
}

func TestToInterfaceArrayUsingSliceOfStringData(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e"}
	res := ToInterfaceArray(data)

	assert.Equal(t, []interface{}{"a", "b", "c", "d", "e"}, res)

	// t.Logf("%#v \n", res)
	// =====> []interface {}{"a", "b", "c", "d", "e"}
}

func TestToInterfaceArrayUsingSliceOfIntData(t *testing.T) {
	data := []int{1, 2, 3}
	res := ToInterfaceArray(data)

	assert.Equal(t, []interface{}{1, 2, 3}, res)

	// t.Logf("%#v \n", res)
	// =====> []interface {}{1, 2, 3}
}

func TestToInterfaceArrayUsingSliceOfMapData(t *testing.T) {
	data := []map[string]int{
		map[string]int{"noval": 12},
		map[string]int{"panjang": 13},
		map[string]int{"ainur": 14},
		map[string]int{"bagus": 15},
	}
	res := ToInterfaceArray(data)

	assert.Equal(t, []interface{}{
		map[string]int{"noval": 12},
		map[string]int{"panjang": 13},
		map[string]int{"ainur": 14},
		map[string]int{"bagus": 15},
	}, res)

	// t.Logf("%#v \n", res)
	// =====> []interface{}{
	//     map[string]int{"noval": 12},
	//     map[string]int{"panjang": 13},
	//     map[string]int{"ainur": 14},
	//     map[string]int{"bagus": 15},
	// }
}

func TestCompareString(t *testing.T) {
	res1 := Compare("noval", "noval", "$eq")
	assert.True(t, res1)
	// t.Logf("%#v \n", res1)
	// =====> true

	res2 := Compare("noval", "agung", "$ne")
	assert.True(t, res2)
	// t.Logf("%#v \n", res2)
	// =====> true

	res3 := Compare("noval", "agung", "$gt")
	assert.True(t, res3)
	// t.Logf("%#v \n", res3)
	// =====> true, because "n" is greater that "a"

	res4 := Compare("agung", "noval", "$lt")
	assert.True(t, res4)
	// t.Logf("%#v \n", res4)
	// =====> true, because "a" is lower that "n"

	res5 := Compare("noval", "agung", "$gte")
	assert.True(t, res5)
	// t.Logf("%#v \n", res5)
	// =====> true, because "n" is greater that "a"

	res6 := Compare("agung", "noval", "$lte")
	assert.True(t, res6)
	// t.Logf("%#v \n", res6)
	// =====> true, because "a" is lower that "n"
}

func TestCompareNumeric(t *testing.T) {
	res1 := Compare(12, 12, "$eq")
	assert.True(t, res1)
	// t.Logf("%#v \n", res1)
	// =====> true

	res2 := Compare(12, 13, "$ne")
	assert.True(t, res2)
	// t.Logf("%#v \n", res2)
	// =====> true

	res3 := Compare(12, 11, "$gt")
	assert.True(t, res3)
	// t.Logf("%#v \n", res3)
	// =====> true

	res4 := Compare(12, 13, "$lt")
	assert.True(t, res4)
	// t.Logf("%#v \n", res4)
	// =====> true

	res5 := Compare(12, 12, "$gte") && Compare(12, 11, "$gte")
	assert.True(t, res5)
	// t.Logf("%#v \n", res5)
	// =====> true

	res6 := Compare(12, 12, "$lte") && Compare(12, 13, "$lte")
	assert.True(t, res6)
	// t.Logf("%#v \n", res6)
	// =====> true

	res7 := Compare(12, float64(12), "$eq")
	assert.True(t, res7)
	// t.Logf("%#v \n", res7)
	// =====> true

	res8 := Compare(float32(12), float64(12), "$eq")
	assert.True(t, res8)
	// t.Logf("%#v \n", res8)
	// =====> true

	res9 := Compare(float32(12), float64(12.0001), "$eq")
	assert.False(t, res9)
	// t.Logf("%#v \n", res9)
	// =====> false

	res10 := Compare(int32(12), float64(12.0001), "$eq")
	assert.False(t, res10)
	// t.Logf("%#v \n", res10)
	// =====> false

	res11 := Compare(float64(12.00000000011), float64(12.00000000011), "$eq")
	assert.True(t, res11)
	// t.Logf("%#v \n", res11)
	// =====> true

	res12 := Compare(12, interface{}("noval"), "$eq")
	assert.False(t, res12)
	// t.Logf("%#v \n", res12)
	// =====> true
}

func TestCompareBool(t *testing.T) {
	res1 := Compare(true, true, "$eq")
	assert.True(t, res1)
	// t.Logf("%#v \n", res1)
	// =====> true

	res2 := Compare(true, true, "$ne")
	assert.False(t, res2)
	// t.Logf("%#v \n", res2)
	// =====> true

	res3 := Compare(true, false, "$eq")
	assert.False(t, res3)
	// t.Logf("%#v \n", res3)
	// =====> true

	res4 := Compare(true, false, "$ne")
	assert.True(t, res4)
	// t.Logf("%#v \n", res4)
	// =====> true
}

func TestCompareTime(t *testing.T) {
	time1, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	assert.NoError(t, err)
	time2, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	assert.NoError(t, err)

	res1 := Compare(time1, time2, "$eq")
	assert.True(t, res1)
	// t.Logf("%#v \n", res1)
	// =====> true

	time1, err = time.Parse(time.RFC3339, "2012-11-01T22:08:42+00:00")
	assert.NoError(t, err)
	time2, err = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	assert.NoError(t, err)

	res2 := Compare(time1, time2, "$eq")
	assert.False(t, res2)
	// t.Logf("%#v \n", res2)
	// =====> true

	res3 := Compare(time1, time2, "$gt") && Compare(time1, time2, "$gte")
	assert.True(t, res3)
	// t.Logf("%#v \n", res3)
	// =====> true

	time1, err = time.Parse(time.RFC3339, "2012-11-01T22:08:42+00:00")
	assert.NoError(t, err)
	time2, err = time.Parse(time.RFC3339, "2012-11-01T22:08:43+00:00")
	assert.NoError(t, err)

	res4 := Compare(time1, time2, "$lt") && Compare(time1, time2, "$lte")
	assert.True(t, res4)
	// t.Logf("%#v \n", res4)
	// =====> true
}
