package toolkit

import (
	"errors"
	"math"
	"strings"
)

var signs string = "()^*/+-"
var signList []string = []string{"^", "*", "/", "+", "-"}

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

	//-- building the parts
	var parts []string
	var brackets []*formulaItem

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
			fi, efi := parseFormula(tmp)
			if efi != nil {
				return nil, errors.New("parseFormula: " + tmp + "," + efi.Error())
			}
			brackets := append(brackets, fi)
			parts = append(parts, Sprintf("@b_%d", len(brackets)-1))
			inBracket = false
			tmp = ""
		} else if i == txtLen-1 {
			tmp += c
			parts = append(parts, tmp)
		} else {
			tmp += c
		}
	}

	formulaTxt = ""
	for _, part := range parts {
		formulaTxt += part
	}
	if !(strings.HasPrefix(formulaTxt, "-") || strings.HasPrefix(formulaTxt, "+")) {
		formulaTxt = "+" + formulaTxt
	}

	var fsigns []string
	var fvalues []string
	txtLen = len(formulaTxt)
	tmp = ""
	for i := 0; i < txtLen; i++ {
		c := string(formulaTxt[txtLen-1-i])
		if HasMember(signList, c) {
			fsigns = append([]string{c}, fsigns...)
			fvalues = append([]string{tmp}, fvalues...)
			tmp = ""
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
	/*
		if len(parts) == 1 {
			txt := parts[0]
			if !(strings.HasPrefix(txt, "-") || strings.HasPrefix(txt, "+")) {
				txt = "+" + txt
			}
			vs, ss := Split(txt, []string{"^", "*", "/", "+", "-"})

			// group by the sign
			var itemParts []string
			tmp = ""
			for i, s := range ss {
				if tmp != "" {
					itemParts = append(itemParts, tmp)
				}
				if ss[i] == "+" || ss[i] == "-" {
					tmp = s
				} else {
					tmp += s
				}
				tmp += vs[i]
			}
		}
	*/

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
