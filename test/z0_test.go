package toolkittest

import (
	. "github.com/eaciit/toolkit"
	"testing"
)

func TestID(t *testing.T) {
	type Data struct {
		ID         string
		Name       string
		Desciption string
	}

	data := &Data{"EACIIT", "EACIIT", "Enhance Enterprise Value"}
	Printf("Data before ID set:\n%s\n\n", JsonString(data))
	f, i := IdInfo(data)
	if f == "" {
		t.Errorf("Unable to find ID")
		return
	}
	Printf("Field: %s\nValue: %v\n\n", f, i)
	e := SetId(data, "EC")
	if e != nil {
		t.Errorf("Unable to set ID: " + e.Error())
		return
	}
	Printf("Data after ID set:\n%s\n\n", JsonString(data))
}
