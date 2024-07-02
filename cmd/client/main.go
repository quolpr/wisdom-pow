package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/quolpr/wisdom-pow/pkg/powsolver"
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

	solution := powsolver.FindSolution(challenge, difficulty)

	_, err = fmt.Fprintf(conn, "%d\n", solution)
	if err != nil {
		log.Println("Error sending response:", err)
		os.Exit(1)
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
