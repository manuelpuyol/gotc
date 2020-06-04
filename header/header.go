package header

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gotc/utils"
	"strconv"
)

type Header struct {
	Nonce uint64
	Prev  [sha256.Size]byte
	Root  [sha256.Size]byte
	Hash  [sha256.Size]byte
}

func NewHeader(nonce uint64, prev, root [sha256.Size]byte) *Header {
	hash := sha256.Sum256(toBytes(prev, root, nonce))

	return &Header{nonce, prev, root, hash}
}

func (h *Header) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"nonce": h.Nonce,
		"prev":  utils.SHAToString(h.Prev),
		"root":  utils.SHAToString(h.Root),
		"hash":  utils.SHAToString(h.Hash),
	}
}

func (h *Header) Print() {
	j, _ := json.MarshalIndent(h.ToJSON(), "", "  ")
	fmt.Println(string(j))
}

func toBytes(prev, root [sha256.Size]byte, nonce uint64) []byte {
	pstr := utils.SHAToString(prev)
	rstr := utils.SHAToString(root)
	nstr := strconv.FormatUint(nonce, 10)

	str := pstr + rstr + nstr
	return []byte(str)
}
