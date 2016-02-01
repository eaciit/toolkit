package toolkit

import (
	"math"
)

var signs string = "()^*/+-"

func Formula(formulaTxt string, in M) float64 {
	ret := float64(0)
	return ret
}

func ParseFormula(formulaTxt string, in M) *formulaElement {
	return nil
}

type formulaElement struct {
	a, b   float64
	fa, fb *formulaElement
	//hasSubFormula bool
	op string
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
