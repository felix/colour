/*
Package colour is an ANSI colour package to output colourized or SGR defined
output to the standard output. The API can be used in several way, pick one
that suits you.

Use simple and default helper functions with predefined foreground colours:

    colour.Cyan("Prints text in cyan.")

    // a newline will be appended automatically
    colour.Blue("Prints %s in blue.", "text")

    // More default foreground colours..
    colour.Red("We have red")
    colour.Yellow("Yellow colour too!")
    colour.Magenta("And many others ..")

    // Hi-intensity colours
    colour.HiGreen("Bright green colour.")
    colour.HiBlack("Bright black means gray..")
    colour.HiWhite("Shiny white colour!")

However there are times where custom colour mixes are required. Below are some
examples to create custom colour objects and use the print functions of each
separate colour object.

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
    whiteBackground.Println("Red text with White background.")

    // Use your own io.Writer output
    colour.New(colour.FgBlue).Fprintln(myWriter, "blue colour!")

    blue := colour.New(colour.FgBlue)
    blue.Fprint(myWriter, "This will print text in blue.")

You can create PrintXxx functions to simplify even more:

    // Create a custom print function for convenient
    red := colour.New(colour.FgRed).PrintfFunc()
    red("warning")
    red("error: %s", err)

    // Mix up multiple attributes
    notice := colour.New(colour.Bold, colour.FgGreen).PrintlnFunc()
    notice("don't forget this...")

You can also FprintXxx functions to pass your own io.Writer:

    blue := colour.New(FgBlue).FprintfFunc()
    blue(myWriter, "important notice: %s", stars)

    // Mix up with multiple attributes
    success := colour.New(colour.Bold, colour.FgGreen).FprintlnFunc()
    success(myWriter, don't forget this...")


Or create SprintXxx functions to mix strings with other non-colourized strings:

    yellow := New(FgYellow).SprintFunc()
    red := New(FgRed).SprintFunc()

    fmt.Printf("this is a %s and this is %s.\n", yellow("warning"), red("error"))

    info := New(FgWhite, BgGreen).SprintFunc()
    fmt.Printf("this %s rocks!\n", info("package"))

Windows support is enabled by default. All Print functions work as intended.
However only for colour.SprintXXX functions, user should use fmt.FprintXXX and
set the output to colour.Output:

    fmt.Fprintf(colour.Output, "Windows support: %s", colour.GreenString("PASS"))

    info := New(FgWhite, BgGreen).SprintFunc()
    fmt.Fprintf(colour.Output, "this %s rocks!\n", info("package"))

Using with existing code is possible. Just use the Set() method to set the
standard output to the given parameters. That way a rewrite of an existing
code is not required.

    // Use handy standard colours.
    colour.Set(colour.FgYellow)

    fmt.Println("Existing text will be now in Yellow")
    fmt.Printf("This one %s\n", "too")

    colour.Unset() // don't forget to unset

    // You can mix up parameters
    colour.Set(colour.FgMagenta, colour.Bold)
    defer colour.Unset() // use it in your function

    fmt.Println("All text will be now bold magenta.")

There might be a case where you want to disable colour output (for example to
pipe the standard output of your app to somewhere else). `Colour` has support to
disable colours both globally and for single colour definition. For example
suppose you have a CLI app and a `--no-colour` bool flag. You can easily disable
the colour output with:

    var flagNoColour = flag.Bool("no-colour", false, "Disable colour output")

    if *flagNoColour {
    	colour.NoColour = true // disables colourized output
    }

It also has support for single colour definitions (local). You can
disable/enable colour output on the fly:

     c := colour.New(colour.FgCyan)
     c.Println("Prints cyan text")

     c.DisableColour()
     c.Println("This is printed without any colour")

     c.EnableColour()
     c.Println("This prints again cyan...")
*/
package colour
