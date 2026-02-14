package chego

import (
	"testing"
)

func TestHuffmanEncoding(t *testing.T) {
	expected := []byte{
		0b11101110, 0b00001010, 0b01100000,
	}

	got := HuffmanEncoding([]int{9, 9, 22, 17})

	for i, b := range expected {
		if got[i] != b {
			t.Fatalf("expected: %v, got: %v", expected, got)
		}
	}
}

func TestHuffmanDecoding(t *testing.T) {
	indices := []int{9, 9, 22, 17}

	w := bitWriter{remainingBits: intSize}
	for _, i := range indices {
		w.write(huffmanCodes[i].code, huffmanCodes[i].size)
	}

	expected := []DecodedMove{
		{NewMove(SE4, SE2, MoveNormal), "e4"},
		{NewMove(SE5, SE7, MoveNormal), "e5"},
		{NewMove(SC4, SF1, MoveNormal), "Bc4"},
		{NewMove(SF6, SG8, MoveNormal), "Nf6"},
	}

	got := HuffmanDecoding(w.content(), len(indices))

	for i, m := range expected {
		if got[i] != m {
			t.Fatalf("expected: %v, got: %v", expected, got)
		}
	}
}

func BenchmarkMakeTrie(b *testing.B) {
	for b.Loop() {
		makeTrie()
	}
}
