package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"sync"
	"sync/atomic"
	"time"
)

var number atomic.Int64

func main() {
	valid := true

	if valid {
		f, err := os.Create("./pprof/profile.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		// Start CPU profiling
		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()

		// Start tracing
		traceFile, err := os.Create("./pprof/trace.out")
		if err != nil {
			panic(err)
		}
		defer traceFile.Close()

		if err := trace.Start(traceFile); err != nil {
			panic(err)
		}
		defer trace.Stop()
	}

	startFn := time.Now()
	countItems := 1000
	messagesInput := make(chan struct{}, 100)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()

		internalWG := &sync.WaitGroup{}
		sem := make(chan struct{}, 100)

		for range messagesInput {
			sem <- struct{}{}
			internalWG.Add(1)

			go func() {

				defer func() {
					internalWG.Done()
					<-sem

				}()

				CountMessages(int64(1))
			}()
		}
		internalWG.Wait()

	}()

	for i := 0; i < countItems; i++ {
		messagesInput <- struct{}{}
	}

	close(messagesInput)

	wg.Wait()
	fmt.Println("time total duration: ", time.Since(startFn))
	fmt.Println("Closing channel...")

}

func CountMessages(numberToIncrement int64) {
	time.Sleep(time.Second * 3)
	number.Add(numberToIncrement)

	fmt.Println(fmt.Sprintf("runtime number: %d", number.Load()))

}
//
