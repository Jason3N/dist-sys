package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/ericlagergren/decimal"
)

var (
	global = 0
	mu     sync.Mutex
	wg     sync.WaitGroup
)

func main() {
	start := time.Now()
	handleConcurrency()
	fmt.Printf("Total time: %v\n", time.Since(start))
	fmt.Printf("%v\n", global)
}

func handleConcurrency() {
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go heavyHandler(&wg)
	}

	wg.Wait()

}

func heavyHandler(wg *sync.WaitGroup) {
	defer wg.Done()
	// some action must be done here, simulated with accessing atomic var
	// this will be placeholder for read/write on db
	// will stop race conditions by placing mutex lock
	mu.Lock()
	global++
	mu.Unlock()
	// assume action takes 1 second
	time.Sleep(1 * time.Second)
}

func mediumHandler(wg *sync.WaitGroup) {
	defer wg.Done()
	// 0.5 second
	mediumTime := new(decimal.Big).SetFloat64(0.5)
	mediumTimeFloat, _ := mediumTime.Float64()
	sleepDuration := time.Duration(mediumTimeFloat * float64(time.Second))
	time.Sleep(sleepDuration)
}

func lowHandler(wg *sync.WaitGroup) {
	defer wg.Done()
	// 0.25 second
	lowTime := new(decimal.Big).SetFloat64(0.5)
	lowFloat, _ := lowTime.Float64()
	sleepDuration := time.Duration(lowFloat * float64(time.Second))
	time.Sleep(sleepDuration)
}
