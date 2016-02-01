package toolkit

import (
	"math"
)

var signs string = "()^*/+-"

type formulaElement struct {
	a, b   float64
	fa, fb *formulaElement
	//hasSubFormula bool
	op string
}

func Formula(formulaTxt string, in M) (ret float64, e error) {
	var fm *formulaElement
	fm, e = parseFormula(formulaTxt)
	if e != nil {
		return
	}
	ret = fm.runFormula(in)
	return
}

func parseFormula(formulaTxt string) (*formulaElement, error) {
	return nil, nil
}

func (f *formulaElement) runFormula(in M) float64 {
	var a, b float64
	if f.fa == nil {
		a = f.a
	} else {
		a = f.fa.runFormula(in)
	}
	if f.fb == nil {
		b = f.b
	} else {
		b = f.fb.runFormula(in)
	}
	op := f.op

	if op == "^" {
		return math.Pow(a, b)
	} else if op == "*" {
		return a * b
	} else if op == "/" {
		if b == 0 {
			return 0
		}
		return a / b
	} else if op == "+" {
		return a + b
	} else if op == "-" {
		return a - b
	} else if a == 0 {
		return b
	} else if b == 0 {
		return a
	}
	return 0
}
