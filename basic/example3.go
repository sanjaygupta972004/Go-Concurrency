package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// count all phrase from the files and also count total number of files

func countOccurrance(filepath, phase string, resChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	content, err := os.ReadFile(filepath)
	if err != nil {
		resChan <- 0
	}

	count := strings.Count(strings.ToLower(string(content)), strings.ToLower(phase))

	resChan <- count
}

func processInOccurance(dir, phrase string) int {

	resultChan := make(chan int)
	var wg sync.WaitGroup

	totalCount := 0

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error while reading directory:", err)
		return 0
	}

	for _, file := range files {
		wg.Add(1)
		go countOccurrance(dir+"/"+file.Name(), phrase, resultChan, &wg)

	}
	go func() {
		wg.Wait()
		close(resultChan)

	}()

	for count := range resultChan {
		totalCount += count
	}

	return totalCount
}

func MainExa3() {

	dir := "RandomFile"

	phrase := "Hello Go"

	totalPhraseOccurrance := processInOccurance(dir, phrase)
	fmt.Printf("Total occurrance is present %d", totalPhraseOccurrance)
}
