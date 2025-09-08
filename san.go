/*
san.go implements serialization of moves into Standard Algebraic Notation.
See https://ia802908.us.archive.org/26/items/pgn-standard-1994-03-12/PGN_standard_1994-03-12.txt Section 8.2.3.
*/

package chego

import (
	"strings"
)

/*
Move2SAN encodes the specified move to its SAN representation.

SAN string consists of these parts:
 1. Piece name, omitted for for pawns;
 2. Optional originating (source) file or rank, used for disambiguation.  If a
    pawn performs a capture, its originating file is always included;
 3. Denotation of capture by 'x'. Mandatory for capture moves;
 4. Destination (to) file and rank;
 5. Denotation of check by '+'. Omitted when the move is a checkmate;
 6. Denotation of checkmate by '#'.

King castling and queen castling are encoded as "O-O" and "O-O-O" respectively.
*/
func Move2SAN(
	m Move,
	pos *Position,
	legalMoves MoveList,
	p Piece,
	isCapture, isCheck, isCheckmate bool,
) string {
	if m.Type() == MoveCastling {
		if m.To() == SA1 || m.To() == SA8 {
			return "O-O-O"
		} else {
			return "O-O"
		}
	}

	var b strings.Builder
	b.Grow(2)

	switch p {
	case PieceWKnight, PieceBKnight:
		b.WriteByte('N')
	case PieceWBishop, PieceBBishop:
		b.WriteByte('B')
	case PieceWRook, PieceBRook:
		b.WriteByte('R')
	case PieceWQueen, PieceBQueen:
		b.WriteByte('Q')
	case PieceWKing, PieceBKing:
		b.WriteByte('K')
	}

	// Resolve the ambiguity if needed.  Skip the pawns since their moves are
	// always ambiguous.
	if p > PieceBPawn {
		for i := range legalMoves.LastMoveIndex {
			if pos.GetPieceFromSquare(1<<legalMoves.Moves[i].From()) == p &&
				legalMoves.Moves[i].To() == m.To() &&
				legalMoves.Moves[i].From() != m.From() {
				b.WriteByte(disambiguate(m.From(), legalMoves.Moves[i].From()))
				break
			}
		}
	}

	if isCapture {
		if p <= PieceBPawn {
			b.WriteByte(files[m.From()%8])
		}
		b.WriteByte('x')
	}

	// Append destination square.
	b.WriteString(Square2String[m.To()])

	// Append promotion info.
	if m.Type() == MovePromotion {
		switch m.PromoPiece() {
		case PromotionKnight:
			b.WriteString("=N")
		case PromotionBishop:
			b.WriteString("=B")
		case PromotionRook:
			b.WriteString("=R")
		case PromotionQueen:
			b.WriteString("=Q")
		}
	}

	if isCheckmate {
		b.WriteByte('#')
	} else if isCheck {
		b.WriteByte('+')
	}

	return b.String()
}

/*
disambiguate resolves the ambiguity that arises when multiple pieces of the same
type can move to the same square.

Steps to resolve:
 1. If the moving pieces can be distinguished by their originating files,
    the originating file letter of the moving piece is inserted immediately
    after the moving piece letter;
 2. If the moving pieces can be distinguished by their originating ranks,
    the originating rank digit of the moving piece is inserted immediately
    after the moving piece letter.
*/
func disambiguate(fromA, fromB int) byte {
	if fromA%8 != fromB%8 {
		return files[fromA%8]
	}
	if fromA/8 != fromB/8 {
		return byte(fromA/8 + 1 + '0')
	}
	panic("cannot disambiguate the move")
}
