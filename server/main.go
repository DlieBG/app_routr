package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("app-routr")
	fmt.Println()

	waiter := sync.WaitGroup{}
	config := loadConfig()

	scheduleProd(config, &waiter)
	scheduleDev(config, &waiter)

	waiter.Wait()
}
