package chego

import (
	"math/rand/v2"
)

// initBishopOccupancy initializes the lookup table of the "relevant occupancy
// squares" for a bishop.  They are the only squares whose occupancy matters when
// generating legal moves of a bishop.
func initBishopOccupancy() {
	// Helper constants.
	const not_A_not_1st = notAFile & not1stRank
	const not_H_not_1st = notHFile & not1stRank
	const not_A_not_8th = notAFile & not8thRank
	const not_H_not_8th = notHFile & not8thRank

	for square := range 64 {
		var occupancy, bishop uint64 = 0, 1 << square

		for i := bishop & notAFile >> 9; i&not_A_not_1st != 0; i >>= 9 {
			occupancy |= i
		}

		for i := bishop & notHFile >> 7; i&not_H_not_1st != 0; i >>= 7 {
			occupancy |= i
		}

		for i := bishop & notAFile << 7; i&not_A_not_8th != 0; i <<= 7 {
			occupancy |= i
		}

		for i := bishop & notHFile << 9; i&not_H_not_8th != 0; i <<= 9 {
			occupancy |= i
		}

		bishopOccupancy[square] = occupancy
	}
}

// initRookOccupancy initializes the lookup table of the "relevant occupancy
// squares" for a rook.  They are the only squares whose occupancy matters when
// generating legal moves of a rook.
func initRookOccupancy() {
	for square := range 64 {
		var occupancy, rook uint64 = 0, 1 << square

		for i := rook & not1stRank >> 8; i&not1stRank != 0; i >>= 8 {
			occupancy |= i
		}

		for i := rook & notAFile >> 1; i&notAFile != 0; i >>= 1 {
			occupancy |= i
		}

		for i := rook & notHFile << 1; i&notHFile != 0; i <<= 1 {
			occupancy |= i
		}

		for i := rook & not8thRank << 8; i&not8thRank != 0; i <<= 8 {
			occupancy |= i
		}

		rookOccupancy[square] = occupancy
	}
}

// initZobristKeys initializes the pseudo-random keys used in the Zobrist
// hashing scheme.
func initZobristKeys() {
	for i := WPawn; i <= BKing; i++ {
		for square := range 64 {
			pieceKeys[i][square] = rand.Uint64()
		}
	}

	for square := range 64 {
		epKeys[square] = rand.Uint64()
	}

	for i := range 16 {
		castlingKeys[i] = rand.Uint64()
	}

	colorKey = rand.Uint64()
}

// genOccupancy returns a bitboard of blocker pieces for the specified attack
// bitboard.
func genOccupancy(key, relevantBitCount int,
	relevantOccupancy uint64) (occupancy uint64) {

	for i := range relevantBitCount {
		square := popLSB(&relevantOccupancy)

		if key&(1<<i) != 0 {
			occupancy |= 1 << square
		}
	}

	return occupancy
}

// initAttackTables initializes the predefined attack tables.
func initAttackTables() {
	initBishopOccupancy()
	initRookOccupancy()

	for square := range 64 {
		bb := uint64(1 << square)

		pawnAttacks[ColorWhite][square] = genPawnAttacks(bb, ColorWhite)
		pawnAttacks[ColorBlack][square] = genPawnAttacks(bb, ColorBlack)

		knightAttacks[square] = genKnightAttacks(bb)

		kingAttacks[square] = genKingAttacks(bb)

		bitCount := bishopBitCount[square]
		for i := 0; i < 1<<bitCount; i++ {
			occupancy := genOccupancy(i, bitCount, bishopOccupancy[square])

			key := occupancy * bishopMagicNumbers[square] >> (64 - bitCount)

			bishopAttacks[square][key] = genBishopAttacks(bb, occupancy)
		}

		bitCount = rookBitCount[square]
		for i := 0; i < 1<<bitCount; i++ {
			occupancy := genOccupancy(i, bitCount, rookOccupancy[square])

			key := occupancy * rookMagicNumbers[square] >> (64 - bitCount)

			rookAttacks[square][key] = genRookAttacks(bb, occupancy)
		}
	}
}

// init is a special function that runs automatically once per
// package, regardless of how many times the package is imported.
func init() {
	initAttackTables()
	initZobristKeys()
}
