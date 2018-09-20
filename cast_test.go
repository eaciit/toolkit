package toolkit

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestValue(t *testing.T) {
	data := "hello"
	assert.Equal(t,
		reflect.Indirect(reflect.ValueOf(data)).Interface(),
		Value(data).Interface(),
	)
}

func TestKind(t *testing.T) {
	data := "hello"
	assert.Equal(t,
		reflect.Indirect(reflect.ValueOf(data)).Kind(),
		Value(data).Kind(),
	)
}

func TestToString(t *testing.T) {
	res1 := ToString("hello")
	assert.Equal(t, "hello", res1)

	res2 := ToString(12)
	assert.Equal(t, "12", res2)

	res3 := ToString(float64(12.123))
	assert.Contains(t, res3, "12.123")

	res4 := ToString(true)
	assert.Equal(t, "true", res4)

	res5 := ToString([]string{"a", "b", "c", "d"})
	assert.Equal(t, "[a b c d]", res5)

	res6 := ToString(interface{}(int64(32)))
	assert.Equal(t, "32", res6)

	res7 := ToString(uint(12))
	assert.Equal(t, "12", res7)

	res8 := ToString(map[string]int{"a": 1, "b": 2})
	assert.True(t, res8 == "map[b:2 a:1]" || res8 == "map[a:1 b:2]")

	data9 := "hello"
	res9 := ToString(&data9)
	assert.Empty(t, res9)

	var data10 *int
	res10 := ToString(data10)
	assert.Empty(t, res10)

	var data11 interface{}
	res11 := ToString(data11)
	assert.Empty(t, res11)
}

func TestSetDefaultDateFormat(t *testing.T) {
	t.Skip("in progress")
}

func TestDefaultDateFormat(t *testing.T) {
	t.Skip("in progress")
}

func TestDate2String(t *testing.T) {
	t.Skip("in progress")
}

func TestString2Date(t *testing.T) {
	t.Skip("in progress")
}

func TestToInt(t *testing.T) {
	res1 := ToInt(12, RoundingAuto)
	assert.Equal(t, 12, res1)

	res2 := ToInt("12", RoundingAuto)
	assert.Equal(t, 12, res2)

	res3 := ToInt(float64(12.00123), RoundingAuto)
	assert.Equal(t, 12, res3)

	res4 := ToInt(float64(12.999), RoundingAuto)
	assert.Equal(t, 13, res4)

	res5 := ToInt(float64(12.999), RoundingDown)
	assert.Equal(t, 12, res5)

	res6 := ToInt(float64(12.00123), RoundingUp)
	assert.Equal(t, 13, res6)

	res7 := ToInt(uint(12), "")
	assert.Equal(t, 12, res7)

	data8 := 12
	res8 := ToInt(&data8, RoundingAuto)
	assert.Empty(t, res8)

	res9 := ToInt("12.0012345", RoundingAuto)
	assert.Equal(t, 12, res9)
}

func TestToFloat32(t *testing.T) {
	res1 := ToFloat32(12, 5, RoundingAuto)
	assert.Equal(t, float32(12), res1)

	res2 := ToFloat32("12", 5, RoundingAuto)
	assert.Equal(t, float32(12), res2)

	res3 := ToFloat32(12.0021, 5, RoundingUp)
	assert.Equal(t, float32(12.0021), res3)

	res4 := ToFloat32(12.0021234567, 5, RoundingAuto)
	assert.Equal(t, float32(12.00212), res4)

	res5 := ToFloat32(12.0021234567, 6, RoundingDown)
	assert.Equal(t, float32(12.002123), res5)

	res6 := ToFloat32(12.0021234567, 7, RoundingDown)
	assert.Equal(t, float32(12.0021234), res6)

	data7 := 12
	res7 := ToFloat32(&data7, 5, RoundingAuto)
	assert.Empty(t, res7)

	res8 := ToFloat32(uint(12), 7, RoundingDown)
	assert.Equal(t, float32(12), res8)
}

func TestToFloat64(t *testing.T) {
	res1 := ToFloat64(12, 5, RoundingAuto)
	assert.Equal(t, float64(12), res1)

	res2 := ToFloat64("12", 5, RoundingAuto)
	assert.Equal(t, float64(12), res2)

	res3 := ToFloat64(12.0021, 5, RoundingUp)
	assert.Equal(t, float64(12.0021), res3)

	res4 := ToFloat64(12.0021234567, 5, RoundingAuto)
	assert.Equal(t, float64(12.00212), res4)

	res5 := ToFloat64(12.0021234567, 6, RoundingDown)
	assert.Equal(t, float64(12.002123), res5)

	res6 := ToFloat64(12.0021234567, 7, RoundingDown)
	assert.Equal(t, float64(12.0021234), res6)

	data7 := 12
	res7 := ToFloat64(&data7, 5, RoundingAuto)
	assert.Empty(t, res7)

	res8 := ToFloat64(uint(12), 7, RoundingDown)
	assert.Equal(t, float64(12), res8)
}

func TestRoundingAuto64(t *testing.T) {
	t.Skip("covered by TestToFloat64()")
}

func TestRoundingDown64(t *testing.T) {
	t.Skip("covered by TestToFloat64()")
}

func TestRoundingUp64(t *testing.T) {
	t.Skip("covered by TestToFloat64()")
}

func TestToDateUsingUnixData(t *testing.T) {
	time1, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	if err != nil {
		t.Fatal(err.Error())
	}

	time2 := ToDate(time1.Unix(), "")
	assert.Equal(t,
		time1.UTC().Format("2006-01-02 15:04:05"),
		time2.UTC().Format("2006-01-02 15:04:05"),
	)
}

func TestToDateUsingDateString(t *testing.T) {
	dateString := "2012-11-01T22:08:41+00:00"
	time1, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		t.Fatal(err.Error())
	}

	time2 := ToDate(dateString, time.RFC3339)
	assert.Equal(t,
		time1.UTC().Format("2006-01-02 15:04:05"),
		time2.UTC().Format("2006-01-02 15:04:05"),
	)
}

func TestToDateUsingTime(t *testing.T) {
	dateString := "2012-11-01T22:08:41+00:00"
	time1, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		t.Fatal(err.Error())
	}

	time2 := ToDate(time1, time.RFC3339)
	assert.Equal(t,
		time1.UTC().Format("2006-01-02 15:04:05"),
		time2.UTC().Format("2006-01-02 15:04:05"),
	)
}

func TestToDuration(t *testing.T) {
	res1 := ToDuration("12")
	assert.Equal(t, time.Duration(12)*time.Second, res1)

	res2 := ToDuration(12)
	assert.Equal(t, time.Duration(12)*time.Second, res2)
}
