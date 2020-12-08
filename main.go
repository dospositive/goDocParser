package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	argK := flag.Int("k", 5, "count of goroutines")
	argQ := flag.String("q", "go", "string to found")
	flag.Parse()
	fmt.Println("k=", *argK, " q=", *argQ)

	stringChan := make(chan string)
	netCounter := GetAsyncNetCounter(*argQ, *argK, stringChan)
	urisAndCountChan := netCounter.GetResultChan()
	netCounter.RunWorkers()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)

		for end := scanner.Scan(); end; {
			stringChan <- scanner.Text()
			end = scanner.Scan()
		}
		close(stringChan)
	}()

	totalCount := 0
	for el := range urisAndCountChan {

		if el.err != nil {
			fmt.Println(el.err, " in", el.uri)
		} else {
			totalCount += el.count
			fmt.Println("Count for ", el.uri, ": ", el.count)
		}
	}

	fmt.Println("Total: ", totalCount)
}
