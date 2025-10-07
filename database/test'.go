package main

import "fmt"

func main() {

	fmt.Println("Testing...")

	a := 50
	b := &a
	c := *b
	a = 100
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(a)

}
