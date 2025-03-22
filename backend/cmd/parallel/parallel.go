package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/ericlagergren/decimal"
)

// enable parralelism by using GOMAXPROCS
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// r := mux.NewRouter()
	start := time.Now()
	// enable concurrency
	handleConcurrency()
	fmt.Printf("Total time: %v\n", time.Since(start))

}

func handleConcurrency() {
	var wg sync.WaitGroup
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go heavyHandler(&wg)
	}
	wg.Wait()
}

func heavyHandler(wg *sync.WaitGroup) {
	defer wg.Done()
	// 1 second
	time.Sleep(2 * time.Second)
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
