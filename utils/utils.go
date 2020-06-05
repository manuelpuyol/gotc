package utils

import (
	"math/rand"
)

const ASCIIStart = 65
const ASCIIEnd = 90
const Length = 10

func RandomBytes() []byte {
	bytes := make([]byte, Length)
	for i := 0; i < Length; i++ {
		bytes[i] = byte(randomInt(ASCIIStart, ASCIIEnd))
	}
	return bytes
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
