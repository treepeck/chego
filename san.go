// san.go implements serialization of moves into Standard Algebraic Notation.
// See http://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm Section 8.2.3.

package chego

import "strings"

// Move2SAN encodes the specified move to its SAN representation.
//
// SAN string consists of these parts:
//  1. Piece name, omitted for for pawns;
//  2. Optional originating (source) file, rank, or both, used for disambiguation.
//     If a pawn performs a capture, its originating file is always included;
//  3. Denotation of capture by 'x'. Mandatory for capture moves;
//  4. Destination (to) file and rank;
//  5. Denotation of check by '+'. Omitted when the move is a checkmate;
//  6. Denotation of checkmate by '#'.
//
// NOTE: Position will be modified by applying the specified move to denote checks
// and checkmates.  MoveList will also be updated with legal moves for the next turn.
// King castling and queen castling are encoded as "O-O" and "O-O-O" respectively.
func Move2SAN(m Move, p *Position, lm *MoveList) string {
	var b strings.Builder
	b.Grow(2)

	moved := p.GetPieceFromSquare(1 << m.From())
	captured := p.GetPieceFromSquare(1 << m.To())

	if m.Type() == MoveCastling {
		if m.To() == SC1 || m.To() == SC8 {
			b.WriteString("O-O-O")
		} else {
			b.WriteString("O-O")
		}
	} else {
		switch moved {
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

		// Resolve the ambiguity if needed.  Skip the pawns since their moves
		// are always ambiguous.
		if moved > PieceBPawn {
			// Find and store all ambiguities.
			ambiguities := make([]int, 0)
			for i := range lm.LastMoveIndex {
				if p.GetPieceFromSquare(1<<lm.Moves[i].From()) == moved &&
					lm.Moves[i].To() == m.To() &&
					lm.Moves[i].From() != m.From() {
					ambiguities = append(ambiguities, lm.Moves[i].From())
				}
			}

			// If there are ambiquities, resolve them.
			if len(ambiguities) != 0 {
				b.WriteString(disambiguate(m.From(), ambiguities))
			}
		}

		if captured != PieceNone || m.Type() == MoveEnPassant {
			if moved <= PieceBPawn {
				b.WriteByte(Square2String[m.From()][0])
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
	}

	p.MakeMove(m, moved, captured)

	GenLegalMoves(*p, lm)

	// Clear en passant target square after each completed move.
	ep := 0
	// If the en passant capture is not possible, clear the en passant target,
	// since it can break the threefold-repetition detection by corrupting the
	// Zobrish hash.  See [IsThreefoldRepetition] function commentary.
	for i := range lm.LastMoveIndex {
		if lm.Moves[i].Type() == MoveEnPassant {
			ep = p.EPTarget
			break
		}
	}
	p.EPTarget = ep

	// The move is check if the opponent's king is under attack.
	isCheck := genAttacks(p.Bitboards, 1^p.ActiveColor)&
		p.Bitboards[PieceWKing+(p.ActiveColor)] != 0

	if isCheck && lm.LastMoveIndex == 0 {
		// If the move results in checkmate, append the '#' symbol to the SAN.
		b.WriteByte('#')
	} else if isCheck {
		// If the move results in check, append the '+' symbol to the SAN.
		b.WriteByte('+')
	}

	return b.String()
}

// disambiguate resolves the ambiguity that arises when multiple pieces of the same
// type can move to the same square.
//
// Steps to resolve:
//  1. If the moving pieces can be distinguished by their originating files, the
//     originating file letter of the moving piece is inserted immediately after
//     the moving piece letter;
//  2. If the moving pieces can be distinguished by their originating ranks, the
//     originating rank digit of the moving piece is inserted immediately after the
//     moving piece letter;
//  3. When both the first and second steps fail, the file letter and rank digit of
//     the moving piece are inserted immediately after the piece letter.
func disambiguate(from int, ambiguities []int) string {
	ranksDiff := 0
	filesDiff := 0

	for i := range ambiguities {
		if ambiguities[i]%8 != from%8 {
			filesDiff++
		}
		if ambiguities[i]/8 != from/8 {
			ranksDiff++
		}
	}

	// Step 1.
	if filesDiff == len(ambiguities) {
		return string(Square2String[from][0])
	}
	// Step 2.
	if ranksDiff == len(ambiguities) {
		return string(Square2String[from][1])
	}
	// Step 3.
	return Square2String[from]
}
