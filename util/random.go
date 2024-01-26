package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)


	for i:= 0; i < n; i++ {
		c := alphabet[RandomInt(0, int64(k-1))]
		sb.WriteByte(c)
	}

	return sb.String();
}

func RandomDate() time.Time {
	min := time.Date(2020, 1, 0, 0, 0, 0, 0, time.Local).Unix()
	max := time.Date(2020, 12, 0, 0, 0, 0, 0, time.Local).Unix()

	sec := RandomInt(min, max)
	return time.Unix(sec, 0)
}

func RandomReminderDate(dueDate time.Time) time.Time {
	min := dueDate.AddDate(0, 0, -1).Unix()
	max := dueDate.Unix()

	sec := RandomInt(min, max)
	return time.Unix(sec, 0)
}

func RandomEmail() string {
	return RandomString(6) + "@gmail.com"
}