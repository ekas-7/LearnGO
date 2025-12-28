package main

import "fmt"

func main() {
	fmt.Println("=== Learning Variables in Go ===")

	// 1. Variable Declaration with var keyword
	fmt.Println("1. Variable Declaration with 'var':")
	var name string = "Alice"
	var age int = 25
	fmt.Printf("   Name: %s, Age: %d\n\n", name, age)

	// 2. Type Inference (Go infers the type)
	fmt.Println("2. Type Inference:")
	var city = "New York" // Go infers this is a string
	var temperature = 72  // Go infers this is an int
	fmt.Printf("   City: %s, Temperature: %d°F\n\n", city, temperature)

	// 3. Short Declaration (most common in Go)
	fmt.Println("3. Short Declaration with ':=':")
	country := "USA"          // string
	population := 331_000_000 // int (underscores for readability)
	isLarge := true           // bool
	fmt.Printf("   Country: %s, Population: %d, Large: %t\n\n", country, population, isLarge)

	// 4. Multiple Variable Declaration
	fmt.Println("4. Multiple Variable Declaration:")
	var x, y, z int = 1, 2, 3
	a, b, c := "Go", "is", "awesome"
	fmt.Printf("   Numbers: %d, %d, %d\n", x, y, z)
	fmt.Printf("   Words: %s %s %s\n\n", a, b, c)

	// 5. Zero Values (default values when not initialized)
	fmt.Println("5. Zero Values:")
	var defaultInt int
	var defaultFloat float64
	var defaultString string
	var defaultBool bool
	fmt.Printf("   int: %d, float64: %f, string: '%s', bool: %t\n\n",
		defaultInt, defaultFloat, defaultString, defaultBool)

	// 6. Constants (values that cannot be changed)
	fmt.Println("6. Constants:")
	const pi = 3.14159
	const greeting = "Hello, Go!"
	fmt.Printf("   Pi: %f\n", pi)
	fmt.Printf("   Greeting: %s\n\n", greeting)

	// 7. Different Data Types
	fmt.Println("7. Common Data Types:")
	var integer int = 42
	var floatingPoint float64 = 3.14
	var text string = "Golang"
	var flag bool = true
	var smallNumber int8 = 127        // -128 to 127
	var bigNumber uint64 = 1844674407 // unsigned (positive only)
	fmt.Printf("   int: %d\n", integer)
	fmt.Printf("   float64: %f\n", floatingPoint)
	fmt.Printf("   string: %s\n", text)
	fmt.Printf("   bool: %t\n", flag)
	fmt.Printf("   int8: %d\n", smallNumber)
	fmt.Printf("   uint64: %d\n\n", bigNumber)

	// 8. Type Conversion
	fmt.Println("8. Type Conversion:")
	var num int = 100
	var decimal float64 = float64(num) // Explicit conversion
	fmt.Printf("   Integer: %d converted to Float: %f\n\n", num, decimal)

	// 9. Block Declaration
	fmt.Println("9. Block Declaration:")
	var (
		username = "gopher"
		email    = "gopher@golang.org"
		active   = true
	)
	fmt.Printf("   User: %s, Email: %s, Active: %t\n\n", username, email, active)

	// 10. Important Notes
	fmt.Println("=== Key Points ===")
	fmt.Println("• Use := for short declaration (most common)")
	fmt.Println("• Use var when you need to declare without initialization")
	fmt.Println("• Use const for values that never change")
	fmt.Println("• Variable names should be camelCase (e.g., firstName)")
	fmt.Println("• Exported variables start with uppercase (e.g., Name)")
	fmt.Println("• Go is statically typed - types are checked at compile time")
}
