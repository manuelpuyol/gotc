package utils

import (
	"fmt"
	"math/rand"
	"time"
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

func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}

func Spinner(delay time.Duration) {
	for {
		for _, r := range `⣽⣾⣷⣯⣟⡿⢿⣻` {
			fmt.Printf("\r Mining... %c ", r)
			time.Sleep(delay)
		}
	}
}
