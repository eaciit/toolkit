package toolkit

import (
	"errors"
	//	"math"
	"strings"
)

var signs string = "()^*/+-"
var signList []string = []string{"^", "*", "/", "+", "-"}

type formulaItem struct {
	Parm        string
	Value       float64
	Negate      bool
	Divide      bool
	SubFormulas []*formulaItem
	BaseOp      string

	Formula string
}

func Formula(formulaTxt string, in M) (ret float64, e error) {
	var fm *formulaItem
	fm, e = parseFormula(formulaTxt)
	if e != nil {
		return
	}
	ret = fm.run(in)
	return
}

func parseFormula(txt string) (*formulaItem, error) {
	return parseFormulaSub(txt, []*formulaItem{})
}

func parseFormulaSub(formulaTxt string, fisubs []*formulaItem) (*formulaItem, error) {

	//-- building the parts
	if formulaTxt == "" {
		return nil, errors.New("parseFormula:Formula is empty")
	}

	var parts []string
	originalFormula := formulaTxt

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
			fi, efi := parseFormulaSub(tmp, fisubs)
			if efi != nil {
				return nil, errors.New("parseFormula: " + tmp + "," + efi.Error())
			} else if fi == nil {
				return nil, errors.New(Sprintf("parseFormula: %s unable to parse it into *formulaItem", tmp))
			}
			fisubs = append(fisubs, fi)
			parts = append(parts, Sprintf("@b_%d", len(fisubs)-1))
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
	//var fsigns []string
	//var fvalues []string

	var fparts []string
	txtLen = len(formulaTxt)
	tmp = ""
	for i := 0; i < txtLen; i++ {
		c := string(formulaTxt[txtLen-1-i])
		tmp = c + tmp
		if HasMember([]string{"+", "-"}, c) {
			fparts = append(fparts, tmp)
			tmp = ""
		}
	}

	ret := new(formulaItem)
	ret.Formula = originalFormula
	if len(fparts) == 1 {
		fpart := fparts[0]
		multiVariable := false
		if splis, _ := Split(fpart, signList); len(splis) > 1 {
			multiVariable = true
		}
		ret.BaseOp = "+"
		if isMultiply(fpart) {
			ret.BaseOp = "*"
		}

		tmp := ""
		fpartlen := len(fpart)
		for fpartidx := fpartlen - 1; fpartidx >= 0; fpartidx-- {
			c := string(fpart[fpartidx])
			if ((c == "*" || c == "/" || c == "^") && ret.BaseOp == "*") ||
				((c == "+" || c == "-") && ret.BaseOp == "+") || fpartidx == 0 {
				var subfi *formulaItem
				isnegate := false
				if fpartidx == 0 && c == "-" {
					tmp = c + tmp
				}
				if strings.HasPrefix(tmp, "-") {
					isnegate = true
					tmp = tmp[1:]
					Println("Negating: ", tmp, " ")
				}
				if strings.Contains(tmp, "@b_") {
					//--- it is a subfunction that already defined
					formulaIndex := int(-1)
					formulaIndex = ToInt(tmp[3:], RoundingAuto)
					if formulaIndex >= len(fisubs) {
						return nil, errors.New(Sprintf("parseFormula: %s Subformula index-%d is not available", tmp, formulaIndex))
					}
					subfi = fisubs[formulaIndex]
				} else {
					//--- not a subfunction already defined
					subfi = new(formulaItem)
					if !strings.Contains(tmp, "@") {
						//-- it is a value
						f64 := ToFloat64(tmp, 4, RoundingAuto)
						if tmp != "0" && f64 == float64(0) {
							return nil, errors.New("parseFormula: " + tmp + " Can not render to float")
						}
						//Printf("%s value is %.2f\n", tmp, f64)
						subfi.Value = f64
						//ret.Value = f64
					} else {
						//-- it is a.Parm
						subfi.Parm = tmp
						//ret.Parm = tmp
					}
				}

				if multiVariable {
					subfi.Divide = c == "/"
					subfi.Negate = isnegate
					ret.SubFormulas = append(ret.SubFormulas, subfi)
				} else {
					ret.Divide = c == "/"
					ret.Negate = isnegate
					ret.Parm = subfi.Parm
					ret.Value = subfi.Value
					ret.SubFormulas = subfi.SubFormulas
				}
				tmp = ""
			} else {
				tmp = c + tmp
			}
		}
	} else {
		for _, fpart := range fparts {
			fi, efi := parseFormulaSub(fpart, fisubs)
			if efi != nil {
				return nil, errors.New("parseFormula: " + fpart + " " + efi.Error())
			}
			if isMultiply(fpart) {
				ret.BaseOp = "*"
			} else {
				ret.BaseOp = "+"
			}
			ret.SubFormulas = append(ret.SubFormulas, fi)
		}
	}

	/*
		Printf("Origin:%s NewFormula:%s FParts:%s\n%s\n\n",
			originalFormula, formulaTxt,
			JsonString(fparts),
			JsonStringIndent(ret, ""))
	*/
	return ret, nil
}

func isMultiply(fpart string) bool {
	if strings.Contains(fpart, "*") || strings.Contains(fpart, "/") || strings.Contains(fpart, "^") {
		return true
	}
	return false
}

func (f *formulaItem) run(in M) float64 {
	var ret float64
	dbg := ""
	if f.BaseOp == "+" {
		ret = float64(0)
	} else {
		ret = float64(1)
	}
	if len(f.SubFormulas) == 0 {
		if f.Parm != "" {
			ret = float64(in.GetFloat64(f.Parm))
		} else {
			ret = f.Value
		}
		if f.Negate {
			ret = -ret
		}
	} else {
		for idx, sf := range f.SubFormulas {
			v := sf.run(in)
			if sf.Divide {
				v = 1.0 / v
			}
			if f.BaseOp == "+" {
				ret += v
			} else {
				ret = ret * v
				//ret = 2016.0 * 0.1
			}
			dbg += Sprintf("%d=%.2f ", idx, v)
		}
	}
	Printf("Formula: %s Value: %.2f Negate:%v BaseOp:%s Trace:%s\n", f.Formula, ret, f.Negate, f.BaseOp, dbg)
	return ret
}
