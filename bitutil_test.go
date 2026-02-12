package chego

import (
	"math/bits"
	"testing"
)

func TestBitWriter(t *testing.T) {
	bw := NewBitWriter()

	expected := 0
	for i := 1; i <= 64; i++ {
		size := bits.Len64(uint64(i))
		expected += size
		bw.Write(uint(i), size)
	}

	got := bw.buff.Len()*8 + (intSize - bw.remainingBits)
	if got != expected {
		t.Fatalf("Expected %d bits. Buffer has %d bits\n", expected, got)
	}

	expectedBuff := []byte{
		0b11011100, 0b10111011, 0b11000100, 0b11010101, 0b11100110, 0b11110111,
		0b11000010, 0b00110010, 0b10011101, 0b00101011, 0b01101011, 0b11100011,
		0b00111010, 0b11011111, 0b00111011, 0b11101111, 0b11000001, 0b00001100,
		0b01010001, 0b11001001, 0b00101100, 0b11010011, 0b11010001, 0b01001101,
		0b01010101, 0b11011001, 0b01101101, 0b11010111, 0b11100001, 0b10001110,
		0b01011001, 0b11101001, 0b10101110, 0b11011011, 0b11110001, 0b11001111,
		0b01011101, 0b11111001, 0b11101111, 0b11011111, 0b11000000,
	}
	for i, b := range bw.Bytes() {
		if expectedBuff[i] != b {
			t.Fatalf("Expected %b, got %b", expectedBuff[i], b)
		}
	}

}

func TestBitScan(t *testing.T) {
	for i := range 64 {
		bb := uint64(1 << i)

		got := bitScan(bb)
		if got != i {
			t.Fatalf("Expected: %d got %d", i, got)
		}
	}
}

func TestPopLSB(t *testing.T) {
	for i := range 64 {
		bb := uint64(1 << i)

		got := popLSB(&bb)
		if got != i {
			t.Fatalf("Expected %d got %d", i, got)
		}
	}
}

func TestCountBits(t *testing.T) {
	bb := uint64(0)

	for i := range 64 {
		bb |= uint64(1 << i)

		got := CountBits(bb)
		if got != i+1 {
			t.Fatalf("Expected: %d got %d", i+1, got)
		}
	}
}

func BenchmarkBitWriter(b *testing.B) {
	bw := NewBitWriter()
	for b.Loop() {
		for i := 1; i <= 64; i++ {
			bw.Write(uint(i), bits.Len64(uint64(i)))
		}
	}
}

func BenchmarkBitScan(b *testing.B) {
	for b.Loop() {
		bitScan(0x8000000000000000)
	}
}

func BenchmarkPopLSB(b *testing.B) {
	var bitboard uint64 = 0xFFFFFFFFFFFFFFFF

	for b.Loop() {
		popLSB(&bitboard)
	}
}

func BenchmarkCountBits(b *testing.B) {
	for b.Loop() {
		CountBits(0xFFFFFFFFFFFFFFFF)
	}
}
