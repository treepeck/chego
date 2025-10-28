/*
bitutil.go implements useful bit utilities which are used in move generation,
huffman coding, and game management.
*/

package chego

import "bytes"

const (
	// For x86-64 CPUs int size is 32 bits. For x64 CPUs int size is 64 bits.
	intSize = (32 << (^uint(0) >> 63))
	// Precalculated magic used to form indices for the bitScanLookup array.
	bitscanMagic uint64 = 0x07EDD5E59A4E28C2
)

/*
BitWriter writes and stores the bit set (aka bit array) of arbitrary size.
Internally, bytes.Buffer is used to prevent excessive memory allocations and
ensure efficient appending of multiple bit chunks. Bit chunks smaller than the
size of an integer are stored in the internal field.
*/
type BitWriter struct {
	buff          bytes.Buffer
	temp          uint
	remainingBits int
}

func NewBitWriter() *BitWriter { return &BitWriter{remainingBits: intSize} }

/*
Write writes data with size bits to the BitWriter.  If size is less than or
equal to the integer size (which depends on the CPU architecture), the data is
stored in an internal integer field. When the field overflows, its contents are
flushed to the internal bytes.Buffer.
*/
func (bw *BitWriter) Write(data uint, size int) {
	// Set all unused bits in data to 0.
	data &= 1<<size - 1

	bw.remainingBits -= size
	if bw.remainingBits >= 0 {
		bw.temp |= data << bw.remainingBits
	} else {
		bw.temp |= data >> -bw.remainingBits
		// Split integer into the byte sequence.
		for i := range intSize / 8 {
			chunk := byte(bw.temp >> (i * 8) & 0xFF)
			// Don't handle error since WriteByte always returns nil.
			bw.buff.WriteByte(chunk)
		}
		bw.remainingBits += intSize
	}
	bw.temp = data << bw.remainingBits
}

/*
CountBits returns the number of bits set within the bitboard.
*/
func CountBits(bitboard uint64) (cnt int) {
	for ; bitboard > 0; cnt++ {
		bitboard &= bitboard - 1
	}
	return cnt
}

/*
bitScan returns the index of the LSB withing the bitboard.  bitboard & -bitboard
gives the LSB which is then run through the hashing scheme to index a lookup.

NOTE: bitScan returns 63 for the empty bitboard.
*/
func bitScan(bitboard uint64) int {
	return bitScanLookup[bitboard&-bitboard*bitscanMagic>>58]
}

/*
popLSB removes the LSB from the bitboard and returns its index.

NOTE: popLSB returns 63 for the empty bitboard.
*/
func popLSB(bitboard *uint64) int {
	lsb := bitScan(*bitboard)
	*bitboard &= *bitboard - 1
	return lsb
}
