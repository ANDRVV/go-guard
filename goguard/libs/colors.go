package libs

import (
	"github.com/fatih/color"
)

var (
	Lmagenta = color.New(color.FgHiMagenta).SprintFunc()
	Magenta  = color.New(color.FgMagenta).SprintFunc()
	Green    = color.New(color.FgGreen).SprintFunc()
	Cyan     = color.New(color.FgHiCyan).SprintFunc()
	Red      = color.New(color.FgHiRed).SprintFunc()
	Yellow   = color.New(color.FgHiYellow).SprintFunc()
	Hgreen   = color.New(color.FgHiGreen).SprintFunc()
	Bblue    = color.New(color.FgHiBlue, color.Bold).SprintFunc()
	Ublue    = color.New(color.FgHiBlue, color.Underline).SprintFunc()
	Bblack   = color.New(color.FgHiBlack, color.Bold).SprintFunc()
	Bwhite   = color.New(color.FgHiWhite, color.Bold).SprintFunc()
	Ubwhite  = color.New(color.FgHiWhite, color.Bold, color.Underline).SprintFunc()
	Cerr     =  color.New(color.FgHiRed, color.Bold, color.Faint, color.Italic).SprintFunc()
	Itawhite =  color.New(color.FgWhite, color.Bold, color.Italic).SprintFunc()
)