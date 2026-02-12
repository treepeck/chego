// init.go contains declaration and initialization of precalculated attack
// tables, zobrist keys, and magic numbers.

package chego

import "math/rand/v2"

// initBishopOccupancy initializes the lookup table of the "relevant occupancy
// squares" for a bishop.  They are the only squares whose occupancy matters when
// generating legal moves of a bishop.
func initBishopOccupancy() [64]uint64 {
	var result [64]uint64
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

		result[square] = occupancy
	}
	return result
}

// initRookOccupancy initializes the lookup table of the "relevant occupancy
// squares" for a rook.  They are the only squares whose occupancy matters when
// generating legal moves of a rook.
func initRookOccupancy() [64]uint64 {
	var result [64]uint64
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

		result[square] = occupancy
	}
	return result
}

// Initializes the lookup table of pawn attacks for every possible square on the
// chessboard.  A pawn's attack pattern depends on its color, so two tables are
// required.
func initPawnAttacks() [2][64]uint64 {
	var attacks [2][64]uint64
	for i := range 64 {
		square := uint64(1 << i)
		attacks[ColorWhite][i] = genPawnAttacks(square, ColorWhite)
		attacks[ColorBlack][i] = genPawnAttacks(square, ColorBlack)
	}
	return attacks
}

// Initializes the lookup table of knight attacks for every possible square on the
// chessboard.
func initKnightAttacks() [64]uint64 {
	var attacks [64]uint64
	for i := range 64 {
		attacks[i] = genKnightAttacks(1 << i)
	}
	return attacks
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

// Initializes the lookup table of bishop attacks for every possible combination
// of square and occupancy.  It's the caller's responsibility to ensure that
// [bishopOccupancy] array is initialized by calling the [initBishopOccupancy]
// function.
func initBishopAttacks() [64][512]uint64 {
	var attacks [64][512]uint64
	for i := range 64 {
		bitCount := bishopBitCount[i]
		for j := 0; j < 1<<bitCount; j++ {
			occupancy := genOccupancy(j, bitCount, bishopOccupancy[i])

			key := occupancy * bishopMagicNumbers[i] >> (64 - bitCount)

			attacks[i][key] = genBishopAttacks(1<<i, occupancy)
		}
	}
	return attacks
}

// Initializes the lookup table of rook attacks for every possible combination
// of square and occupancy. It's the caller's responsibility to ensure that
// [rookOccupancy] array is initialized by calling the [initRookOccupancy] function.
func initRookAttacks() [64][4096]uint64 {
	var attacks [64][4096]uint64
	for i := range 64 {
		bitCount := rookBitCount[i]
		for j := 0; j < 1<<bitCount; j++ {
			occupancy := genOccupancy(j, bitCount, rookOccupancy[i])

			key := occupancy * rookMagicNumbers[i] >> (64 - bitCount)

			attacks[i][key] = genRookAttacks(1<<i, occupancy)
		}
	}
	return attacks
}

// Initializes the lookup table of king attacks for every possible square on the
// chessboard.
func initKingAttacks() [64]uint64 {
	var attacks [64]uint64
	for i := range 64 {
		attacks[i] = genKingAttacks(1 << i)
	}
	return attacks
}

// Initializes the piece placement keys for the Zobrist hashing scheme.
func initPieceKeys() [12][64]uint64 {
	var keys [12][64]uint64
	for i := WPawn; i <= BKing; i++ {
		for square := range 64 {
			keys[i][square] = rand.Uint64()
		}
	}
	return keys
}

// Initializes the en passant keys for the Zobrist hashing scheme.
func initEnPassantKeys() [64]uint64 {
	var keys [64]uint64
	for square := range 64 {
		keys[square] = rand.Uint64()
	}
	return keys
}

// Initializes the caslting keys for the Zobrist hashing scheme.
func initCastlingKeys() [16]uint64 {
	var keys [16]uint64
	for i := range 16 {
		keys[i] = rand.Uint64()
	}
	return keys
}

// Precalculated lookup tables used to speed up the move generation process.
var (
	// Leaper pieces attacks.

	pawnAttacks   = initPawnAttacks()
	knightAttacks = initKnightAttacks()
	kingAttacks   = initKingAttacks()

	// Precalculated lookup table of the bishop relevant occupancy bit count for
	// every square.
	bishopBitCount = [64]int{
		6, 5, 5, 5, 5, 5, 5, 6,
		5, 5, 5, 5, 5, 5, 5, 5,
		5, 5, 7, 7, 7, 7, 5, 5,
		5, 5, 7, 9, 9, 7, 5, 5,
		5, 5, 7, 9, 9, 7, 5, 5,
		5, 5, 7, 7, 7, 7, 5, 5,
		5, 5, 5, 5, 5, 5, 5, 5,
		6, 5, 5, 5, 5, 5, 5, 6,
	}

	// Precalculated lookup table of the rook relevant occupancy bit count for
	// every square.
	rookBitCount = [64]int{
		12, 11, 11, 11, 11, 11, 11, 12,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		12, 11, 11, 11, 11, 11, 11, 12,
	}

	// Occupancies to initialize the sliding pieces attacks.

	bishopOccupancy = initBishopOccupancy()
	rookOccupancy   = initRookOccupancy()

	// Sliding pieces attacks.

	bishopAttacks = initBishopAttacks()
	rookAttacks   = initRookAttacks()
)

// Zobrist Keys are used to hash each possible position into the unique number.
// Each key is generated randomly and large enough, so the probability of hash
// collisions is negligible.
var (
	pieceKeys = initPieceKeys()
	// Used only when black is the active color.
	epKeys       = initEnPassantKeys()
	castlingKeys = initCastlingKeys()
	// Used only when black is the active color.
	colorKey = rand.Uint64()
)

var (
	// bishopMagicNumbers is a precalculated lookup table of magic
	// numbers for a bishop.
	bishopMagicNumbers = [64]uint64{
		0x11410121040100,
		0x2084820928010,
		0xa010208481080040,
		0x214240082000610,
		0x4d104000400480,
		0x1012010804408,
		0x42044101452000c,
		0x2844804050104880,
		0x814204290a0a00,
		0x10280688224500,
		0x1080410101010084,
		0x10020a108408004,
		0x2482020210c80080,
		0x480104a0040400,
		0x411006404200810,
		0x1024010908024292,
		0x1004401001011a,
		0x810006081220080,
		0x1040404206004100,
		0x58080000820041ce,
		0x3406000422010890,
		0x1a004100520210,
		0x202a000048040400,
		0x225004441180110,
		0x8064240102240,
		0x1424200404010402,
		0x1041100041024200,
		0x8082002012008200,
		0x1010008104000,
		0x8808004000806000,
		0x380a000080c400,
		0x31040100042d0101,
		0x110109008082220,
		0x4010880204201,
		0x4006462082100300,
		0x4002010040140041,
		0x40090200250880,
		0x2010100c40c08040,
		0x12800ac01910104,
		0x10b20051020100,
		0x210894104828c000,
		0x50440220004800,
		0x1002011044180800,
		0x4220404010410204,
		0x1002204a2020401,
		0x21021001000210,
		0x4880081009402,
		0xc208088c088e0040,
		0x4188464200080,
		0x3810440618022200,
		0xc020310401040420,
		0x2000008208800e0,
		0x4c910240020,
		0x425100a8602a0,
		0x20c4206a0c030510,
		0x4c10010801184000,
		0x200202020a026200,
		0x6000004400841080,
		0xc14004121082200,
		0x400324804208800,
		0x1802200040504100,
		0x1820000848488820,
		0x8620682a908400,
		0x8010600084204240,
	}
	// rookMagicNumbers is a precalculated lookup table of magic
	// numbers for a rook.
	rookMagicNumbers = [64]uint64{
		0x2080008040002010,
		0x40200010004000,
		0x100090010200040,
		0x2080080010000480,
		0x880040080080102,
		0x8200106200042108,
		0x410041000408b200,
		0x100009a00402100,
		0x5800800020804000,
		0x848404010002000,
		0x101001820010041,
		0x10a0040100420080,
		0x8a02002006001008,
		0x926000844110200,
		0x8000800200800100,
		0x28060001008c2042,
		0x10818002204000,
		0x10004020004001,
		0x110002008002400,
		0x11a020010082040,
		0x2001010008000410,
		0x42010100080400,
		0x4004040008020110,
		0x820000840041,
		0x400080208000,
		0x2080200040005000,
		0x8000200080100080,
		0x4400080180500080,
		0x4900080080040080,
		0x4004004480020080,
		0x8006000200040108,
		0xc481000100006396,
		0x1000400080800020,
		0x201004400040,
		0x10008010802000,
		0x204012000a00,
		0x800400800802,
		0x284000200800480,
		0x3000403000200,
		0x840a6000514,
		0x4080c000228012,
		0x10002000444010,
		0x620001000808020,
		0xc210010010009,
		0x100c001008010100,
		0xc10020004008080,
		0x20100802040001,
		0x808008305420014,
		0xc010800840043080,
		0x208401020890100,
		0x10b0081020028280,
		0x6087001001220900,
		0xc080011000500,
		0x9810200040080,
		0x2000010882100400,
		0x2000050880540200,
		0x800020104200810a,
		0x6220250242008016,
		0x9180402202900a,
		0x40210500100009,
		0x6000814102026,
		0x410100080a040013,
		0x10405008022d1184,
		0x1000009400410822,
	}
)
