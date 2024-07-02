package pow

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GenerateChallenge(ctx context.Context) (string, error) {
	max := new(big.Int).Lsh(big.NewInt(1), uint(512))
	n, err := rand.Int(rand.Reader, max)

	if err != nil {
		return "", fmt.Errorf("error generating random number: %w", err)
	}

	return n.Text(10), nil
}

func (s *Service) ValidateChallenge(ctx context.Context, difficulty int, challenge string, response string) bool {
	hash := sha256.Sum256([]byte(challenge + response))
	hashString := hex.EncodeToString(hash[:])

	return hashString[:difficulty/4] == strings.Repeat("0", difficulty/4)
}
