package tcptransport_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/quolpr/wisdom-pow/internal/service/pow"
	"github.com/quolpr/wisdom-pow/internal/service/quotes"
	"github.com/quolpr/wisdom-pow/internal/service/quotes/model"
	"github.com/quolpr/wisdom-pow/internal/service/quotes/repo/jsonquote"
	. "github.com/quolpr/wisdom-pow/internal/tcptransport"
	"github.com/quolpr/wisdom-pow/pkg/powsolver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	powService := pow.NewService()
	repo, err := jsonquote.NewRepo()
	require.NoError(t, err)
	quoteService := quotes.NewService(repo)

	addr := ""
	for i := 0; i < 10; i++ {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err == nil {
			addr = ln.Addr().String()
			err := ln.Close()
			if err != nil {
				t.Log("Error closing listener:", err)
			}
			break
		}
	}
	require.NotEqual(t, "", addr, "could not find available port after 10 attempts")

	handler := NewHandler(powService, quoteService, 1)

	// Create a listener on an address
	ln, err := net.Listen("tcp", addr)
	require.NoError(t, err)
	defer func() {
		if err := ln.Close(); err != nil {
			t.Log("Error closing listener:", err)
		}
	}()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go handler.Handle(ctx, conn)
		}
	}()

	time.Sleep(1 * time.Second) // Ensure listener startup before client dials it.

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	defer func() {
		if err := conn.Close(); err != nil {
			t.Log("Error closing connection:", err)
		}
	}()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	// Read the difficulty and challenge from the handler:
	readLine := func() string {
		line, err := reader.ReadString('\n')
		require.NoError(t, err)
		return strings.TrimSpace(line)
	}

	// First read expected difficulty and challenge
	difficulty := readLine()
	challenge := readLine()

	assert.Equal(t, "1", difficulty, "Expected difficulty to be 1")

	solution := powsolver.FindSolution(challenge, 1)

	_, err = fmt.Fprintf(conn, "%d\n", solution)
	require.NoError(t, err)
	err = writer.Flush()

	assert.NoError(t, err)

	// Read the quote JSON
	quoteJSON, err := reader.ReadBytes('\n')
	require.ErrorIs(t, err, io.EOF)

	// Parse and check quote
	var quote model.Quote
	err = json.Unmarshal(quoteJSON, &quote)
	require.NoError(t, err)
}
