package goguard

import (
	"fmt"
	"gg/goguard/libs"
	"reflect"
	"runtime"
	"strings"
	"github.com/shurcooL/go/reflectsource"
)

func GG(args ...interface{}) (file string, line int, timestamp int64, values string) {
	_, file, line, ok := runtime.Caller(1)
	strTime, timestamp := libs.GetTime(false)
	var base string = libs.GetPrefix() + libs.GetPosx(file, line, ok) + strTime
	if len(args) > 0 {
		base += libs.Bblack(" [")
		var baselength = len(base) - 149
		for counter, value := range args {
			var hasPointer bool = false
			var pointers []string
			var bslAdder int
			if counter > 0 {
				base += strings.Repeat(" ", baselength - 4)
			}
			var srcname string = libs.GetVarSourceName(reflectsource.GetParentArgExprAsString(uint32(counter)))
			for {
				switch libs.GetInterfaceDataType(value) {
					case "func":
						values, benchmark, timediff := libs.GetValueListFromCall(value)
						var time1, time2 string
						var changed bool
						bslAdder = 2
						for i := 0; i < 15; i++ {
							if benchmark[0][i] != benchmark[1][i] {
								time1 += libs.Red(string(benchmark[0][i]))
								time2 += libs.Red(string(benchmark[1][i]))
								changed = true
							} else {
								if string(benchmark[0][i]) == "." || string(benchmark[0][i]) == ":" {
									time1 += libs.Bwhite(string(benchmark[0][i]))
									time2 += libs.Bwhite(string(benchmark[1][i]))
								} else {
									if changed {
										time1 += libs.Red(string(benchmark[0][i]))
										time2 += libs.Red(string(benchmark[1][i]))
									} else {
										time1 += libs.Hgreen(string(benchmark[0][i]))
										time2 += libs.Hgreen(string(benchmark[1][i]))
									}
								}
							}
						}
						base += libs.ElemFUNC(srcname) + libs.MakeValuesList(values, baselength)
						base += "\n" + strings.Repeat(" ", baselength + 2) + libs.Itawhite("benchmark details ") + time1 + libs.Bwhite(" - ") + libs.Itawhite("start")
						base += "\n" + strings.Repeat(" ", baselength + 2) + libs.Itawhite("~" + timediff + "ms ") + strings.Repeat(" ", 14 - len(timediff)) + time2 + libs.Bwhite(" - ") + libs.Itawhite("end")
					case "error":
						var title string = value.(error).Error()
						var stacklink []string
						if rlink := libs.GetStackoverflowLink(title, false); rlink != "" {
							stacklink = append(stacklink, libs.Ublue(rlink))
						}
						stacklink = append(stacklink, libs.Ublue(libs.GetStackoverflowLink(title, true)))	
						var elem string = libs.ElemERROR(srcname)
						base += elem + libs.HighlightValue(libs.Cerr(title)) + "\n" + strings.Repeat(" ", baselength + len(srcname) - 2) + libs.Itawhite("search for ") + libs.MakeValuesList(stacklink, baselength + len(srcname) - 2) 
					case "nil":
						base += libs.ElemNIL(srcname)
					case "*":
							var reftype reflect.Kind = reflect.ValueOf(value).Kind()
							if reftype == reflect.Pointer {
								hasPointer = true
								pointers = append(pointers, "0x" + strings.ToUpper(fmt.Sprintf("%p", value)[2:]))
							}
							switch reftype {
								case reflect.Pointer:
									value = reflect.ValueOf(value).Elem().Interface()
									continue;
								case reflect.Slice:
									if reflect.ValueOf(value).Len() > 1 {
										base += libs.ElemOBJECT(srcname) + libs.GetStrListFromInterface(value, baselength + len(srcname))
									}
								case reflect.Map:
									base += libs.ElemOBJECT(srcname) + libs.GetMapFromInterface(value, baselength + len(srcname), "map")
								case reflect.Struct:
									base += libs.ElemOBJECT(srcname) + libs.GetMapFromInterface(value, baselength + len(srcname), "struct")
								default:
									base += libs.ElemOBJECT(srcname) + libs.GetStrFromInterface(value, baselength + len(srcname))
							}
					default:
						base += libs.ElemOBJECT(srcname) + libs.GetErrorMSG("DataTypeNotRecognized")
				}
				break
			}
			if hasPointer {
				base += "\n" + strings.Repeat(" ", baselength + len(srcname) - 2 + bslAdder) + libs.Lmagenta("@") + libs.Bblack("(") + func()string{
					var formatted string
					for _, pointer := range pointers {
						formatted += libs.Bwhite(pointer) + libs.Bblack("<-")
					}
					return strings.TrimRight(formatted, libs.Bblack("<-"))
				}() + libs.Bblack(")") 
				bslAdder = 0
			}
			base += libs.Bwhite(",") + strings.Repeat(" ", baselength) + "\n"
		}
		base = strings.TrimRight(strings.TrimRight(strings.TrimRight(base, "\n"), strings.Repeat(" ", baselength)), libs.Bwhite(",")) + libs.Bblack("]")
		values = libs.ErasePrefixAndColors(base, baselength)
	}
	libs.Println(base)
	return
}