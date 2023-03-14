package main

import "fmt"

func main() {
	fmt.Println(Add(1, 2))
	fmt.Println(Add("1", "2"))
}

func Add[T int | string](a, b T) T {
	return a + b
}
