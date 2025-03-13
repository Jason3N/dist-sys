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
	// router
	// r := mux.NewRouter()
	start := time.Now()
	// enable concurrency
	handleConcurrency()
	fmt.Printf("Total time: %v\n", time.Since(start))
	// have low-level, medium-level, high-level
	// r.HandleFunc("/highLevel", heavyHandler)
	// r.HandleFunc("/mediumLevel", mediumHandler)
	// r.HandleFunc("/lowLevel", lowHandler)

	// err := http.ListenAndServe(":8080", r)
	// if err != nil {
	// 	fmt.Printf("error running server")
	// 	return
	// }
}

// jesus christ?? it can run 10,000 groups of
// mutex
func handleConcurrency() {
	var wg sync.WaitGroup
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go heavyHandler(&wg)
	}
	wg.Wait()
}

func heavyHandler(wg *sync.WaitGroup) {
	defer wg.Done()
	// 1 second
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
