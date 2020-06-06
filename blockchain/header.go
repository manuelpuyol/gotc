package blockchain

import (
	"encoding/json"
	"fmt"
	"gotc/hash"
	"gotc/utils"
	"strconv"
)

type Header struct {
	Nonce uint32 // number used for minig
	Prev  string // hash of the last block
	Root  string // merkle root of the transaction list
	Hash  string // hash of the block
}

func NewHeader(nonce uint32, prev, root string) *Header {
	h := hash.BTCHash(prev + root + strconv.FormatUint(uint64(nonce), 10))

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
