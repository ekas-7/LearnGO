package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
	result := myFunction("GoLang")
	fmt.Println("Back in main")
	fmt.Println(result)
}

func myFunction(p1 string) string {
	fmt.Println("Inside myFunction")
	fmt.Println(p1)
	return p1 + p1
}