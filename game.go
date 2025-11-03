/*
game.go impements chess game state management.
*/

package chego

/*
Game represents a game state that can be converted to or parsed from the PGN
string.

It's the user's responsibility to spin up a time.Ticker and handle time ticks
by calling the [DecrementTime] function.  The value of timeBonus is added to
the player's timer during the [PushMove] function, so the user must ensure that
time ticks and moves are not handled concurrently (prevent race conditions).

NOTE: Call [InitAttackTables] and [InitZobristKeys] ONCE before creating a
[Game].
*/
type Game struct {
	LegalMoves MoveList
	position   Position
	// Repetition keys are stored as a map of Zobrist keys to the number of
	// times each position has occurred.
	repetitions map[uint64]int
	Result      Result
	Termination Termination
	whiteTime   int
	blackTime   int
	timeBonus   int
}

func NewGame() *Game {
	g := &Game{
		position:    ParseFEN(InitialPos),
		repetitions: make(map[uint64]int, 1),
		Result:      ResultUnknown,
		Termination: TerminationUnterminated,
	}

	GenLegalMoves(g.position, &g.LegalMoves)

	// Initialize Zobrist key for the initial position.
	g.repetitions[g.position.zobristKey()] = 1

	return g
}

/*
PushMove updates the game state by performing the specified move and returns its
Standard Algebraic Notation.  It's a caller's responsibility to ensure that the
specified move is legal.  Not safe for concurrent use.
*/
func (g *Game) PushMove(m Move) string {
	moved := g.position.GetPieceFromSquare(1 << m.From())
	captured := g.position.GetPieceFromSquare(1 << m.To())
	isCapture := captured != PieceNone

	// Encode the move in the Standard Algebraic Notation.  Note that the check
	// and checkmate sybmols must be added later.
	// Move2SAN also perform the move and generates legal moves for next turn.
	san := Move2SAN(m, &g.position, &g.LegalMoves)

	// Clear the repetitions map after applying the irreversable move.
	// See https://www.chessprogramming.org/Irreversible_Moves
	if isCapture || m.Type() == MoveCastling || m.Type() == MovePromotion ||
		moved <= PieceBPawn {
		clear(g.repetitions)
	}

	// Increment the repitition key entry.
	// TODO: optimize by updating the hash incrementally.
	g.repetitions[g.position.zobristKey()]++

	return san
}

/*
IsThreefoldRepetition checks whether the game has reached a threefold repetition.

Two positions are considered identical if all of the following conditions are met:
  - Active colors are the same.
  - Pieces occupy the same squares.
  - Legal moves are the same.
  - Castling rights are identical.

NOTE: Positions are identical even if the en passant target square differs,
provided that no en passant capture is possible.
*/
func (g *Game) IsThreefoldRepetition() bool {
	for _, numOfReps := range g.repetitions {
		if numOfReps >= 3 {
			return true
		}
	}
	return false
}

/*
IsInsufficientMaterial returns true if one of the following statements is true:
  - Both sides have a bare king.
  - One side has a king and a minor piece against a bare king.
  - Both sides have a king and a bishop, the bishops standing on the same color.
  - Both sides have a king and a knight.
*/
func (g *Game) IsInsufficientMaterial() bool {
	// Bitmask for all dark squares.
	dark := uint64(0xAA55AA55AA55AA55)
	material := g.position.calculateMaterial()

	if material == 0 || (material == 3 && g.position.Bitboards[PieceWPawn] == 0 &&
		g.position.Bitboards[PieceBPawn] == 0) {
		return true
	}

	if material == 6 {
		wb := g.position.Bitboards[PieceWBishop]
		bb := g.position.Bitboards[PieceBBishop]

		// If there are two bishops both standing on the same colored squares.
		return (wb != 0 && bb != 0 && ((wb&dark > 0 && bb&dark > 0) ||
			(wb&dark == 0 && bb&dark == 0))) ||
			// Or if there are two knights.
			(g.position.Bitboards[PieceWKnight] != 0 &&
				g.position.Bitboards[PieceBKnight] != 0)
	}
	return false
}

/*
IsCheckmate returns true if both of the following statements are true:
  - There are no legal moves available for the current turn.
  - The king of the side to move is in check.

NOTE: If there are no legal moves, but the king is not in check, the position is
a stalemate.
*/
func (g *Game) IsCheckmate() bool {
	return GenChecksCounter(g.position.Bitboards, 1^g.position.ActiveColor) > 0 &&
		g.LegalMoves.LastMoveIndex == 0
}

/*
IsMoveLegal checks if the specified move is legal by comparing it with moves,
stored in the LegalMoves field.
*/
func (g *Game) IsMoveLegal(m Move) bool {
	for i := range g.LegalMoves.LastMoveIndex {
		lm := g.LegalMoves.Moves[i]
		if lm.From() == m.From() && lm.To() == m.To() && lm.Type() == m.Type() &&
			lm.PromoPiece() == m.PromoPiece() {
			return true
		}
	}
	return false
}

/*
SetClock sets the players’ remaining time and increment (bonus) values.  It
expects these values to be specified in seconds.
*/
func (g *Game) SetClock(control, bonus int) {
	g.whiteTime = control
	g.blackTime = control
	g.timeBonus = bonus
}
