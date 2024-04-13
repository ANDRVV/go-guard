package libs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"reflect"
	"os"
	"unsafe"
	"net/url"
	"unicode"
)

func GetTime(base bool) (string, int64) {
	var now time.Time = time.Now()
	var sixformat func(string)string = func(text string)string{if len(text) > 6 {return text[:6]} else {return fmt.Sprintf("%06s", text)}}
	if base {
		return fmt.Sprintf("%02s", strconv.Itoa(now.Hour())) + ":" +
		fmt.Sprintf("%02s", strconv.Itoa(now.Minute())) + ":" +
		fmt.Sprintf("%02s", strconv.Itoa(now.Second())) + "." +
		sixformat(strconv.Itoa(now.Nanosecond() / 1e3)), now.UnixMicro()
	} else {
		return Bwhite(" at ") + 
		Hgreen(fmt.Sprintf("%02s", strconv.Itoa(now.Hour())))   + Bwhite(":") +
		Hgreen(fmt.Sprintf("%02s", strconv.Itoa(now.Minute()))) + Bwhite(":") +
		Hgreen(fmt.Sprintf("%02s", strconv.Itoa(now.Second()))) + Bwhite(".") +
		Hgreen(sixformat(strconv.Itoa(now.Nanosecond() / 1e3))), now.UnixMicro()
	}
}

func GetPosx(file string, line int, ok bool) (string) {
	if ok {
		var slash string = "./"
		if runtime.GOOS == "windows" { slash = ".\\" }
		return Ubwhite(fmt.Sprintf("%s%s:%d", slash, filepath.Base(file), line))
	} else {
		return Ubwhite("UNK-LINE")
	}
}

func GetPrefix() (string) {
	return Bblue("GG-DEBUG") + Bblack(" -> ")
}

func GetStackoverflowLink(title string, base bool) (string) {
	var getTitle func(string)string = func(title string)string{
		var newtitle string = ""; 
		for _, t := range strings.Split(strings.ReplaceAll(strings.ReplaceAll(title, " ", "+"), ":", "") + "+go", "+") {
			if !strings.HasPrefix(t, "\"") && !strings.HasPrefix(t, "'") {
				newtitle += t + "+"
			}
		}
		return strings.ToLower(strings.TrimRight(newtitle, "+"))
	}
	if base {
		return "https://stackoverflow.com/search?q=" + getTitle(title)
	} else {
		if resp, err := http.Get(strings.ReplaceAll("https://api.stackexchange.com/2.3/search/advanced?order=desc&sort=relevance&q=TITLE_INPUT&site=stackoverflow", "TITLE_INPUT", getTitle(title))); err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			if body, err2 := io.ReadAll(resp.Body); err2 == nil {
				var jmap map[string]interface{}
				json.Unmarshal([]byte(string(body)), &jmap)
				if items, exist := jmap["items"].([]interface{}); exist && len(items) >= 1 {
					if firstItem, exist2 := items[0].(map[string]interface{}); exist2 {
						if link := firstItem["link"].(string); len(link) > 25 {
							return link
						}
					}
				}
			}
		}
	}
	return ""
}

func GetInterfaceDataType(value interface{}) (string) {
	var interfacetype string
	if value == nil {
		interfacetype = "nil"
	} else {
		var refValue reflect.Value = reflect.ValueOf(value)
		if refValue.Kind() == reflect.Func {
			interfacetype = "func"
		} else if refValue.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			interfacetype = "error"
		} else if refValue.Kind() == reflect.Invalid {
			interfacetype = ""
		} else {
			interfacetype = "*"
		}
	}
	return interfacetype
}

func Println(text string) {
	os.Stdout.Write(unsafe.Slice(unsafe.StringData(text + "\n"), len(text) + 1))
}

func GetValueListFromCall(value interface{}) ([]string, []string, string) {
	time1, timestamp1 := GetTime(true)
	var returnArgs []reflect.Value = reflect.ValueOf(value).Call(make([]reflect.Value, 0))
	time2, timestamp2 := GetTime(true)
	var values []string
	for _, val := range returnArgs {
		values = append(values, fmt.Sprintf("%v", val.Interface())) 
	}
	return values, []string{time1, time2}, strconv.FormatFloat(float64(time.UnixMicro(timestamp2 - timestamp1).UnixMicro()) / 1e3, 'f', -1, 64)
}

func MakeValuesList(values []string, baselength int) (string) {
	if len(values) > 1 {
		var vlist string = Bblack("{")
		for count, value := range values {
			vlist += Bwhite("'") + Green(value) + Bwhite("'") + Bwhite(",") + " "
			if len(value) >= 45 {
				vlist += "\n"
				if count != len(values) {
					if isLink(value) {
						vlist += strings.Repeat(" ", baselength + 12)
					} else {
						vlist += strings.Repeat(" ", baselength - 1)
					}
				}
			}
			
		}
		return strings.TrimRight(vlist, Bwhite(",") + " \n") + Bblack("}")
	} else {
		if len(values) < 1 {
			return Bwhite("''")
		}
		return Bwhite("'") + Green(values[0]) + Bwhite("'")
	}
}

func isLink(value string) (bool) {
	if u, err := url.Parse(value[10:len(value) - 10]); err == nil && u.Scheme != "" && u.Host != "" {
		return true
	}
	return false
}

func HighlightValue(value string) (string) {
	return Bblack("<") + value + Bblack(">")
}

func GetStrFromInterface(value interface{}, baselength int) (string) {
	return MakeValuesList([]string{strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v", reflect.ValueOf(value).Interface()), "[", ""), "]", "")}, baselength)
}

func colorizeObjects(jsonText string) (string) {
	var fJSON strings.Builder
	var inString bool
	for i := 0; i < len(jsonText); i++ {
		var char string = string(jsonText[i])
		if char == "\"" {
			if i > 0 && string(jsonText[i - 1]) != "\\" {
				inString = !inString
				fJSON.WriteString(Bwhite("'"))
				continue
			}
		}
		if i + 1 == len(jsonText) {
			if string(jsonText[i - 2]) != "\"" {
				fJSON.WriteString(Bwhite("'") + Bblack("}"))
				break
			}
		}
		if inString {
			if char == "," && string(jsonText[i + 1]) == "\"" {
				fJSON.WriteString(Bwhite("',") + " ")
				inString = false
			} else {
				fJSON.WriteString(Green(char))
			}
		} else {
			switch char {
				case ":", ",":
					if char == ":" && string(jsonText[i + 1]) != "\"" && string(jsonText[i + 1]) != "[" && string(jsonText[i + 1]) != "{" {
						fJSON.WriteString(Bwhite(char) + " " + Bwhite("'"))
						inString = true
					} else {
						fJSON.WriteString(Bwhite(char) + " ")
					}
				default:
					fJSON.WriteString(Bblack(char))
			}
		}
	}
	return fJSON.String()
}

func GetMapFromInterface(value interface{}, baselength int, kind string) (string) {
	if jsoniter, err := json.Marshal(value); err == nil {
		if kind == "struct" {
			var className string = func(cnames []string)string{return cnames[len(cnames) - 1]}(strings.Split(fmt.Sprintf("%#v", value)[:strings.Index(fmt.Sprintf("%#v", value), "{")], "."))
			return Bblack("<") + Magenta("struct") + Bblack("::") + Lmagenta(className) + Bblack(">") + colorizeObjects(string(jsoniter))
		}
		return Bblack("<") + Magenta("map") + Bblack(">") + colorizeObjects(string(jsoniter))
    }
	return GetErrorMSG("MapReadingOrIterating")
}

func GetStrListFromInterface(value interface{}, baselength int) (string) {
	var vref reflect.Value = reflect.ValueOf(value)
	var strValues []string
	for i := 0; i < vref.Len(); i++ {
		strValues = append(strValues, fmt.Sprintf("%v", vref.Index(i).Interface()))
	}
	return MakeValuesList(strValues, baselength)
}

func GetVarSourceName(sourceName string) string {
	sourceName = strings.TrimLeft(sourceName, "&")
	if unicode.IsDigit(rune(sourceName[0])) || strings.ContainsAny(sourceName, "!@#$%^&*()-+={}[]|\\;:'\",.<>/?") || sourceName == "nil" {
		return "object" 
	}
	return sourceName
}

func ErasePrefixAndColors(text string, baselength int) (string) {
	var toErase []string = []string{strings.Repeat(" ", 45), strings.Repeat(" ", baselength + 4), "\x1b[0;22;24m", "\x1b[97;1;4m", "\x1b[0;22;22;23m", "\x1b[93m", "\x1b[32m", "\x1b[0;22;23m", "\x1b[37;1;3m", "\x1b[32m", "\x1b[96m", "\x1b[0m", "\x1b[94;1m", "\x1b[0;22m", "\x1b[90;1m", "\x1b[97;1m", "\x1b[0;24m", "\x1b[32m62", "\x1b[92m", "\x1b[35m", "\x1b[94;4m", "\x1b[91;1;2;3m", "\x1b[91m", "\x1b[95m"}
	for _, e := range toErase {
		text = strings.ReplaceAll(text, e, "")
	}
	return "[" + text[baselength - 4:] + "]"
}