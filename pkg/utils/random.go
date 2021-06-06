package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())

}

func RandInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)

}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabets)

	for i := 0; i < n; i++ {
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomKey() string {
	return RandomString(int(RandInt(5, 15)))
}

func RandomValue() string {
	return "v" + RandomString(int(RandInt(5, 15)))
}

func RandomPinCode() int64 {
	return RandInt(700000, 800000)
}

func RandomUserName() string {
	return RandomString(6)
}
