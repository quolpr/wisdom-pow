package jsonquote

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/quolpr/wisdom-pow/internal/service/quotes/model"
)

type JSONQuote struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

//go:embed quotes.json
var s []byte

type Repo struct {
	quotes []model.Quote
}

func NewRepo() (*Repo, error) {
	var jsonQuotes []JSONQuote
	err := json.Unmarshal(s, &jsonQuotes)

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling quotes json: %w", err)
	}

	modelQuotes := make([]model.Quote, len(jsonQuotes))
	for i, j := range jsonQuotes {
		modelQuotes[i] = model.Quote{
			Quote:  j.Quote,
			Author: j.Author,
		}
	}

	return &Repo{
		quotes: modelQuotes,
	}, nil
}

func (r *Repo) GetRandomQuote(ctx context.Context) model.Quote {
	return r.quotes[rand.Intn(len(r.quotes))] //nolint:gosec
}
