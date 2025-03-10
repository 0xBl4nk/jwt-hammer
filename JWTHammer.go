package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"sync"
	"bufio"
)

// Split JWT into header, payload, and signature parts
func splitJWT(token string) (string, string, string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid JWT format")
	}
	return parts[0], parts[1], parts[2], nil
}

// Compute HMAC-SHA256 signature
func computeSignature(header, payload, secret string) string {
	data := header + "." + payload
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	signature := h.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(signature)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run jwt_brute.go <jwt_token> <wordlist_file>")
		os.Exit(1)
	}

	jwtToken := os.Args[1]
	wordlistFile := os.Args[2]

	header, payload, signature, err := splitJWT(jwtToken)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Open wordlist file
	file, err := os.Open(wordlistFile)
	if err != nil {
		fmt.Printf("Error opening wordlist file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// For larger files, you may need to increase the buffer size
	const maxCapacity = 512 * 1024 // 512KB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	// Configure workers for parallel processing
	numWorkers := 4
	var wg sync.WaitGroup
	passwords := make(chan string)
	found := make(chan string)
	done := make(chan struct{})

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for password := range passwords {
				calculatedSignature := computeSignature(header, payload, password)
				if calculatedSignature == signature {
					found <- password
					return
				}
			}
		}()
	}

	// Start a goroutine to signal when all workers are done
	go func() {
		wg.Wait()
		close(done)
	}()

	// Start a goroutine to read passwords from the file
	go func() {
		defer close(passwords)
		for scanner.Scan() {
			password := scanner.Text()
			select {
			case passwords <- password:
			case <-done:
				return
			}
		}
	}()

	// Wait for a match or completion
	select {
	case password := <-found:
		fmt.Printf("Success! Secret key found: %s\n", password)
	case <-done:
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading wordlist: %v\n", err)
		} else {
			fmt.Println("No matching secret key found in wordlist")
		}
	}
}
