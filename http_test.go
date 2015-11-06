package toolkit

import (
	"fmt"
	"testing"
	"time"
)

func Test_Http(t *testing.T) {
	tstart := time.Now()
	url := "http://localhost:27123"

	for i := 0; i < 150; i++ {
		go func(i int) {
			r, e := HttpCall(url, "GET", nil, false, "", "")
			if e != nil {
				t.Errorf("(%d) Unable to call %s, got %s", i, url, e.Error())
				return
			}
			s := HttpContentString(r)
			if s == "" {
				t.Errorf("(%d) Got nothing", i)
			}
			//fmt.Println("Return: " + s)
		}(i)
	}

	time.Sleep(2 * time.Second)
	fmt.Printf("Done in %v \n", time.Since(tstart))
}
