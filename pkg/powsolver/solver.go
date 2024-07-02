package powsolver

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func FindSolution(challenge string, difficulty int) int {
	for nonce := 0; ; nonce++ {
		test := fmt.Sprintf("%s%d", challenge, nonce)
		hash := sha256.Sum256([]byte(test))
		hashString := hex.EncodeToString(hash[:])
		if hashString[:difficulty/4] == strings.Repeat("0", difficulty/4) {
			return nonce
		}
	}
}
