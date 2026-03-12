package utils

import (
    "fmt"
    "math/rand"
    "time"
)

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateId() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	nums1 := rng.Intn(9000) + 1000

	alpha := make([]byte, 3)
	for i := range alpha {
		alpha[i] = letters[rng.Intn(len(letters))]
	}

	nums2 := rng.Intn(9000) + 1000

	return fmt.Sprintf("%d-%s-%d", nums1, alpha, nums2)
}