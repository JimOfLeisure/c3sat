package main

import "C"
import "fmt"

//export ReadFile
func ReadFile() {
	fmt.Println("Stub test response during development. Will eventually accept a path and return a byte array.")
}

func main() {
	fmt.Print("This code was meant to be compiled as a shared libary, not an executable.")
}
