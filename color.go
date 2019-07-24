package colour

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

var (
	// NoColour defines if the output is colourized or not. It's dynamically set to
	// false or true based on the stdout's file descriptor referring to a terminal
	// or not. This is a global option and affects all colours. For more control
	// over each colour block use the methods DisableColour() individually.
	NoColour = os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))

	// Output defines the standard output of the print functions. By default
	// os.Stdout is used.
	Output = colorable.NewColorableStdout()

	// Error defines a colour supporting writer for os.Stderr.
	Error = colorable.NewColorableStderr()

	// coloursCache is used to reduce the count of created Colour objects and
	// allows to reuse already created objects with required Attribute.
	coloursCache   = make(map[Attribute]*Colour)
	coloursCacheMu sync.Mutex // protects coloursCache
)

// Colour defines a custom colour object which is defined by SGR parameters.
type Colour struct {
	params   []Attribute
	noColour *bool
}

// Attribute defines a single SGR Code
type Attribute int

const escape = "\x1b"

// Base attributes
const (
	Reset Attribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground text colours
const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity text colours
const (
	FgHiBlack Attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colours
const (
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colours
const (
	BgHiBlack Attribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

// New returns a newly created colour object.
func New(value ...Attribute) *Colour {
	c := &Colour{params: make([]Attribute, 0)}
	c.Add(value...)
	return c
}

// Set sets the given parameters immediately. It will change the colour of
// output with the given SGR parameters until colour.Unset() is called.
func Set(p ...Attribute) *Colour {
	c := New(p...)
	c.Set()
	return c
}

// Unset resets all escape attributes and clears the output. Usually should
// be called after Set().
func Unset() {
	if NoColour {
		return
	}

	fmt.Fprintf(Output, "%s[%dm", escape, Reset)
}

// Set sets the SGR sequence.
func (c *Colour) Set() *Colour {
	if c.isNoColourSet() {
		return c
	}

	fmt.Fprintf(Output, c.format())
	return c
}

func (c *Colour) unset() {
	if c.isNoColourSet() {
		return
	}

	Unset()
}

func (c *Colour) setWriter(w io.Writer) *Colour {
	if c.isNoColourSet() {
		return c
	}

	fmt.Fprintf(w, c.format())
	return c
}

func (c *Colour) unsetWriter(w io.Writer) {
	if c.isNoColourSet() {
		return
	}

	if NoColour {
		return
	}

	fmt.Fprintf(w, "%s[%dm", escape, Reset)
}

// Add is used to chain SGR parameters. Use as many as parameters to combine
// and create custom colour objects. Example: Add(colour.FgRed, colour.Underline).
func (c *Colour) Add(value ...Attribute) *Colour {
	c.params = append(c.params, value...)
	return c
}

func (c *Colour) prepend(value Attribute) {
	c.params = append(c.params, 0)
	copy(c.params[1:], c.params[0:])
	c.params[0] = value
}

// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
// On Windows, users should wrap w with colorable.NewColorable() if w is of
// type *os.File.
func (c *Colour) Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	c.setWriter(w)
	defer c.unsetWriter(w)

	return fmt.Fprint(w, a...)
}

// Print formats using the default formats for its operands and writes to
// standard output. Spaces are added between operands when neither is a
// string. It returns the number of bytes written and any write error
// encountered. This is the standard fmt.Print() method wrapped with the given
// colour.
func (c *Colour) Print(a ...interface{}) (n int, err error) {
	c.Set()
	defer c.unset()

	return fmt.Fprint(Output, a...)
}

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
// On Windows, users should wrap w with colorable.NewColorable() if w is of
// type *os.File.
func (c *Colour) Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	c.setWriter(w)
	defer c.unsetWriter(w)

	return fmt.Fprintf(w, format, a...)
}

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
// This is the standard fmt.Printf() method wrapped with the given colour.
func (c *Colour) Printf(format string, a ...interface{}) (n int, err error) {
	c.Set()
	defer c.unset()

	return fmt.Fprintf(Output, format, a...)
}

// Fprintln formats using the default formats for its operands and writes to w.
// Spaces are always added between operands and a newline is appended.
// On Windows, users should wrap w with colorable.NewColorable() if w is of
// type *os.File.
func (c *Colour) Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	c.setWriter(w)
	defer c.unsetWriter(w)

	return fmt.Fprintln(w, a...)
}

// Println formats using the default formats for its operands and writes to
// standard output. Spaces are always added between operands and a newline is
// appended. It returns the number of bytes written and any write error
// encountered. This is the standard fmt.Print() method wrapped with the given
// colour.
func (c *Colour) Println(a ...interface{}) (n int, err error) {
	c.Set()
	defer c.unset()

	return fmt.Fprintln(Output, a...)
}

// Sprint is just like Print, but returns a string instead of printing it.
func (c *Colour) Sprint(a ...interface{}) string {
	return c.wrap(fmt.Sprint(a...))
}

// Sprintln is just like Println, but returns a string instead of printing it.
func (c *Colour) Sprintln(a ...interface{}) string {
	return c.wrap(fmt.Sprintln(a...))
}

// Sprintf is just like Printf, but returns a string instead of printing it.
func (c *Colour) Sprintf(format string, a ...interface{}) string {
	return c.wrap(fmt.Sprintf(format, a...))
}

// FprintFunc returns a new function that prints the passed arguments as
// colourized with colour.Fprint().
func (c *Colour) FprintFunc() func(w io.Writer, a ...interface{}) {
	return func(w io.Writer, a ...interface{}) {
		c.Fprint(w, a...)
	}
}

// PrintFunc returns a new function that prints the passed arguments as
// colourized with colour.Print().
func (c *Colour) PrintFunc() func(a ...interface{}) {
	return func(a ...interface{}) {
		c.Print(a...)
	}
}

// FprintfFunc returns a new function that prints the passed arguments as
// colourized with colour.Fprintf().
func (c *Colour) FprintfFunc() func(w io.Writer, format string, a ...interface{}) {
	return func(w io.Writer, format string, a ...interface{}) {
		c.Fprintf(w, format, a...)
	}
}

// PrintfFunc returns a new function that prints the passed arguments as
// colourized with colour.Printf().
func (c *Colour) PrintfFunc() func(format string, a ...interface{}) {
	return func(format string, a ...interface{}) {
		c.Printf(format, a...)
	}
}

// FprintlnFunc returns a new function that prints the passed arguments as
// colourized with colour.Fprintln().
func (c *Colour) FprintlnFunc() func(w io.Writer, a ...interface{}) {
	return func(w io.Writer, a ...interface{}) {
		c.Fprintln(w, a...)
	}
}

// PrintlnFunc returns a new function that prints the passed arguments as
// colourized with colour.Println().
func (c *Colour) PrintlnFunc() func(a ...interface{}) {
	return func(a ...interface{}) {
		c.Println(a...)
	}
}

// SprintFunc returns a new function that returns colourized strings for the
// given arguments with fmt.Sprint(). Useful to put into or mix into other
// string. Windows users should use this in conjunction with colour.Output, example:
//
//	put := New(FgYellow).SprintFunc()
//	fmt.Fprintf(colour.Output, "This is a %s", put("warning"))
func (c *Colour) SprintFunc() func(a ...interface{}) string {
	return func(a ...interface{}) string {
		return c.wrap(fmt.Sprint(a...))
	}
}

// SprintfFunc returns a new function that returns colourized strings for the
// given arguments with fmt.Sprintf(). Useful to put into or mix into other
// string. Windows users should use this in conjunction with colour.Output.
func (c *Colour) SprintfFunc() func(format string, a ...interface{}) string {
	return func(format string, a ...interface{}) string {
		return c.wrap(fmt.Sprintf(format, a...))
	}
}

// SprintlnFunc returns a new function that returns colourized strings for the
// given arguments with fmt.Sprintln(). Useful to put into or mix into other
// string. Windows users should use this in conjunction with colour.Output.
func (c *Colour) SprintlnFunc() func(a ...interface{}) string {
	return func(a ...interface{}) string {
		return c.wrap(fmt.Sprintln(a...))
	}
}

// sequence returns a formatted SGR sequence to be plugged into a "\x1b[...m"
// an example output might be: "1;36" -> bold cyan
func (c *Colour) sequence() string {
	format := make([]string, len(c.params))
	for i, v := range c.params {
		format[i] = strconv.Itoa(int(v))
	}

	return strings.Join(format, ";")
}

// wrap wraps the s string with the colours attributes. The string is ready to
// be printed.
func (c *Colour) wrap(s string) string {
	if c.isNoColourSet() {
		return s
	}

	return c.format() + s + c.unformat()
}

func (c *Colour) format() string {
	return fmt.Sprintf("%s[%sm", escape, c.sequence())
}

func (c *Colour) unformat() string {
	return fmt.Sprintf("%s[%dm", escape, Reset)
}

// DisableColour disables the colour output. Useful to not change any existing
// code and still being able to output. Can be used for flags like
// "--no-colour". To enable back use EnableColour() method.
func (c *Colour) DisableColour() {
	c.noColour = boolPtr(true)
}

// EnableColour enables the colour output. Use it in conjunction with
// DisableColour(). Otherwise this method has no side effects.
func (c *Colour) EnableColour() {
	c.noColour = boolPtr(false)
}

func (c *Colour) isNoColourSet() bool {
	// check first if we have user setted action
	if c.noColour != nil {
		return *c.noColour
	}

	// if not return the global option, which is disabled by default
	return NoColour
}

// Equals returns a boolean value indicating whether two colours are equal.
func (c *Colour) Equals(c2 *Colour) bool {
	if len(c.params) != len(c2.params) {
		return false
	}

	for _, attr := range c.params {
		if !c2.attrExists(attr) {
			return false
		}
	}

	return true
}

func (c *Colour) attrExists(a Attribute) bool {
	for _, attr := range c.params {
		if attr == a {
			return true
		}
	}

	return false
}

func boolPtr(v bool) *bool {
	return &v
}

func getCachedColour(p Attribute) *Colour {
	coloursCacheMu.Lock()
	defer coloursCacheMu.Unlock()

	c, ok := coloursCache[p]
	if !ok {
		c = New(p)
		coloursCache[p] = c
	}

	return c
}

func colourPrint(format string, p Attribute, a ...interface{}) {
	c := getCachedColour(p)

	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}

	if len(a) == 0 {
		c.Print(format)
	} else {
		c.Printf(format, a...)
	}
}

func colourString(format string, p Attribute, a ...interface{}) string {
	c := getCachedColour(p)

	if len(a) == 0 {
		return c.SprintFunc()(format)
	}

	return c.SprintfFunc()(format, a...)
}

// Black is a convenient helper function to print with black foreground. A
// newline is appended to format by default.
func Black(format string, a ...interface{}) { colourPrint(format, FgBlack, a...) }

// Red is a convenient helper function to print with red foreground. A
// newline is appended to format by default.
func Red(format string, a ...interface{}) { colourPrint(format, FgRed, a...) }

// Green is a convenient helper function to print with green foreground. A
// newline is appended to format by default.
func Green(format string, a ...interface{}) { colourPrint(format, FgGreen, a...) }

// Yellow is a convenient helper function to print with yellow foreground.
// A newline is appended to format by default.
func Yellow(format string, a ...interface{}) { colourPrint(format, FgYellow, a...) }

// Blue is a convenient helper function to print with blue foreground. A
// newline is appended to format by default.
func Blue(format string, a ...interface{}) { colourPrint(format, FgBlue, a...) }

// Magenta is a convenient helper function to print with magenta foreground.
// A newline is appended to format by default.
func Magenta(format string, a ...interface{}) { colourPrint(format, FgMagenta, a...) }

// Cyan is a convenient helper function to print with cyan foreground. A
// newline is appended to format by default.
func Cyan(format string, a ...interface{}) { colourPrint(format, FgCyan, a...) }

// White is a convenient helper function to print with white foreground. A
// newline is appended to format by default.
func White(format string, a ...interface{}) { colourPrint(format, FgWhite, a...) }

// BlackString is a convenient helper function to return a string with black
// foreground.
func BlackString(format string, a ...interface{}) string { return colourString(format, FgBlack, a...) }

// RedString is a convenient helper function to return a string with red
// foreground.
func RedString(format string, a ...interface{}) string { return colourString(format, FgRed, a...) }

// GreenString is a convenient helper function to return a string with green
// foreground.
func GreenString(format string, a ...interface{}) string { return colourString(format, FgGreen, a...) }

// YellowString is a convenient helper function to return a string with yellow
// foreground.
func YellowString(format string, a ...interface{}) string { return colourString(format, FgYellow, a...) }

// BlueString is a convenient helper function to return a string with blue
// foreground.
func BlueString(format string, a ...interface{}) string { return colourString(format, FgBlue, a...) }

// MagentaString is a convenient helper function to return a string with magenta
// foreground.
func MagentaString(format string, a ...interface{}) string {
	return colourString(format, FgMagenta, a...)
}

// CyanString is a convenient helper function to return a string with cyan
// foreground.
func CyanString(format string, a ...interface{}) string { return colourString(format, FgCyan, a...) }

// WhiteString is a convenient helper function to return a string with white
// foreground.
func WhiteString(format string, a ...interface{}) string { return colourString(format, FgWhite, a...) }

// HiBlack is a convenient helper function to print with hi-intensity black foreground. A
// newline is appended to format by default.
func HiBlack(format string, a ...interface{}) { colourPrint(format, FgHiBlack, a...) }

// HiRed is a convenient helper function to print with hi-intensity red foreground. A
// newline is appended to format by default.
func HiRed(format string, a ...interface{}) { colourPrint(format, FgHiRed, a...) }

// HiGreen is a convenient helper function to print with hi-intensity green foreground. A
// newline is appended to format by default.
func HiGreen(format string, a ...interface{}) { colourPrint(format, FgHiGreen, a...) }

// HiYellow is a convenient helper function to print with hi-intensity yellow foreground.
// A newline is appended to format by default.
func HiYellow(format string, a ...interface{}) { colourPrint(format, FgHiYellow, a...) }

// HiBlue is a convenient helper function to print with hi-intensity blue foreground. A
// newline is appended to format by default.
func HiBlue(format string, a ...interface{}) { colourPrint(format, FgHiBlue, a...) }

// HiMagenta is a convenient helper function to print with hi-intensity magenta foreground.
// A newline is appended to format by default.
func HiMagenta(format string, a ...interface{}) { colourPrint(format, FgHiMagenta, a...) }

// HiCyan is a convenient helper function to print with hi-intensity cyan foreground. A
// newline is appended to format by default.
func HiCyan(format string, a ...interface{}) { colourPrint(format, FgHiCyan, a...) }

// HiWhite is a convenient helper function to print with hi-intensity white foreground. A
// newline is appended to format by default.
func HiWhite(format string, a ...interface{}) { colourPrint(format, FgHiWhite, a...) }

// HiBlackString is a convenient helper function to return a string with hi-intensity black
// foreground.
func HiBlackString(format string, a ...interface{}) string {
	return colourString(format, FgHiBlack, a...)
}

// HiRedString is a convenient helper function to return a string with hi-intensity red
// foreground.
func HiRedString(format string, a ...interface{}) string { return colourString(format, FgHiRed, a...) }

// HiGreenString is a convenient helper function to return a string with hi-intensity green
// foreground.
func HiGreenString(format string, a ...interface{}) string {
	return colourString(format, FgHiGreen, a...)
}

// HiYellowString is a convenient helper function to return a string with hi-intensity yellow
// foreground.
func HiYellowString(format string, a ...interface{}) string {
	return colourString(format, FgHiYellow, a...)
}

// HiBlueString is a convenient helper function to return a string with hi-intensity blue
// foreground.
func HiBlueString(format string, a ...interface{}) string { return colourString(format, FgHiBlue, a...) }

// HiMagentaString is a convenient helper function to return a string with hi-intensity magenta
// foreground.
func HiMagentaString(format string, a ...interface{}) string {
	return colourString(format, FgHiMagenta, a...)
}

// HiCyanString is a convenient helper function to return a string with hi-intensity cyan
// foreground.
func HiCyanString(format string, a ...interface{}) string { return colourString(format, FgHiCyan, a...) }

// HiWhiteString is a convenient helper function to return a string with hi-intensity white
// foreground.
func HiWhiteString(format string, a ...interface{}) string {
	return colourString(format, FgHiWhite, a...)
}
