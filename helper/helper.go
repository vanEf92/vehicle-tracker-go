package helper

import (
	"log"
	"os"
)

func GetEnvRequired(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Error: Environment variable '%s' is empty or missing", key)
	}
	return val
}