package main

import "fmt"

type Number interface {
	int8 | int16 | int32 | int64 | float32 | float64
}

func newGenericFunc[age Number](a age) {
	val := float32(a) + 1
	fmt.Println("New generic func: ", val)
}

func BubbleSort[T Number](input []T) []T {
	n := len(input)
	swapped := true
	for swapped {
		swapped = false
		for i := 0; i < n-1; i++ {
			if input[i] > input[i+1] {
				input[i+1], input[i] = input[i], input[i+1]
				swapped = true
			}
		}
	}
	return input
}

func main() {
	fmt.Println("Go generics tutorial")

	var testAge int64 = 28
	var testAge2 float32 = 28.6

	newGenericFunc(testAge)
	newGenericFunc(testAge2)

	sorted := BubbleSort([]int32{111, 22, 333, 44, 555})
	fmt.Println("Sorted: ", sorted)
}
