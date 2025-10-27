/*
precalc.go contains declarations of precalculated attack tables, magic numbers,
huffman codes, and other useful constants.
*/

package chego

var (
	// Precalculated lookup table of LSB indices for 64-bit unsigned integers.
	//
	// See http://pradu.us/old/Nov27_2008/Buzz/research/magic/Bitboards.pdf
	// section 3.2.
	bitScanLookup = [64]int{
		63, 0, 58, 1, 59, 47, 53, 2,
		60, 39, 48, 27, 54, 33, 42, 3,
		61, 51, 37, 40, 49, 18, 28, 20,
		55, 30, 34, 11, 43, 14, 22, 4,
		62, 57, 46, 52, 38, 26, 32, 41,
		50, 36, 17, 19, 29, 10, 13, 21,
		56, 45, 25, 31, 35, 16, 9, 12,
		44, 24, 15, 8, 23, 7, 6, 5,
	}
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

	// Precalculated lookup tables used to speed up the move generation process.

	// Pawn's attack pattern depends on the color, so it is necessary to store
	// two tables.
	pawnAttacks     [2][64]uint64
	knightAttacks   [64]uint64
	kingAttacks     [64]uint64
	bishopOccupancy [64]uint64
	rookOccupancy   [64]uint64
	// Lookup bishop attack table for every possible combination of
	// square/occupancy.
	bishopAttacks [64][512]uint64
	// Lookup rook attack table for every possible combination of
	// square/occupancy.
	rookAttacks [64][4096]uint64
	// Precalculated lookup table of bishop relevant occupancy bit count for
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
	// Precalculated lookup table of rook relevant occupancy bit count for
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
	// Each path includes the king square.
	// 0 : White O-O castling path.
	// 1 : White O-O-O castling path.
	// 2 : Black O-O castling path.
	// 3 : Black O-O-O castling path.
	castlingPath = [4]uint64{
		0x70, 0x1E, 0x7000000000000000, 0x1E00000000000000,
	}
	castlingAttackPath = [4]uint64{
		0x70, 0x1C, 0x7000000000000000, 0x1C00000000000000,
	}
	// Each piece weight used to calculate material on the board.
	// Use Piece type as index to get it's weight.
	pieceWeights = [10]int{1, 1, 3, 3, 3, 3, 5, 5, 9, 9}
	// PieceSymbols maps each piece type to its symbol.
	PieceSymbols = [12]byte{
		'P', 'p', 'N', 'n', 'B', 'b',
		'R', 'r', 'Q', 'q', 'K', 'k',
	}
	// Square2String maps each board square to its string representation.
	Square2String = [64]string{
		"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
		"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
		"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
		"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
		"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
		"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
		"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
		"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8",
	}

	// Huffman codes for legal move list indices.
	// To calculate them, 10164006 games with 685863447 moves in total were
	// analyzed.  See README.md for more details.
	//
	// The maximum number of legal moves in a chess position appears to be 218,
	// so the array contains 218 elements. The latter half of the indices didn't
	// actually occur in the dataset, but each was assigned a frequency of 1 to
	// ensure valid Huffman codes are generated.
	// If a move has a frequency of 0, no code would be produced, which is
	// undesirable since the move is still theoretically possible.
	//
	// The frequencies were generated using data from the Lichess database exports
	// (https://database.lichess.org), which are released under the Creative
	// Commons CC0 license.
	huffmanCodes = [218]huffmanEntry{
		{0b1011, 4},                            // index 0 | played 35516075 times
		{0b00011, 5},                           // index 1 | played 28863637 times
		{0b1100, 4},                            // index 2 | played 33697520 times
		{0b1111, 4},                            // index 3 | played 31340990 times
		{0b00110, 5},                           // index 4 | played 26616335 times
		{0b00100, 5},                           // index 5 | played 26967376 times
		{0b00111, 5},                           // index 6 | played 26599119 times
		{0b00000, 5},                           // index 7 | played 30127529 times
		{0b00101, 5},                           // index 8 | played 26726290 times
		{0b1110, 4},                            // index 9 | played 31546838 times
		{0b01011, 5},                           // index 10 | played 21719881 times
		{0b01101, 5},                           // index 11 | played 20960808 times
		{0b01111, 5},                           // index 12 | played 20924693 times
		{0b10001, 5},                           // index 13 | played 20426220 times
		{0b10000, 5},                           // index 14 | played 20450176 times
		{0b10010, 5},                           // index 15 | played 20288330 times
		{0b01100, 5},                           // index 16 | played 21182180 times
		{0b10011, 5},                           // index 17 | played 19779373 times
		{0b01010, 5},                           // index 18 | played 22055062 times
		{0b10101, 5},                           // index 19 | played 18959904 times
		{0b11011, 5},                           // index 20 | played 16182542 times
		{0b000011, 6},                          // index 21 | played 14643685 times
		{0b000010, 6},                          // index 22 | played 15035699 times
		{0b000100, 6},                          // index 23 | played 14551558 times
		{0b010000, 6},                          // index 24 | played 12841369 times
		{0b010010, 6},                          // index 25 | played 12121516 times
		{0b011100, 6},                          // index 26 | played 11024918 times
		{0b011101, 6},                          // index 27 | played 9908166 times
		{0b101001, 6},                          // index 28 | played 9388606 times
		{0b110101, 6},                          // index 29 | played 8215047 times
		{0b0001010, 7},                         // index 30 | played 7382257 times
		{0b0100010, 7},                         // index 31 | played 6656836 times
		{0b0100011, 7},                         // index 32 | played 6157014 times
		{0b0100111, 7},                         // index 33 | played 5400835 times
		{0b1010001, 7},                         // index 34 | played 4790308 times
		{0b1101000, 7},                         // index 35 | played 4378929 times
		{0b00010110, 8},                        // index 36 | played 3779824 times
		{0b01001100, 8},                        // index 37 | played 3261509 times
		{0b01001101, 8},                        // index 38 | played 2846448 times
		{0b10100001, 8},                        // index 39 | played 2399087 times
		{0b11010011, 8},                        // index 40 | played 2045159 times
		{0b000101110, 9},                       // index 41 | played 1707181 times
		{0b101000000, 9},                       // index 42 | played 1390278 times
		{0b110100100, 9},                       // index 43 | played 1139651 times
		{0b110100101, 9},                       // index 44 | played 932421 times
		{0b0001011111, 10},                     // index 45 | played 722679 times
		{0b1010000011, 10},                     // index 46 | played 623129 times
		{0b00010111101, 11},                    // index 47 | played 423358 times
		{0b10100000101, 11},                    // index 48 | played 320010 times
		{0b000101111001, 12},                   // index 49 | played 235655 times
		{0b101000001001, 12},                   // index 50 | played 175233 times
		{0b0001011110000, 13},                  // index 51 | played 127442 times
		{0b1010000010000, 13},                  // index 52 | played 91111 times
		{0b00010111100010, 14},                 // index 53 | played 64858 times
		{0b10100000100010, 14},                 // index 54 | played 46568 times
		{0b000101111000110, 15},                // index 55 | played 31905 times
		{0b101000001000110, 15},                // index 56 | played 22068 times
		{0b0001011110001110, 16},               // index 57 | played 15412 times
		{0b1010000010001110, 16},               // index 58 | played 10561 times
		{0b00010111100011111, 17},              // index 59 | played 7044 times
		{0b10100000100011111, 17},              // index 60 | played 4775 times
		{0b000101111000111101, 18},             // index 61 | played 3372 times
		{0b101000001000111101, 18},             // index 62 | played 2320 times
		{0b1010000010001111000, 19},            // index 63 | played 1633 times
		{0b00010111100011110000, 20},           // index 64 | played 1138 times
		{0b00010111100011110011, 20},           // index 65 | played 821 times
		{0b10100000100011110011, 20},           // index 66 | played 646 times
		{0b000101111000111100100, 21},          // index 67 | played 454 times
		{0b101000001000111100101, 21},          // index 68 | played 338 times
		{0b0001011110001111000100, 22},         // index 69 | played 294 times
		{0b0001011110001111001010, 22},         // index 70 | played 207 times
		{0b1010000010001111001000, 22},         // index 71 | played 195 times
		{0b00010111100011110001010, 23},        // index 72 | played 148 times
		{0b00010111100011110001100, 23},        // index 73 | played 134 times
		{0b00010111100011110010111, 23},        // index 74 | played 90 times
		{0b10100000100011110010010, 23},        // index 75 | played 85 times
		{0b000101111000111100010111, 24},       // index 76 | played 71 times
		{0b000101111000111100011100, 24},       // index 77 | played 62 times
		{0b000101111000111100011111, 24},       // index 78 | played 54 times
		{0b000101111000111100011101, 24},       // index 79 | played 59 times
		{0b0001011110001111000111100, 25},      // index 80 | played 30 times
		{0b101000001000111100100111, 24},       // index 81 | played 42 times
		{0b0001011110001111001011000, 25},      // index 82 | played 27 times
		{0b0001011110001111001011010, 25},      // index 83 | played 26 times
		{0b0001011110001111000111101, 25},      // index 84 | played 28 times
		{0b1010000010001111001001100, 25},      // index 85 | played 22 times
		{0b1010000010001111001001101, 25},      // index 86 | played 21 times
		{0b0001011110001111001011001, 25},      // index 87 | played 27 times
		{0b00010111100011110001011001, 26},     // index 88 | played 18 times
		{0b00010111100011110001011011, 26},     // index 89 | played 16 times
		{0b00010111100011110001101010, 26},     // index 90 | played 16 times
		{0b00010111100011110010110111, 26},     // index 91 | played 12 times
		{0b00010111100011110010110110, 26},     // index 92 | played 14 times
		{0b00010111100011110001011000010, 29},  // index 93 | played 3 times
		{0b0001011110001111000101100000, 28},   // index 94 | played 6 times
		{0b0001011110001111000101101001, 28},   // index 95 | played 4 times
		{0b000101111000111100010110001, 27},    // index 96 | played 9 times
		{0b00010111100011110001011000011, 29},  // index 97 | played 3 times
		{0b00010111100011110001011010001, 29},  // index 98 | played 2 times
		{0b00010111100011110001011010000, 29},  // index 99 | played 3 times
		{0b000101111000111100011010111010, 30}, // index 100 | played 1 times
		{0b00010111100011110001011010110, 29},  // index 101 | played 2 times
		{0b000101111000111100011010111011, 30}, // index 102 | played 1 times
		{0b000101111000111100011010111000, 30}, // index 103 | played 1 times
		{0b000101111000111100011010111001, 30}, // index 104 | played 1 times
		{0b000101111000111100011010111110, 30}, // index 105 | played 1 times
		{0b000101111000111100011010111111, 30}, // index 106 | played 1 times
		{0b000101111000111100011010111100, 30}, // index 107 | played 1 times
		{0b000101111000111100011010111101, 30}, // index 108 | played 1 times
		{0b00010111100011110001011010111, 29},  // index 109 | played 2 times
		{0b000101111000111100011010110010, 30}, // index 110 | played 1 times
		{0b000101111000111100011010110011, 30}, // index 111 | played 1 times
		{0b000101111000111100011010110000, 30}, // index 112 | played 1 times
		{0b000101111000111100011010110001, 30}, // index 113 | played 1 times
		{0b000101111000111100011010110110, 30}, // index 114 | played 1 times
		{0b000101111000111100011010110111, 30}, // index 115 | played 1 times
		{0b000101111000111100011010110100, 30}, // index 116 | played 1 times
		{0b000101111000111100011010110101, 30}, // index 117 | played 1 times
		{0b000101111000111100011010001010, 30}, // index 118 | played 1 times
		{0b000101111000111100011010001011, 30}, // index 119 | played 1 times
		{0b000101111000111100011010001000, 30}, // index 120 | played 1 times
		{0b000101111000111100011010001001, 30}, // index 121 | played 1 times
		{0b000101111000111100011010001110, 30}, // index 122 | played 1 times
		{0b000101111000111100011010001111, 30}, // index 123 | played 1 times
		{0b000101111000111100011010001100, 30}, // index 124 | played 1 times
		{0b000101111000111100011010001101, 30}, // index 125 | played 1 times
		{0b000101111000111100011010000010, 30}, // index 126 | played 1 times
		{0b000101111000111100011010000011, 30}, // index 127 | played 1 times
		{0b000101111000111100011010000000, 30}, // index 128 | played 1 times
		{0b000101111000111100011010000001, 30}, // index 129 | played 1 times
		{0b000101111000111100011010000110, 30}, // index 130 | played 1 times
		{0b000101111000111100011010000111, 30}, // index 131 | played 1 times
		{0b000101111000111100011010000100, 30}, // index 132 | played 1 times
		{0b000101111000111100011010000101, 30}, // index 133 | played 1 times
		{0b000101111000111100011010011010, 30}, // index 134 | played 1 times
		{0b000101111000111100011010011011, 30}, // index 135 | played 1 times
		{0b000101111000111100011010011000, 30}, // index 136 | played 1 times
		{0b000101111000111100011010011001, 30}, // index 137 | played 1 times
		{0b000101111000111100011010011110, 30}, // index 138 | played 1 times
		{0b000101111000111100011010011111, 30}, // index 139 | played 1 times
		{0b000101111000111100011010011100, 30}, // index 140 | played 1 times
		{0b000101111000111100011010011101, 30}, // index 141 | played 1 times
		{0b000101111000111100011010010010, 30}, // index 142 | played 1 times
		{0b000101111000111100011010010011, 30}, // index 143 | played 1 times
		{0b000101111000111100011010010000, 30}, // index 144 | played 1 times
		{0b000101111000111100011010010001, 30}, // index 145 | played 1 times
		{0b000101111000111100011010010110, 30}, // index 146 | played 1 times
		{0b000101111000111100011010010111, 30}, // index 147 | played 1 times
		{0b000101111000111100011010010100, 30}, // index 148 | played 1 times
		{0b000101111000111100011010010101, 30}, // index 149 | played 1 times
		{0b000101111000111100011011101010, 30}, // index 150 | played 1 times
		{0b000101111000111100011011101011, 30}, // index 151 | played 1 times
		{0b000101111000111100011011101000, 30}, // index 152 | played 1 times
		{0b000101111000111100011011101001, 30}, // index 153 | played 1 times
		{0b000101111000111100011011101110, 30}, // index 154 | played 1 times
		{0b000101111000111100011011101111, 30}, // index 155 | played 1 times
		{0b000101111000111100011011101100, 30}, // index 156 | played 1 times
		{0b000101111000111100011011101101, 30}, // index 157 | played 1 times
		{0b000101111000111100011011100010, 30}, // index 158 | played 1 times
		{0b000101111000111100011011100011, 30}, // index 159 | played 1 times
		{0b000101111000111100011011100000, 30}, // index 160 | played 1 times
		{0b000101111000111100011011100001, 30}, // index 161 | played 1 times
		{0b000101111000111100011011100110, 30}, // index 162 | played 1 times
		{0b000101111000111100011011100111, 30}, // index 163 | played 1 times
		{0b000101111000111100011011100100, 30}, // index 164 | played 1 times
		{0b000101111000111100011011100101, 30}, // index 165 | played 1 times
		{0b000101111000111100011011111010, 30}, // index 166 | played 1 times
		{0b000101111000111100011011111011, 30}, // index 167 | played 1 times
		{0b000101111000111100011011111000, 30}, // index 168 | played 1 times
		{0b000101111000111100011011111001, 30}, // index 169 | played 1 times
		{0b000101111000111100011011111110, 30}, // index 170 | played 1 times
		{0b000101111000111100011011111111, 30}, // index 171 | played 1 times
		{0b000101111000111100011011111100, 30}, // index 172 | played 1 times
		{0b000101111000111100011011111101, 30}, // index 173 | played 1 times
		{0b000101111000111100011011110010, 30}, // index 174 | played 1 times
		{0b000101111000111100011011110011, 30}, // index 175 | played 1 times
		{0b000101111000111100011011110000, 30}, // index 176 | played 1 times
		{0b000101111000111100011011110001, 30}, // index 177 | played 1 times
		{0b000101111000111100011011110110, 30}, // index 178 | played 1 times
		{0b000101111000111100011011110111, 30}, // index 179 | played 1 times
		{0b000101111000111100011011110100, 30}, // index 180 | played 1 times
		{0b000101111000111100011011110101, 30}, // index 181 | played 1 times
		{0b000101111000111100011011001010, 30}, // index 182 | played 1 times
		{0b000101111000111100011011001011, 30}, // index 183 | played 1 times
		{0b000101111000111100011011001000, 30}, // index 184 | played 1 times
		{0b000101111000111100011011001001, 30}, // index 185 | played 1 times
		{0b000101111000111100011011001110, 30}, // index 186 | played 1 times
		{0b000101111000111100011011001111, 30}, // index 187 | played 1 times
		{0b000101111000111100011011001100, 30}, // index 188 | played 1 times
		{0b000101111000111100011011001101, 30}, // index 189 | played 1 times
		{0b000101111000111100011011000010, 30}, // index 190 | played 1 times
		{0b000101111000111100011011000011, 30}, // index 191 | played 1 times
		{0b000101111000111100011011000000, 30}, // index 192 | played 1 times
		{0b000101111000111100011011000001, 30}, // index 193 | played 1 times
		{0b000101111000111100011011000110, 30}, // index 194 | played 1 times
		{0b000101111000111100011011000111, 30}, // index 195 | played 1 times
		{0b000101111000111100011011000100, 30}, // index 196 | played 1 times
		{0b000101111000111100011011000101, 30}, // index 197 | played 1 times
		{0b000101111000111100011011011010, 30}, // index 198 | played 1 times
		{0b000101111000111100011011011011, 30}, // index 199 | played 1 times
		{0b000101111000111100011011011000, 30}, // index 200 | played 1 times
		{0b000101111000111100011011011001, 30}, // index 201 | played 1 times
		{0b000101111000111100011011011110, 30}, // index 202 | played 1 times
		{0b000101111000111100011011011111, 30}, // index 203 | played 1 times
		{0b000101111000111100011011011100, 30}, // index 204 | played 1 times
		{0b000101111000111100011011011101, 30}, // index 205 | played 1 times
		{0b000101111000111100011011010010, 30}, // index 206 | played 1 times
		{0b000101111000111100011011010011, 30}, // index 207 | played 1 times
		{0b000101111000111100011011010000, 30}, // index 208 | played 1 times
		{0b000101111000111100011011010001, 30}, // index 209 | played 1 times
		{0b000101111000111100011011010110, 30}, // index 210 | played 1 times
		{0b000101111000111100011011010111, 30}, // index 211 | played 1 times
		{0b000101111000111100011011010100, 30}, // index 212 | played 1 times
		{0b000101111000111100011011010101, 30}, // index 213 | played 1 times
		{0b000101111000111100010110101010, 30}, // index 214 | played 1 times
		{0b000101111000111100010110101011, 30}, // index 215 | played 1 times
		{0b000101111000111100010110101000, 30}, // index 216 | played 1 times
		{0b000101111000111100010110101001, 30}, // index 217 | played 1 times
	}
)

// Standard initial chess position.
const InitialPos = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Bitboards of each square.
const (
	A1 uint64 = 1 << iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1
	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
	ALL_SQUARES = 0xFFFFFFFFFFFFFFFF
)

// Indicies of each square.
const (
	SA1 int = iota
	SB1
	SC1
	SD1
	SE1
	SF1
	SG1
	SH1
	SA2
	SB2
	SC2
	SD2
	SE2
	SF2
	SG2
	SH2
	SA3
	SB3
	SC3
	SD3
	SE3
	SF3
	SG3
	SH3
	SA4
	SB4
	SC4
	SD4
	SE4
	SF4
	SG4
	SH4
	SA5
	SB5
	SC5
	SD5
	SE5
	SF5
	SG5
	SH5
	SA6
	SB6
	SC6
	SD6
	SE6
	SF6
	SG6
	SH6
	SA7
	SB7
	SC7
	SD7
	SE7
	SF7
	SG7
	SH7
	SA8
	SB8
	SC8
	SD8
	SE8
	SF8
	SG8
	SH8
)
