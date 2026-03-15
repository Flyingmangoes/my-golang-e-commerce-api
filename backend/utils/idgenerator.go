package utils

import (
    "fmt"
    "math/rand"
    "time"
)

func GenerateId() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	nums := rng.Intn(900000000) + 100000000
	return fmt.Sprintf("%d", nums)
}