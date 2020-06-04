package header

import (
	"crypto/sha256"
)

type Header struct {
	Nonce uint              `json:"nonce"`
	Prev  [sha256.Size]byte `json:"prev"`
	Root  [sha256.Size]byte `json:"root"`
	Hash  [sha256.Size]byte `json:"hash"`
}

func NewHeader(nonce uint, prev, root [sha256.Size]byte) *Header {
	hash := sha256.Sum256(toBytes(prev, root, nonce))

	return &Header{nonce, prev, root, hash}
}

func toBytes(prev, root [sha256.Size]byte, nonce uint) []byte {
	bytes := append(prev[:], root[:]...)
	return append(bytes, byte(nonce))
}
