package main

import (
	"fmt"
)

// HelloWorld is a simple function that returns a greeting
func HelloWorld(name string) string {
	// Return a personalized greeting
	if name == "" {
		return "Hello, World!"
	}
	return fmt.Sprintf("Hello, %s!", name)
}

func main() {
	// Print a greeting to the console
	message := HelloWorld("Gopher")
	fmt.Println(message)

	// Test with different inputs
	fmt.Println(HelloWorld(""))

	// Example of using a multiline string
	multiline := `This is a
multiline string
in Go!`
	fmt.Println(multiline)
}
