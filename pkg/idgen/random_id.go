package idgen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateReadableID generates a unique, readable ID with a prefix.
// Example: ORD-20251031191045-4fa2b3
func GenerateReadableID(prefix string) string {
	timestamp := time.Now().Format("20060102150405") // YYYYMMDDHHMMSS
	randomBytes := make([]byte, 3)                   // 3 bytes = 6 hex chars
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err) // Or handle gracefully if you have a logger
	}
	randomPart := hex.EncodeToString(randomBytes)

	if prefix != "" {
		return fmt.Sprintf("%s-%s-%s", prefix, timestamp, randomPart)
	}
	return fmt.Sprintf("%s-%s", timestamp, randomPart)
}
