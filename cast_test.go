package toolkit

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
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
}
