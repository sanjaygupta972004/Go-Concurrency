package main

import (
	"fmt"
	"time"
)

func generateNumber(num []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range num {
			fmt.Println("Generating numbers", num)
			out <- n
			time.Sleep(1 * time.Second)
		}
		close(out)
	}()

	return out
}

func squareNumber(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {

		for num := range in {
			fmt.Println("Receiving number from generater go routine", num)
			sq := num * num
			time.Sleep(1 * time.Second)
			out <- sq
		}
		close(out)
	}()

	return out
}

func totalSquare(in <-chan int) <-chan int {
	res := make(chan int)
	var sum int

	go func() {
		for sqNum := range in {
			fmt.Println("Adding total numbers", sqNum)
			sum = sqNum + sum
			time.Sleep(1 * time.Second)
		}
		res <- sum
		close(res)
	}()

	return res
}

func PipelineExa() {
	fmt.Println("Solving pipeline example")

	nums := []int{1, 2, 3, 5, 6, 7, 7, 8, 9}

	generatdNumber := generateNumber(nums)
	squareRes := squareNumber(generatdNumber)
	res := totalSquare(squareRes)

	fmt.Printf("Total sum of generated numbers :%d \n", <-res)
}
