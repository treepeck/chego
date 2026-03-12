package chego

import "testing"

func TestMove2SAN(t *testing.T) {
	cases := []struct {
		move     Move
		pos      *Position
		expected string
	}{
		{
			NewMove(SE2, SC3, MoveNormal),
			ParseFen("8/8/8/8/8/2N5/8/4K1N1 w - - 0 1"),
			"Nce2",
		},
		// Similar case to the previous one, except the knight c3 is pinned by
		// the black bishop, so the disambiguation is not needed.
		{
			NewMove(SE2, SG1, MoveNormal),
			ParseFen("8/8/8/8/1b6/2N5/8/4K1N1 w - - 0 1"),
			"Ne2",
		},
		{
			NewMove(SB7, SA6, MoveNormal),
			ParseFen("2k5/Qr6/Q7/8/8/8/8/3R4 w - - 0 1"),
			"Q6xb7#",
		},
		{
			NewPromotionMove(SE8, SD7, PromotionQueen),
			ParseFen("4b3/3P1P2/8/8/8/8/8/8 w - - 0 1"),
			"dxe8=Q",
		},
		{
			NewMove(SE4, SF6, MoveNormal),
			ParseFen("rnbqkb1r/pppppppp/5n2/8/3PP3/8/PPP2PPP/RNBQKBNR b KQkq - 0 1"),
			"Nxe4",
		},
		{
			NewMove(SD4, SE5, MoveNormal),
			ParseFen("8/8/8/4p3/3P4/2K5/8/8 b - - 0 1"),
			"exd4+",
		},
		{
			NewMove(SE7, SF7, MoveNormal),
			ParseFen("r1bk3r/ppqpbQpp/2p4n/6B1/2BpP3/3P1P2/PPP3PP/RN3RK1 w - - 0 1"),
			"Qxe7#",
		},
		{
			NewMove(SB8, SE5, MoveNormal),
			ParseFen("Q3Q2Q/8/8/4Q3/4P3/2N5/3k2P1/R5K1 w - - 0 1"),
			"Q5b8",
		},
	}

	for _, tc := range cases {
		var legal MoveList
		GenLegalMoves(*tc.pos, &legal)

		got := Move2SAN(tc.move, tc.pos, &legal)
		if got != tc.expected {
			t.Fatalf("expected: %s, got: %s", tc.expected, got)
		}
	}
}

func BenchmarkMove2SAN(b *testing.B) {
	p := ParseFen("r1bk3r/ppqpbQpp/2p4n/6B1/2BpP3/3P1P2/PPP3PP/RN3RK1 w - - 0 1")
	var legal MoveList
	GenLegalMoves(*p, &legal)

	for b.Loop() {
		prev := *p
		Move2SAN(NewMove(SE7, SF7, MoveNormal), p, &legal)
		p = &prev
	}
}
