package repo

import (
	"context"

	"github.com/quolpr/wisdom-pow/internal/service/quotes/model"
)

type QuoteRepo interface {
	GetRandomQuote(ctx context.Context) model.Quote
}
