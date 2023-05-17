package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func searchFiles(folderPath, lookupValue string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			lineNumber := 0
			for scanner.Scan() {
				lineNumber++
				line := scanner.Text()
				if strings.Contains(line, lookupValue) {
					results <- fmt.Sprintf("found in db: %s, line: %d\n%s", path, lineNumber, line)
				}
			}

			if err := scanner.Err(); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		results <- fmt.Sprintf("err: %s", err)
	}
}

func main() {
  fmt.Print("(example: test@gmail.com) enter the lookup value: ")
	var lookupValue string
	fmt.Scanln(&lookupValue)

  fmt.Print("(example: ./main) enter the folder path: ")
	var folderPath string
	fmt.Scanln(&folderPath)

	var wg sync.WaitGroup

	results := make(chan string)

	wg.Add(1)
	go searchFiles(folderPath, lookupValue, results, &wg)

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}
