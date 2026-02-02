package chego

var (
	// Each piece weight used to calculate material on the board.
	// Use Piece type as index to get it's weight.
	pieceWeights = [10]int{1, 1, 3, 3, 3, 3, 5, 5, 9, 9}
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
)

// Position represents a chessboard state that can be converted to or parsed from
// the FEN string.
type Position struct {
	Bitboards      [15]uint64
	ActiveColor    Color
	CastlingRights CastlingRights
	EPTarget       int
	HalfmoveCnt    int
	FullmoveCnt    int
}

// MakeMove modifies the position by applying the specified move.  It is the
// caller’s responsibility to ensure that the specified move is at least
// pseudo-legal.
//
// Not only is the piece placement updated, but also the entire position, including
// castling rights, en passant target, halfmove counter, fullmove counter, and the
// active color.
func (p *Position) MakeMove(m Move, moved, captured Piece) {
	to := uint64(1 << m.To())
	from := uint64(1 << m.From())

	// Clear the origin square.
	p.removePiece(moved, from)

	// Increment halfmove counter to detect 50-move rule draw.
	// This will be reset if the move is a capture or a pawn push.
	p.HalfmoveCnt++

	// Remove the captured piece from the board.
	// This skips en passant captures, since the captured
	// pawn does not occupy the square the capturing piece moves to.
	if captured != PieceNone {
		p.removePiece(captured, to)
		// Reset the halfmove counter after capture.
		p.HalfmoveCnt = 0
	}

	switch m.Type() {
	case MoveNormal:
		p.placePiece(moved, to)

	case MoveEnPassant:
		p.placePiece(moved, to)
		// Remove the captured piece from the board.
		if moved == WPawn {
			p.removePiece(BPawn, to>>8)
		} else {
			p.removePiece(WPawn, to<<8)
		}

	case MoveCastling:
		p.placePiece(moved, to)
		// Update the rook position.
		switch to {
		case G1: // White O-O.
			p.removePiece(WRook, H1)
			p.placePiece(WRook, F1)
		case G8: // Black O-O.
			p.removePiece(BRook, H8)
			p.placePiece(BRook, F8)
		case C1: // White O-O-O.
			p.removePiece(WRook, A1)
			p.placePiece(WRook, D1)
		case C8: // Black O-O-O.
			p.removePiece(BRook, A8)
			p.placePiece(BRook, D8)
		}

	case MovePromotion:
		switch m.PromoPiece() {
		case PromotionKnight:
			p.placePiece(WKnight+p.ActiveColor, to)
		case PromotionBishop:
			p.placePiece(WBishop+p.ActiveColor, to)
		case PromotionRook:
			p.placePiece(WRook+p.ActiveColor, to)
		case PromotionQueen:
			p.placePiece(WQueen+p.ActiveColor, to)
		}
	}

	// Reset the en passant target since the en passant capture
	// is only legal for 1 move.
	p.EPTarget = 0

	switch moved {
	// Set en passant target square in case of double pawn push.
	case WPawn, BPawn:
		if m.To()+16 == m.From() {
			p.EPTarget = m.To() + 8
		} else if m.To()-16 == m.From() {
			p.EPTarget = m.To() - 8
		}
		// Reset the halfmove counter after pawn moves.
		p.HalfmoveCnt = 0
	// The king cannot castle with a rook that has already moved.
	case WRook:
		switch m.From() {
		case SA1:
			p.CastlingRights &= ^CastlingWhiteLong
		case SH1:
			p.CastlingRights &= ^CastlingWhiteShort
		}
	// The king cannot castle with a rook that has already moved.
	case BRook:
		switch m.From() {
		case SA8:
			p.CastlingRights &= ^CastlingBlackLong
		case SH8:
			p.CastlingRights &= ^CastlingBlackShort
		}
	// Disable white castling rights.
	case WKing:
		p.CastlingRights &= ^(CastlingWhiteShort | CastlingWhiteLong)
	// Disable black castling rights.
	case BKing:
		p.CastlingRights &= ^(CastlingBlackShort | CastlingBlackLong)
	}

	// Increment the full move counter after black moves.
	if p.ActiveColor == ColorBlack {
		p.FullmoveCnt++
	}

	// Switch the active color.
	p.ActiveColor ^= 1
}

// GetPieceFromSquare returns the type of the piece that stands on the specified
// square, or [PieceNone] if the square is empty.
func (p *Position) GetPieceFromSquare(square uint64) Piece {
	for i := range p.Bitboards {
		if square&p.Bitboards[i] != 0 {
			return i
		}
	}
	return PieceNone
}

// canCastle checks whether the king can peform castling in the specified direction.
//
// side represents a castling type:
//   - 1 -> White O-O.
//   - 2 -> White O-O-O.
//   - 4 -> Black O-O.
//   - 8 -> Black O-O-O.
func (p *Position) canCastle(side int, attacks, occupancy uint64) bool {
	c := bitScan(uint64(side))
	path := castlingPath[c]
	return p.CastlingRights&side != 0 &&
		attacks&castlingAttackPath[c] == 0 &&
		occupancy&path == 0
}

// placePiece places the piece on the specified square as well as updates the
// occupancy and allies bitboards.
func (p *Position) placePiece(piece Piece, square uint64) {
	// Place the piece.
	p.Bitboards[piece] |= square
	// Update allies bitboard.
	p.Bitboards[12+(piece%2)] |= square
	// Update occupancy bitboard.
	p.Bitboards[14] |= square
}

// removePiece removes the piece from the specified square as well as updates the
// occupancy and allies bitboards.
//
// NOTE: If a piece of the specified type is not present on the specified square,
// it will be placed rather than removed.
func (p *Position) removePiece(piece Piece, square uint64) {
	// Remove the piece.
	p.Bitboards[piece] ^= square
	// Update allies bitboard.
	p.Bitboards[12+(piece%2)] ^= square
	// Update occupancy bitboard.
	p.Bitboards[14] ^= square
}

// calculateMaterial calculates the piece valies of each side.  Used to determine
// a draw by insufficient material.
func (p *Position) calculateMaterial() (material int) {
	for piece := range WKing {
		material += CountBits(p.Bitboards[piece]) * pieceWeights[piece]
	}
	return material
}

// zobristKey hashes the position into a 64-bit unsigned integer.   This allows
// positions to be used as lookup keys and stored or compared efficiently.
func (p Position) zobristKey() (key uint64) {
	for i := WPawn; i <= BKing; i++ {
		for p.Bitboards[i] > 0 {
			key ^= pieceKeys[i][popLSB(&p.Bitboards[i])]
		}
	}

	key ^= epKeys[p.EPTarget]

	key ^= castlingKeys[p.CastlingRights]

	key ^= colorKey & uint64(p.ActiveColor)

	return key
}
