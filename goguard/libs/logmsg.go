package libs

import (
	"strings"
)

var (
	ElemERROR   func(string)string = func(varname string)string{return Red(varname) + Bwhite(": ")}
	ElemFUNC    func(string)string = func(funcname string)string{var lfuncname []string = strings.Split(funcname, "."); return Yellow(lfuncname[len(lfuncname) - 1]) + Bwhite("(): ")}
	ElemOBJECT  func(string)string = func(varname string)string{return Cyan(varname) + Bwhite(": ")}
	ElemNIL     func(string)string = func(varname string)string{return ElemOBJECT(varname) + Magenta("nil")}
	GetErrorMSG func(string)string = func(msg string)string{return Red("ERR-" + msg)}
)