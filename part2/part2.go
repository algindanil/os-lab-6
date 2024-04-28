package main

import (
	"fmt"
	"sync"
)

const iterCount = int(10e8)

var counter int
var mutex sync.Mutex

func SafeIncrement() {
	mutex.Lock()

	//fmt.Println(counter)
	defer mutex.Unlock()
	counter++
}

func UnsafeIncrement() {
	//fmt.Println(counter)
	counter++
}

func IncrementJob(wg *sync.WaitGroup, IncFunc func()) {
	defer wg.Done()
	for range iterCount {
		IncFunc()
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go IncrementJob(&wg, SafeIncrement)
	go IncrementJob(&wg, SafeIncrement)

	wg.Wait()

	fmt.Printf("Counter value with mutexes: %d\n", counter)

	counter = 0
	wg = sync.WaitGroup{}
	wg.Add(2)

	go IncrementJob(&wg, UnsafeIncrement)
	go IncrementJob(&wg, UnsafeIncrement)

	wg.Wait()

	fmt.Printf("Counter value without mutexes: %d", counter)
}
