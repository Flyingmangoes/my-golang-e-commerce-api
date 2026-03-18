package utils

import (
    "math/rand"
    "time"
)

func GenerateId() int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	nums := rng.Intn(900000000) + 100000000
	return nums
}