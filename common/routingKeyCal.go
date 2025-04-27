package common

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
)

func RoutingKeyCalculator(key string) string {
	if len(key) != 64 {
		log.Println("Ivalid Validator Key")
		return ""
	}

	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
