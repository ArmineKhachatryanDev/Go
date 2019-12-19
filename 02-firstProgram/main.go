package main

import "fmt"

// "func main" is special.
//
// Go has to know where to start
//
// func main creates a starting point for Go
//
// After compiling the code,
// Go runtime will first run this function
func main() {
	// after: import "fmt"
	// Println function of "fmt" package becomes available

	// Look at what it looks like by typing in the console:
	//   go doc -src fmt Println

	// Println is just an exported function from
	//   "fmt" package

	// Exported = First Letter is uppercase
	fmt.Println("Hello Armine!")

	// Go cannot call Println function by itself.
	// That's why you need to call it here.
	// It only calls `func main` automatically.

	// -----

	// Go supports Unicode characters in string literals

}
