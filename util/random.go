package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomRetailer() string {
	return RandomString(6)
}

func CurrentDate() string {
	now := time.Now()
	return now.Format("2006-01-02")
}

func CurrentTime() string {
	now := time.Now()
	return now.Format("15:04")
}

func RandomDescription() string {
	return RandomString(12)
}

func RandomPrice() float64 {
	return rand.Float64()
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
