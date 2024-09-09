package random

import (
	"math/rand/v2"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

var retry = 0

const retrylimit = 3

func NewRandomString(length int) string {
	randStr := make([]rune, length)
	for i := 0; i < length; i++ {
		switch rand.IntN(3) {
		case 0:
			randStr[i] = rune(int('a') + rand.IntN(int('z')-int('a')))
		case 1:
			randStr[i] = rune(int('A') + rand.IntN(int('Z')-int('A')))
		case 2:
			randStr[i] = rune(int('1') + rand.IntN(int('9')-int('0')))
		}
	}
	return string(randStr)
}
