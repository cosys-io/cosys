package main

import "fmt"

func main() {
	dataChan := make(chan int, 1)
	dataChan <- 123

	n := <-dataChan
	fmt.Printf("%d\n", n)
}
