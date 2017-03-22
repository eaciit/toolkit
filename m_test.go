package toolkit

import (
	"encoding/json"
	"fmt"
	"testing"
)

func ExampleM_PathGet() {
	input := []byte(`{
		"student": {
			"name": "John Doe",
			"class": {
				"name": "A5",
				"year": 2006,
				"major": {
					"name": "Biomedic"
				}
			}
		}
	}`)

	var obj M
	err := json.Unmarshal(input, &obj)
	if err != nil {
		fmt.Println("error: ", err)
	}

	majorName, err := obj.PathGet("student.class.major.name")
	if err != nil {
		fmt.Println("error: ", err)
	}

	fmt.Printf("%v", majorName)
	// Output:
	// Biomedic
}

func TestMPathGetInvalidPath(t *testing.T) {
	obj := M{
		"student": M{
			"name": "Charlie",
		},
	}

	var result interface{}
	var err error
	result, err = obj.PathGet("student.class.major.name")
	if err == nil {
		t.Fail()
	}
	if result != nil {
		t.Fail()
	}

	result, err = obj.PathGet("student.name.fullname")
	if err == nil {
		t.Fail()
	}
	if result != nil {
		t.Fail()
	}
}
