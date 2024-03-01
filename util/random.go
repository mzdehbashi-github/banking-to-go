package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghigklmnopqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandString generates a random string of length n
func RandString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwnerName generates a random Account.Owner
func RandomOwner() string {
	return RandString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random Account.Currency
func RandomCurrency() string {
	currencies := []string{"EUR", "USD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
