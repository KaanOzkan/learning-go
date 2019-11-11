package main

import "fmt"

func main() {
	fmt.Println("go" + "lang")

	fmt.Println("1+1 = ", 1+1)

	// Type is determined at compile time?? Interesting
	interpolated := fmt.Sprintf("2+2 = %d, hello:%s", 4, "world")
	fmt.Println(interpolated)

	fmt.Println(!true)
}
