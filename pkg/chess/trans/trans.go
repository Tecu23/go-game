package trans

var randPcSq [12 * 64]uint64 // keyvalues for 'pc on sq'

// PcSqKey returns the keyvalue for the piece on square
func PcSqKey(pc, sq int) uint64 {
	return randPcSq[pc*64+sq]
}
