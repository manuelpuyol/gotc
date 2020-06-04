package hash

import (
	"crypto/sha256"
	"fmt"
)

func BTCHash(data []byte) [sha256.Size]byte {
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	return sha256.Sum256([]byte(hash))
}
