package tcpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"
)

type TCPHandler interface {
	Handle(ctx context.Context, conn net.Conn)
}

func StartServer(ctx context.Context, port int, handler TCPHandler) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	slog.InfoContext(ctx, "Server started", "port", port)

	defer func() {
		if err := listener.Close(); err != nil {
			slog.WarnContext(ctx, "Failed to close listener:", "error", err)
		}
	}()

	connPool := make(chan net.Conn)

	go func() {
		defer close(connPool)
		for {
			conn, err := listener.Accept()

			if ctx.Err() != nil {
				if err == nil {
					closeErr := conn.Close() // Don't leak connections

					if closeErr != nil {
						slog.WarnContext(ctx, "Error closing connection:", "error", closeErr)
					}
				}

				break
			}

			if err != nil {
				slog.WarnContext(ctx, "Error accepting connection:", "error", err)

				continue
			}

			connPool <- conn
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case conn, ok := <-connPool:
			if !ok {
				return ctx.Err()
			}

			slog.InfoContext(ctx, "Connection accepted", "remote", conn.RemoteAddr())
			go handler.Handle(ctx, conn)
		}
	}
}
