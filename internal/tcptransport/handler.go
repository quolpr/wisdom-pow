package tcptransport

import (
	"bufio"
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"strconv"
	"strings"

	"github.com/quolpr/wisdom-pow/internal/service/quotes/model"
)

type PowService interface {
	GenerateChallenge(ctx context.Context) (string, error)
	ValidateChallenge(ctx context.Context, difficulty int, challenge string, response string) bool
}

type QuoteService interface {
	GetRandomQuote(ctx context.Context) model.Quote
}

type Handler struct {
	powService   PowService
	quoteService QuoteService
	difficulty   int
}

func NewHandler(powService PowService, quoteService QuoteService, difficulty int) *Handler {
	return &Handler{powService: powService, quoteService: quoteService, difficulty: difficulty}
}

func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			slog.WarnContext(ctx, "Failed to close connection:", "error", err)
		}
	}()

	reader := bufio.NewReader(conn)

	n, err := h.powService.GenerateChallenge(ctx)

	if err != nil {
		slog.WarnContext(ctx, "Error generating random number", "error", err)

		return
	}

	slog.InfoContext(ctx, "Challenge generated", "challenge", n)

	_, err = conn.Write([]byte(strconv.Itoa(h.difficulty) + "\n" + n + "\n"))

	if err != nil {
		slog.WarnContext(ctx, "Error writing to connection", "error", err)
	}

	response, err := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	if !h.powService.ValidateChallenge(ctx, h.difficulty, n, response) {
		slog.WarnContext(ctx, "Challenge validation failed")
		return
	}

	slog.InfoContext(ctx, "Challenge validated")

	if err != nil {
		slog.WarnContext(ctx, "Error reading response", "error", err)

		return
	}

	res, err := json.Marshal(h.quoteService.GetRandomQuote(ctx))

	if err != nil {
		slog.WarnContext(ctx, "Error marshalling quote", "error", err)

		return
	}

	_, err = conn.Write(res)

	if err != nil {
		slog.WarnContext(ctx, "Error writing quote", "error", err)

		return
	}

	slog.InfoContext(ctx, "Quote sent")
}
