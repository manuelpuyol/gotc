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

func Spinner(message string) {
	for {
		for _, r := range `⣽⣾⣷⣯⣟⡿⢿⣻` {
			fmt.Printf("\r %s %c ", message, r)
			time.Sleep(constants.SpinnerDelay * time.Millisecond)
		}
	}
}
