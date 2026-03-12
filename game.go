// game.go impements chess game management.

package chego

// PlayedMove represents the played chess move.
type PlayedMove struct {
	// Standard Algebraic Notation.
	San string `json:"s"`
	// Forsyth-Edwards Notation. Describes the [Position] after played move.
	Fen string `json:"f"`
}

// Game represents the state of a chess game.
//
// Methods are not safe for concurrent use.
type Game struct {
	Played   []PlayedMove `json:"p"`
	Legal    *MoveList    `json:"l"`
	Position *Position    `json:"-"`
	// Repetition keys are stored as a map of Zobrist keys to the number of
	// times each Position has occurred.
	repetitions map[uint64]int
	Result      Result      `json:"r"`
	Termination Termination `json:"t"`
}

// NewGame initializes [Game] fields and generates legal moves on the [InitialPos].
func NewGame() *Game {
	g := &Game{
		Played:      make([]PlayedMove, 0),
		Position:    ParseFen(InitialPos),
		Legal:       &MoveList{},
		repetitions: make(map[uint64]int, 1),
		Result:      Unknown,
		Termination: Unterminated,
	}
	GenLegalMoves(*g.Position, g.Legal)
	// Initialize Zobrist key for the initial position.
	g.repetitions[g.Position.zobristKey()] = 1
	return g
}

// Push updates the [Game] by performing [Move] with specified index in legal
// [MoveList]. It's the caller's responsibility to validate m.
//
// Game [Result] and [Termination] will not be modified.
func (g *Game) Push(m Move) {
	moved := g.Position.GetPieceFromSquare(1 << m.From())
	captured := g.Position.GetPieceFromSquare(1 << m.To())
	isCapture := captured != PieceNone

	// Move2SAN updates the position and generates legal moves for next turn.
	san := Move2SAN(m, g.Position, g.Legal)

	// Clear the repetitions map after applying the irreversable move.
	// See https://www.chessprogramming.org/Irreversible_Moves
	if isCapture || m.Type() == MoveCastling || m.Type() == MovePromotion ||
		moved <= BPawn {
		clear(g.repetitions)
	}

	// Increment the repitition key entry.
	g.repetitions[g.Position.zobristKey()]++

	// Store played move.
	g.Played = append(g.Played, PlayedMove{
		San: san,
		Fen: SerializeFen(g.Position),
	})
}

// Pop discards the last pushed move and restores the [Game] state.
func (g *Game) Pop() {
	if len(g.Played) == 0 {
		return
	}

	// Decrement the repetition key entry.
	g.repetitions[g.Position.zobristKey()]--

	m := g.Played[len(g.Played)-1]
	g.Played = g.Played[:len(g.Played)-1]

	g.Position = ParseFen(m.Fen)

	GenLegalMoves(*g.Position, g.Legal)
}

// IsThreefoldRepetition checks whether the game has reached a threefold repetition.
//
// Two positions are considered identical if all of the following conditions are met:
//   - Active colors are the same.
//   - Pieces occupy the same squares.
//   - Legal moves are the same.
//   - Castling rights are identical.
//
// NOTE: Positions are identical even if the en passant target square differs,
// provided that no en passant capture is possible.
func (g *Game) IsThreefoldRepetition() bool {
	for _, numOfReps := range g.repetitions {
		if numOfReps >= 3 {
			return true
		}
	}
	return false
}

// IsInsufficientMaterial returns true if one of the following statements is true:
//   - Both sides have a bare king.
//   - One side has a king and a minor piece against a bare king.
//   - Both sides have a king and a bishop, the bishops standing on the same color.
//   - Both sides have a king and a knight.
func (g *Game) IsInsufficientMaterial() bool {
	// Bitmask for all dark squares.
	dark := uint64(0xAA55AA55AA55AA55)
	material := g.Position.calculateMaterial()
	if material == 0 || (material == 3 && g.Position.Bitboards[WPawn] == 0 &&
		g.Position.Bitboards[BPawn] == 0) {
		return true
	}

	if material == 6 {
		wb := g.Position.Bitboards[WBishop]
		bb := g.Position.Bitboards[BBishop]

		// If there are two bishops both standing on the same colored squares.
		return (wb != 0 && bb != 0 && ((wb&dark > 0 && bb&dark > 0) ||
			(wb&dark == 0 && bb&dark == 0))) ||
			// Or if there are two knights.
			(g.Position.Bitboards[WKnight] != 0 &&
				g.Position.Bitboards[BKnight] != 0)
	}
	return false
}

// IsCheckmate returns true if both of the following statements are true:
//   - There are no legal moves available for the current turn.
//   - The king of the side to move is in check.
//
// NOTE: If there are no legal moves, but the king is not in check, the position
// is a stalemate.
func (g *Game) IsCheckmate() bool {
	return GenChecksCounter(g.Position.Bitboards, 1^g.Position.ActiveColor) > 0 &&
		g.Legal.LastMoveIndex == 0
}

// IsMoveLegal checks if the specified move is legal by comparing it with moves,
// stored in the Legal field.
func (g *Game) IsMoveLegal(m Move) bool {
	for i := range g.Legal.LastMoveIndex {
		lm := g.Legal.Moves[i]
		if lm.From() == m.From() && lm.To() == m.To() && lm.Type() == m.Type() &&
			lm.PromoPiece() == m.PromoPiece() {
			return true
		}
	}
	return false
}
