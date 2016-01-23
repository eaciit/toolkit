package toolkittest

import (
	. "github.com/eaciit/toolkit"
	"testing"
	"time"
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

func TestCast(t *testing.T) {
	var (
		f64        float64 = 405.21
		f64other   float64
		int64      int = 201
		int64other int

		date       time.Time
		dateOther  time.Time
		dateString string
	)

	date = time.Date(1980, 4, 1, 0, 0, 0, 0, time.UTC)
	dateOther = String2Date("1-April-1980", "d-MMMM-yyyy")
	if dateOther != date {
		t.Errorf("Strign2Date fail. Expect %v got %v", date, dateOther)
	} else {
		t.Logf("String2Date: %v", dateOther)
	}

	dateString = Date2String(date, "dd MMM yy")
	if dateString != "01 Apr 80" {
		t.Errorf("Date2String fail. Expect %s got %s", "01 Apr 80", dateString)
	} else {
		t.Logf("Date2String: %s", dateString)
	}

	int64other = ToInt(f64, RoundingAuto)
	if int64other != 405 {
		t.Errorf("ToInt fail. Expect %d got %d", 405, int64other)
	}

	f64other = ToFloat64(&int64, 0, RoundingAuto) + 0.758
	if f64other != 201.758 {
		t.Errorf("ToFloat64 fail. Expect %f got %f", 201.758, f64other)
	}
}
