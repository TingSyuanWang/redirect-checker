package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var redirects int

func main() {
	// Set the concurrency limit
	concurrencyLimit := 50

	// Initialize the number of redirects to 0
	redirects = 0

	// Open the CSV file
	file, err := os.Open("urls.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the lines from the CSV file
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// Open the log file for writing
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Set the log output to the file
	log.SetOutput(logFile)

	tasks := make(chan string, concurrencyLimit)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < concurrencyLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				checkHasRedirect(task)
			}
		}()
	}

	for key, line := range lines {
		fmt.Println("line" + strconv.Itoa(key))

		if len(line) > 0 {
			url := line[0]
			tasks <- url
		}
	}

	close(tasks)

	wg.Wait()

	log.Println("Total redirects:", redirects)
}

func checkHasRedirect(url string) {
	// Send an HTTP GET request to the API endpoint
	client := http.Client{
		CheckRedirect: redirectHandler,
	}

	response, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error making the request: %s\n", err)
		return
	}
	defer response.Body.Close()
}

func redirectHandler(req *http.Request, via []*http.Request) error {
	log.Printf("===== Redirect %d =====\n", redirects+1)
	log.Printf("Redirected from: %s\n", via[len(via)-1].URL)
	log.Printf("Redirected to: %s\n", req.URL)
	redirects++
	return nil
}
