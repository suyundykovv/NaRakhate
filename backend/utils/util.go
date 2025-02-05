package utils

import (
	"log"
	"math/rand"
	"time"
)

func CatchCriticalPoint() {
	if r := recover(); r != nil {
		log.Printf("Recovered from error")
	}
}

func contains(validStatuses []string, status string) bool {
	for _, validStatus := range validStatuses {
		if validStatus == status {
			return true
		}
	}
	return false
}

// RandomInt generates a random integer between min and max (inclusive).
func RandomInt(min, max int) int {
	// Initialize the random number generator with the current time as the seed
	rand.Seed(time.Now().UnixNano())

	// Return a random integer between min and max (inclusive)
	return rand.Intn(max-min+1) + min
}
