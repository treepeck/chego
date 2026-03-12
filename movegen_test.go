package chego

import (
	"testing"
)

func BenchmarkGenPawnAttacks(b *testing.B) {
	for b.Loop() {
		genPawnAttacks(B4, ColorWhite)
	}
}

func BenchmarkGenKnightAttacks(b *testing.B) {
	for b.Loop() {
		genKnightAttacks(B4)
	}
}

func BenchmarkGenKingAttakcs(b *testing.B) {
	for b.Loop() {
		genKingAttacks(B4)
	}
}

func BenchmarkGenBishopAttacks(b *testing.B) {
	for b.Loop() {
		genBishopAttacks(D5, B3)
	}
}

func BenchmarkGenRookAttacks(b *testing.B) {
	for b.Loop() {
		genRookAttacks(D5, B3)
	}
}

func BenchmarkInitBishopOccupancy(b *testing.B) {
	for b.Loop() {
		initBishopOccupancy()
	}
}

func BenchmarkInitRookOccupancy(b *testing.B) {
	for b.Loop() {
		initRookOccupancy()
	}
}

func BenchmarkLookupBishopAttacks(b *testing.B) {
	for b.Loop() {
		lookupBishopAttacks(35, 0x0)
	}
}

func BenchmarkLookupRookAttacks(b *testing.B) {
	for b.Loop() {
		lookupRookAttacks(35, 0x0)
	}
}

func BenchmarkLookupQueenAttacks(b *testing.B) {
	for b.Loop() {
		lookupQueenAttacks(35, 0x0)
	}
}

func BenchmarkGenKingMoves(b *testing.B) {
	pos := ParseFen("8/8/8/8/8/8/8/R3K2R w - - 0 1")

	for b.Loop() {
		genKingMoves(*pos, &MoveList{})
	}
}

func BenchmarkGenLegalMoves(b *testing.B) {
	pos := ParseFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")

	for b.Loop() {
		lm := MoveList{}
		GenLegalMoves(*pos, &lm)
	}
}
