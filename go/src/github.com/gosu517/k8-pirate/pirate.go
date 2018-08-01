package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello!")
	Sum(5, 5)
}

func Sum(x int, y int) int {
    return x + y
}

