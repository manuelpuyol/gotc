package blockchain

type Slice []*Transaction

func (ps Slice) Len() int      { return len(ps) }
func (ps Slice) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }
