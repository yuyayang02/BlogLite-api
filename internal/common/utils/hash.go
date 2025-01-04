package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func ShortHash(inputs ...string) string {
	hasher := sha256.New()
	for _, s := range inputs {
		hasher.Write([]byte(s))
	}
	return hex.EncodeToString(hasher.Sum(nil))[:8]
}
