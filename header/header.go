package header

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type Header struct {
	Nonce uint
	Prev  [sha256.Size]byte
	Root  [sha256.Size]byte
	Hash  [sha256.Size]byte
}

func NewHeader(nonce uint, prev, root [sha256.Size]byte) *Header {
	hash := sha256.Sum256(toBytes(prev, root, nonce))

	return &Header{nonce, prev, root, hash}
}

func (h *Header) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"nonce": h.Nonce,
		"prev":  fmt.Sprintf("%x", h.Prev),
		"root":  fmt.Sprintf("%x", h.Root),
		"hash":  fmt.Sprintf("%x", h.Hash),
	}
}

func (h *Header) Print() {
	j, _ := json.MarshalIndent(h.ToJSON(), "", "  ")
	fmt.Println(string(j))
}

func toBytes(prev, root [sha256.Size]byte, nonce uint) []byte {
	bytes := append(prev[:], root[:]...)
	return append(bytes, byte(nonce))
}
