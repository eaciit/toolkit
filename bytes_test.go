package toolkit

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToBytesJSON(t *testing.T) {
	data := map[string]interface{}{
		"Name":     "noval",
		"Projects": []string{"dbox", "knot"},
		"Age":      12,
		"IsMale":   true,
	}

	bts := ToBytes(data, "json")
	assert.Equal(t, `{"Age":12,"IsMale":true,"Name":"noval","Projects":["dbox","knot"]}`, string(bts))

	// t.Logf("%s \n", string(bts))
	// =====> {"Age":12,"IsMale":true,"Name":"noval","Projects":["dbox","knot"]}
}

func TestToBytesGOBAndFromBytesGOB(t *testing.T) {
	data := map[string]interface{}{
		"Name":     "noval",
		"Projects": []string{"dbox", "knot"},
		"Age":      12,
		"IsMale":   true,
	}

	bts := ToBytes(data, "gob")

	res := make(M)
	err := FromBytes(bts, "gob", &res)
	assert.NoError(t, err)

	for key, value := range res {
		switch key {
		case "Name":
			assert.True(t, data[key] == value.(string))
		case "Projects":
			assert.Len(t, value.([]string), len(data[key].([]string)))
			assert.True(t, value.([]string)[0] == "dbox" || value.([]string)[0] == "knot")
			assert.True(t, value.([]string)[1] == "dbox" || value.([]string)[1] == "knot")
		case "Age":
			assert.True(t, data[key] == value.(int))
		case "IsMale":
			assert.True(t, data[key] == value.(bool))
		}
	}

	// t.Logf("%#v \n", res)
	// =====> map[string]interface{}{"Name":"noval", "Projects":[]string{"dbox", "knot"}, "Age":12, "IsMale":true}
}

func TestToBytesWithEmptyEncoderId(t *testing.T) {
	data := map[string]interface{}{
		"Name":     "noval",
		"Projects": []string{"dbox", "knot"},
		"Age":      12,
		"IsMale":   true,
	}

	bts := ToBytes(data, "")
	assert.Equal(t, `{"Age":12,"IsMale":true,"Name":"noval","Projects":["dbox","knot"]}`, string(bts))

	// t.Logf("%s \n", string(bts))
	// =====> {"Age":12,"IsMale":true,"Name":"noval","Projects":["dbox","knot"]}
}

func TestToBytesWithInvalidEncoderId(t *testing.T) {
	data := map[string]interface{}{
		"Name":     "noval",
		"Projects": []string{"dbox", "knot"},
		"Age":      12,
		"IsMale":   true,
	}

	bts := ToBytes(data, "lalala")
	assert.Equal(t, "", string(bts))

	// t.Logf("%s \n", string(bts))
	// =====> {"Age":12,"IsMale":true,"Name":"noval","Projects":["dbox","knot"]}
}

func TestFromBytes(t *testing.T) {
	jsonString := `{"Age":12,"IsMale":true,"Name":"noval","Projects":["dbox","knot"]}`
	res := make(map[string]interface{})
	err := FromBytes([]byte(jsonString), "json", &res)
	assert.NoError(t, err)

	for key, value := range res {
		switch key {
		case "Name":
			assert.True(t, value == "noval")
		case "Projects":
			valueCasted := value.([]interface{})
			assert.Len(t, valueCasted, 2)
			assert.True(t, valueCasted[0] == "dbox" || valueCasted[0] == "knot")
			assert.True(t, valueCasted[1] == "dbox" || valueCasted[1] == "knot")
		case "Age":
			assert.True(t, fmt.Sprintf("%s", value) == "12")
		case "IsMale":
			assert.True(t, value == true)
		}
	}

	// t.Logf("%#v \n", res)
	// =====> map[string]interface {}{"Age":"12", "IsMale":true, "Name":"noval", "Projects":[]interface {}{"dbox", "knot"}}
}

func TestFromBytesWithEmptyID(t *testing.T) {
	jsonString := `{"Age":12,"IsMale":true,"Name":"noval","Projects":["dbox","knot"]}`
	res := make(map[string]interface{})
	err := FromBytes([]byte(jsonString), "", &res)
	assert.NoError(t, err)

	for key, value := range res {
		switch key {
		case "Name":
			assert.True(t, value == "noval")
		case "Projects":
			valueCasted := value.([]interface{})
			assert.Len(t, valueCasted, 2)
			assert.True(t, valueCasted[0] == "dbox" || valueCasted[0] == "knot")
			assert.True(t, valueCasted[1] == "dbox" || valueCasted[1] == "knot")
		case "Age":
			assert.True(t, fmt.Sprintf("%s", value) == "12")
		case "IsMale":
			assert.True(t, value == true)
		}
	}

	// t.Logf("%#v \n", res)
	// =====> map[string]interface {}{"Age":"12", "IsMale":true, "Name":"noval", "Projects":[]interface {}{"dbox", "knot"}}
}

func TestGetEncodeByte(t *testing.T) {
	data := map[string]interface{}{
		"Name":     "noval",
		"Projects": []string{"dbox", "knot"},
		"Age":      12,
		"IsMale":   true,
	}

	res := make(map[string]interface{})
	err := DecodeByte(GetEncodeByte(data), &res)
	assert.NoError(t, err)

	// t.Logf("%#v \n", res)
	// ======> map[string]interface {}{"Projects":[]string{"dbox", "knot"}, "Age":12, "IsMale":true, "Name":"noval"}
}
