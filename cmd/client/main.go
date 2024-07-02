package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type JSONQuote struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error connecting:", err)
		os.Exit(1)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Error closing connection:", err)
		}
	}()

	reader := bufio.NewReader(conn)
	difficultyStr, _ := reader.ReadString('\n')
	challenge, _ := reader.ReadString('\n')
	challenge = strings.TrimSpace(challenge)

	difficulty, err := strconv.Atoi(strings.TrimSpace(difficultyStr))
	if err != nil {
		log.Println("Error parsing difficulty:", err)

		os.Exit(1) //nolint:gocritic
	}

	for nonce := 0; ; nonce++ {
		test := fmt.Sprintf("%s%d", challenge, nonce)
		hash := sha256.Sum256([]byte(test))
		hashString := hex.EncodeToString(hash[:])
		if hashString[:difficulty/4] == strings.Repeat("0", difficulty/4) {
			_, err := fmt.Fprintf(conn, "%d\n", nonce)
			if err != nil {
				log.Println("Error sending response:", err)
				os.Exit(1)
			}
			break
		}
	}

	res, err := reader.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		log.Println("Error reading response:", err)
		os.Exit(1)
	}

	quote := JSONQuote{}
	err = json.Unmarshal(res, &quote)
	if err != nil {
		log.Println("Error unmarshalling quote:", err)
		os.Exit(1)
	}

	log.Printf("Quote: %s\n", quote.Quote)
	log.Printf("Author: %s\n", quote.Author)
}
