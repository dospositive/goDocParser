package main

import (
	"net/http"
	"sync"
	"time"
)

// AsyncNetCounter struct -
type AsyncNetCounter struct {
	substring         string
	countOfGoroutines int
	inURI             chan string

	outInfo chan URIInfo
}

// GetAsyncNetCounter - create a struct for work
func GetAsyncNetCounter(aSubstring string, aCountGoroutines int,
	aDataChan chan string) *AsyncNetCounter {
	a := AsyncNetCounter{aSubstring, aCountGoroutines, aDataChan, make(chan URIInfo)}
	return &a
}

// GetResultChan Getter - to outInfo chan URIInfo AsyncNetCounter
func (a *AsyncNetCounter) GetResultChan() <-chan URIInfo {
	return a.outInfo
}

func worker(countOf func(string) (int, error), in chan string, out chan URIInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	for val := range in {
		count, err := countOf(val)
		out <- URIInfo{val, count, err}
	}
}

func getNetCounter(targetSubstr string) func(string) (int, error) {
	return func(uri string) (int, error) {
		resp, err := http.Get(uri)

		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		var s StringCounter
		s.SetReader(resp.Body)
		s.SetSubstring(targetSubstr)
		return s.SafetyCount()
	}
}

// RunWorkers - run pool workers, they process inURI chan
func (a *AsyncNetCounter) RunWorkers() {
	go func() {
		var wg sync.WaitGroup
		netCounter := getNetCounter(a.substring)
		uriChan := make(chan string)
		countOfListener := 0

		for uri := range a.inURI {

			for isReadyToNext := false; !isReadyToNext; {
				select {
				case uriChan <- uri:
					isReadyToNext = true

				default:
					if countOfListener < a.countOfGoroutines {
						countOfListener++
						wg.Add(1)
						go worker(netCounter, uriChan, a.outInfo, &wg)
					} else {
						time.Sleep(5 * time.Millisecond)
					}
				}
			}
		}
		close(uriChan)
		wg.Wait()

		close(a.outInfo)
	}()
}
