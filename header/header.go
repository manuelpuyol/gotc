package header

import (
	"encoding/json"
	"fmt"
	"gotc/hash"
	"gotc/utils"
	"strconv"
)

type Header struct {
	Nonce uint64
	Prev  string
	Root  string
	Hash  string
}

func NewHeader(nonce uint64, prev, root string) *Header {
	h := hash.BTCHash(prev + root + strconv.FormatUint(nonce, 10))

	return &Header{nonce, prev, root, h}
}

func (h *Header) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"nonce": h.Nonce,
		"prev":  h.Prev,
		"root":  h.Root,
		"hash":  h.Hash,
	}
}

func (h *Header) Print() {
	j, err := json.MarshalIndent(h.ToJSON(), "", "  ")
	utils.CheckErr(err)
	fmt.Println(string(j))
}
