package quotes

import (
	"context"

	"github.com/quolpr/wisdom-pow/internal/service/quotes/model"
	"github.com/quolpr/wisdom-pow/internal/service/quotes/repo"
)

type Service struct {
	repo repo.QuoteRepo
}

func NewService(repo repo.QuoteRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetRandomQuote(ctx context.Context) model.Quote {
	return s.repo.GetRandomQuote(ctx)
}
