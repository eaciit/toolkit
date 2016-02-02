package toolkit

import (
	"math"
	"strings"
)

var signs string = "()^*/+-"

type formulaItem struct {
	a, b   float64
	fa, fb *formulaItem
	//hasSubFormula bool
	op string
}

func Formula(formulaTxt string, in M) (ret float64, e error) {
	var fm *formulaItem
	fm, e = parseFormula(formulaTxt)
	if e != nil {
		return
	}
	ret = fm.runFormula(in)
	return
}

func parseFormula(formulaTxt string) (*formulaItem, error) {
	var parts []string
	txtLen := len(formulaTxt)
	tmp := ""
	inBracket := false
	for i := 0; i < txtLen; i++ {
		c := string(formulaTxt[i])
		if c == "(" && !inBracket {
			if tmp != "" {
				parts = append(parts, tmp)
			}
			inBracket = true
			tmp = ""
		} else if c == ")" && inBracket {
			parts = append(parts, tmp)
			inBracket = false
			tmp = ""
		} else if i == txtLen-1 {
			tmp += c
			parts = append(parts, tmp)
		} else {
			tmp += c
		}
	}

	/*
		= 3+2*5/2-2*6
		= +2*5/2-2*6+3
		= +2*5/2+(-2)*6+3
		= +5-12+3
		= -4

		= +2+3*6/2+7 = 18
		= +3*6/2+2+7 = 18
	*/
	if len(parts) == 1 {
		var fs1, fs2 []float64
		var ss1, ss2 []string

		txt := parts[0]
		if !(strings.HasPrefix(txt, "-") || strings.HasPrefix(txt, "+")) {
			txt = "+" + txt
		}
		vs, ss := Split(txt, []string{"^", "*", "/", "+", "-"})

		// translate - to +
		for i, s := range ss {
			if s == "-" {
				ss[i] = "+"
				vs[i+1] = "-" + vs[i+1]
			}
		}

		// sort the signs, all operand not + should be placed first
		for i, s := range ss {
			f := ToFloat64(vs[i], 4, RoundingAuto)
			if s != "+" {
				ss1 = append(ss1, s)
				fs1 = append(fs1, f)
			} else {
				ss2 = append(ss2, s)
				fs2 = append(fs2, f)
			}
		}
	}

	return nil, nil
}

func (f *formulaItem) runFormula(in M) float64 {
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
