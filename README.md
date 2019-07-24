# Colour

Colour lets you use colourized outputs in terms of [ANSI Escape
Codes](http://en.wikipedia.org/wiki/ANSI_escape_code#Colours) in Go (Golang). It
has support for Windows too! The API can be used in several ways, pick one that
suits you.

## Install

```bash
go get github.com/felix/colour
```

## Examples

### Standard colours

```go
// Print with default helper functions
colour.Cyan("Prints text in cyan.")

// A newline will be appended automatically
colour.Blue("Prints %s in blue.", "text")

// These are using the default foreground colours
colour.Red("We have red")
colour.Magenta("And many others ..")

```

### Mix and reuse colours

```go
// Create a new colour object
c := colour.New(colour.FgCyan).Add(colour.Underline)
c.Println("Prints cyan text with an underline.")

// Or just add them to New()
d := colour.New(colour.FgCyan, colour.Bold)
d.Printf("This prints bold cyan %s\n", "too!.")

// Mix up foreground and background colours, create new mixes!
red := colour.New(colour.FgRed)

boldRed := red.Add(colour.Bold)
boldRed.Println("This will print text in bold red.")

whiteBackground := red.Add(colour.BgWhite)
whiteBackground.Println("Red text with white background.")
```

### Use your own output (io.Writer)

```go
// Use your own io.Writer output
colour.New(colour.FgBlue).Fprintln(myWriter, "blue colour!")

blue := colour.New(colour.FgBlue)
blue.Fprint(writer, "This will print text in blue.")
```

### Custom print functions (PrintFunc)

```go
// Create a custom print function for convenience
red := colour.New(colour.FgRed).PrintfFunc()
red("Warning")
red("Error: %s", err)

// Mix up multiple attributes
notice := colour.New(colour.Bold, colour.FgGreen).PrintlnFunc()
notice("Don't forget this...")
```

### Custom fprint functions (FprintFunc)

```go
blue := colour.New(FgBlue).FprintfFunc()
blue(myWriter, "important notice: %s", stars)

// Mix up with multiple attributes
success := colour.New(colour.Bold, colour.FgGreen).FprintlnFunc()
success(myWriter, "Don't forget this...")
```

### Insert into noncolour strings (SprintFunc)

```go
// Create SprintXxx functions to mix strings with other non-colourized strings:
yellow := colour.New(colour.FgYellow).SprintFunc()
red := colour.New(colour.FgRed).SprintFunc()
fmt.Printf("This is a %s and this is %s.\n", yellow("warning"), red("error"))

info := colour.New(colour.FgWhite, colour.BgGreen).SprintFunc()
fmt.Printf("This %s rocks!\n", info("package"))

// Use helper functions
fmt.Println("This", colour.RedString("warning"), "should be not neglected.")
fmt.Printf("%v %v\n", colour.GreenString("Info:"), "an important message.")

// Windows supported too! Just don't forget to change the output to colour.Output
fmt.Fprintf(colour.Output, "Windows support: %s", colour.GreenString("PASS"))
```

### Plug into existing code

```go
// Use handy standard colours
colour.Set(colour.FgYellow)

fmt.Println("Existing text will now be in yellow")
fmt.Printf("This one %s\n", "too")

colour.Unset() // Don't forget to unset

// You can mix up parameters
colour.Set(colour.FgMagenta, colour.Bold)
defer colour.Unset() // Use it in your function

fmt.Println("All text will now be bold magenta.")
```

### Disable/Enable colour

There might be a case where you want to explicitly disable/enable colour output.
the `go-isatty` package will automatically disable colour output for non-tty
output streams (for example if the output were piped directly to `less`)

`Colour` has support to disable/enable colours both globally and for single colour
definitions. For example suppose you have a CLI app and a `--no-colour` bool
flag. You can easily disable the colour output with:

```go

var flagNoColour = flag.Bool("no-colour", false, "Disable colour output")

if *flagNoColour {
	colour.NoColour = true // disables colourized output
}
```

It also has support for single colour definitions (local). You can
disable/enable colour output on the fly:

```go
c := colour.New(colour.FgCyan)
c.Println("Prints cyan text")

c.DisableColour()
c.Println("This is printed without any colour")

c.EnableColour()
c.Println("This prints again cyan...")
```

## License

The MIT License (MIT).

