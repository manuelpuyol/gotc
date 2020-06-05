package utils

import (
	"fmt"
	"gotc/constants"
	"math/rand"
	"time"
)

func RandomBytes() []byte {
	bytes := make([]byte, constants.Length)
	for i := 0; i < constants.Length; i++ {
		bytes[i] = byte(randomInt(constants.ASCIIStart, constants.ASCIIEnd))
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
