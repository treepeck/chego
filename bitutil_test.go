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
		t.Logf("Expected %d bits. Buffer has %d bits\n", expected, got)
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
