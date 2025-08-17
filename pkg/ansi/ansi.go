// Package ansi provides functionality for working with ANSI escape sequences
// for terminal text coloring and formatting.
package ansi

// Color codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"
)

// StyleMap maps style names to their ANSI escape codes
var StyleMap = map[string]string{
	"reset":     Reset,
	"bold":      Bold,
	"italic":    Italic,
	"underline": Underline,

	"black":   Black,
	"red":     Red,
	"green":   Green,
	"yellow":  Yellow,
	"blue":    Blue,
	"magenta": Magenta,
	"cyan":    Cyan,
	"white":   White,

	"brightblack":   BrightBlack,
	"brightred":     BrightRed,
	"brightgreen":   BrightGreen,
	"brightyellow":  BrightYellow,
	"brightblue":    BrightBlue,
	"brightmagenta": BrightMagenta,
	"brightcyan":    BrightCyan,
	"brightwhite":   BrightWhite,
}

// GetStyleCode returns the ANSI escape code for the specified style
// If the style is not found, it returns the Reset code
func GetStyleCode(style string) string {
	if code, ok := StyleMap[style]; ok {
		return code
	}
	return Reset
}

// Colorize applies the specified style to the text and adds a reset code at the end
func Colorize(text, style string) string {
	styleCode := GetStyleCode(style)
	return styleCode + text + Reset
}
