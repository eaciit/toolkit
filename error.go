package toolkit

import (
	"fmt"
	"runtime/debug"
	"strings"
)

func Error(txt string) error {
	return fmt.Errorf(txt)
}

func Errorf(txt string, obj ...interface{}) error {
	return fmt.Errorf(txt, obj...)
}

func StackTrace(filters ...string) string {
	traceStr := string(debug.Stack())
	traces := strings.Split(traceStr, "\n")
	outs := []string{}
	for _, trace := range traces {
		if strings.Contains(trace, ":") {
			if len(filters) == 0 {
				outs = append(outs, trace)
			} else {
				for _, filter := range filters {
					if strings.Contains(trace, filter) {
						outs = append(outs, trace)
						break
					}
				}
			}
		}
	}
	return strings.Join(outs, "\n")
}
