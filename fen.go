package chego

import (
	"strconv"
	"strings"
)

// FEN of the standard initial chess position.
const InitialPos = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var (
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
	// PieceSymbols maps each piece type to its symbol.
	PieceSymbols = [12]byte{
		'P', 'p', 'N', 'n', 'B', 'b',
		'R', 'r', 'Q', 'q', 'K', 'k',
	}
)

// ParseFEN parses the given FEN string into a [Position].  It's the caller's
// responsibility to validate the provided FEN string.
//
// Each FEN string consists of six parts, separated by a space:
//  1. Piece placement: will be parsed into the array of bitboards.
//  2. Active color:
//     "w" means that White is to move;
//     "b" means that Black is to move.
//  3. Castling rights: if neither side has the ability to castle,
//     this field uses the character "-".
//  4. En passant target square: if there is no en passant target square,
//     this field uses the character "-".
//  5. Halfmove clock: used for the fifty-move rule.
//  6. Fullmove number: The number of the full moves.
func ParseFEN(fen string) (p Position) {
	// Separate FEN fields.
	fields := strings.SplitN(fen, " ", 6)

	// Parse piece placement.
	p.Bitboards = ParseBitboards(fields[0])

	// Parse active color.
	// p will have ColorWhite by default.
	if fields[1] == "b" {
		p.ActiveColor = ColorBlack
	}

	// Parse castling rights.
	for i := range len(fields[2]) {
		switch fields[2][i] {
		case 'K':
			p.CastlingRights |= CastlingWhiteShort
		case 'Q':
			p.CastlingRights |= CastlingWhiteLong
		case 'k':
			p.CastlingRights |= CastlingBlackShort
		case 'q':
			p.CastlingRights |= CastlingBlackLong
		}
	}

	// Parse en passant target square.
	for i := range Square2String {
		if Square2String[i] == fields[3] {
			p.EPTarget = i
		}
	}

	// Parse halfmove counter.
	var err error
	p.HalfmoveCnt, err = strconv.Atoi(fields[4])
	if err != nil {
		panic("cannot parse halfmove counter from FEN string")
	}

	// Parse fullmove counter.
	p.FullmoveCnt, err = strconv.Atoi(fields[5])
	if err != nil {
		panic("cannot parse fullmove counter from FEN string")
	}

	return p
}

// SerializeFEN serializes the specified [Position] into a FEN string.  It's
// the caller's responsibility to validate the provided position.
func SerializeFEN(p Position) string {
	var fen strings.Builder
	fen.Grow(64)

	// 1 field: piece placement.
	fen.WriteString(SerializeBitboards(p.Bitboards))

	// 2 field: active color.
	if p.ActiveColor == ColorWhite {
		fen.WriteString(" w ")
	} else {
		fen.WriteString(" b ")
	}

	// 3 field: castling rights.
	cnt := 4
	if p.CastlingRights&CastlingWhiteShort != 0 {
		fen.WriteByte('K')
		cnt--
	}
	if p.CastlingRights&CastlingWhiteLong != 0 {
		fen.WriteByte('Q')
		cnt--
	}
	if p.CastlingRights&CastlingBlackShort != 0 {
		fen.WriteByte('k')
		cnt--
	}
	if p.CastlingRights&CastlingBlackLong != 0 {
		fen.WriteByte('q')
		cnt--
	}
	if cnt == 4 {
		fen.WriteByte('-')
	}
	fen.WriteByte(' ')

	// 4 field: en passant target square.
	if p.EPTarget == 0 {
		fen.WriteString("- ")
	} else {
		files := "abcdefgh"
		fen.WriteByte(files[p.EPTarget%8])
		fen.WriteByte('0' + byte(p.EPTarget/8+1))
		fen.WriteByte(' ')
	}

	// 5 field: the number of halfmoves.
	fen.WriteString(strconv.Itoa(p.HalfmoveCnt))
	fen.WriteByte(' ')

	// 6 field: the number of fullmoves.
	fen.WriteString(strconv.Itoa(p.FullmoveCnt))

	return fen.String()
}

// ParseBitboards converts the first part of a FEN string into an array of
// bitboards.  May panic if the provided string is not valid.
func ParseBitboards(piecePlacement string) (bitboards [15]uint64) {
	square := 56

	// Piece placement data describes each rank beginning from the eigth.
	for i := range len(piecePlacement) {
		char := piecePlacement[i]

		if char == '/' { // Rank separator.
			square -= 16
			// Number of consecutive empty squares.
		} else if char >= '1' && char <= '8' {
			// Convert byte to the integer it represents.
			square += int(char - '0')
		} else { // There is piece on a square.
			piece := PieceWPawn
			// Manual switch construction is ~3x faster than map approach.
			switch char {
			case 'N':
				piece = PieceWKnight
			case 'B':
				piece = PieceWBishop
			case 'R':
				piece = PieceWRook
			case 'Q':
				piece = PieceWQueen
			case 'K':
				piece = PieceWKing
			case 'p':
				piece = PieceBPawn
			case 'n':
				piece = PieceBKnight
			case 'b':
				piece = PieceBBishop
			case 'r':
				piece = PieceBRook
			case 'q':
				piece = PieceBQueen
			case 'k':
				piece = PieceBKing
			}
			// Set the bit on the bitboards to place a piece.
			bb := uint64(1 << square)

			bitboards[piece] |= bb
			if piece%2 == 0 {
				bitboards[12] |= bb
			} else {
				bitboards[13] |= bb
			}
			bitboards[14] |= bb

			square++
		}
	}

	return bitboards
}

// SerializeBitboards converts the array of bitboards into the first part of FEN
// string.
func SerializeBitboards(bitboards [15]uint64) string {
	// Used to add characters to a string without extra memory allocations.
	b := strings.Builder{}
	b.Grow(30)

	board := [64]byte{}

	for i := 0; i <= PieceBKing; i++ {
		// Go through all pieces on a bitboard.
		for bitboards[i] > 0 {
			square := popLSB(&bitboards[i])
			// Add piece on board.
			board[square] = PieceSymbols[i]
		}
	}

	emptySquares := byte(0)
	for rank := 7; rank >= 0; rank-- {
		for file := range 8 {
			square := 8*rank + file
			if board[square] == 0 { // Empty square.
				emptySquares++
			} else { // Piece on square.
				if emptySquares > 0 {
					b.WriteByte('0' + emptySquares)
					emptySquares = 0
				}
				b.WriteByte(board[square])
			}

			// To add rank separators.
			if (square+1)%8 == 0 {
				if emptySquares > 0 {
					b.WriteByte('0' + emptySquares)
					emptySquares = 0
				}
				// Do not add separator in the end of the string.
				if square != 7 {
					b.WriteByte('/')
				}
			}
		}
	}

	return b.String()
}
