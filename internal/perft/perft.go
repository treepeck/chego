// peft.go implements debugging and testing functions for the move generator.
//
// It is internal, as it is only used for testing purposes.
//
// TODO: fix verbose perft.  It doesn't print the resulting information correctly.

package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/treepeck/chego"
)

// result information is printed to the console when the verbose flag is used.
type result struct {
	nodes        int
	captures     int
	epCaptures   int
	castles      int
	promotions   int
	checks       int
	doubleChecks int
	checkmates   int
}

// perft is a debugging function that walks through the move generation tree of
// strictly legal moves to a given depth and counts the number of visited leaf
// nodes. The resulting count is then compared to predetermined values.
//
// See https://www.chessprogramming.org/Perft_Results
func perft(p chego.Position, depth int) int {
	l := chego.MoveList{}
	nodes := 0

	chego.GenLegalMoves(p, &l)

	if depth == 1 {
		return int(l.LastMoveIndex)
	}

	var prev chego.Position
	var moved, captured chego.Piece

	for i := range l.LastMoveIndex {
		prev = p
		moved = p.GetPieceFromSquare(1 << l.Moves[i].From())
		captured = p.GetPieceFromSquare(1 << l.Moves[i].To())
		p.MakeMove(l.Moves[i], moved, captured)

		nodes += perft(p, depth-1)

		p = prev
	}

	return nodes
}

// perftVerbose follows the same principle as the perft function, except it
// writes detailed move debugging information to r. Use this function to debug
// and find invalid branches in the move generation tree, not to measure
// performance.
func perftVerbose(p chego.Position, depth int, r *result, isRoot bool) int {
	l := chego.MoveList{}
	nodes := 0

	chego.GenLegalMoves(p, &l)

	if depth == 1 {
		return int(l.LastMoveIndex)
	}

	c := p.ActiveColor
	var prev chego.Position
	var moved, captured chego.Piece

	for i := range l.LastMoveIndex {
		if p.GetPieceFromSquare(1<<l.Moves[i].To()) != chego.PieceNone {
			r.captures++
		}

		prev = p
		moved = p.GetPieceFromSquare(1 << l.Moves[i].From())
		captured = p.GetPieceFromSquare(1 << l.Moves[i].To())
		p.MakeMove(l.Moves[i], moved, captured)

		cnt := chego.GenChecksCounter(p.Bitboards, 1^c)
		if cnt > 0 {
			r.checks++
		}
		if cnt > 1 {
			r.doubleChecks++
		}

		cnt = perftVerbose(p, depth-1, r, false)
		if isRoot {
			log.Printf("%s %d", move2UCI(l.Moves[i]), cnt)
		}
		nodes += cnt

		switch l.Moves[i].Type() {
		case chego.MoveCastling:
			r.castles++
		case chego.MoveEnPassant:
			r.epCaptures++
		case chego.MovePromotion:
			r.promotions++
		}

		p = prev
	}

	return nodes
}

// move2UCI converts the move into a long algebraic notation string.
//
// Examples: e2e4, e7e5, e1g1 (white short castling), e7e8q (for promotion).
func move2UCI(m chego.Move) string {
	var b strings.Builder
	b.Grow(4)

	b.WriteString(chego.Square2String[m.From()])
	b.WriteString(chego.Square2String[m.To()])

	if m.Type() == chego.MovePromotion {
		switch m.PromoPiece() {
		case chego.PromotionKnight:
			b.WriteByte('n')
		case chego.PromotionBishop:
			b.WriteByte('b')
		case chego.PromotionRook:
			b.WriteByte('r')
		case chego.PromotionQueen:
			b.WriteByte('q')
		}
	}

	return b.String()
}

// main runs the perft and measures it's execution time.
func main() {
	depth := flag.Int("depth", 1, "Performance test depth")
	verbose := flag.Bool("verbose", false, "Wether to print the debug info")
	cpuprofile := flag.String("cpuprofile", "", "File to write a cpu profile")
	memprofile := flag.String("memprofile", "", "File to write a memory profile")

	flag.Parse()

	r := &result{}

	fen := chego.InitialPos
	p := chego.ParseFEN(fen)

	start := time.Now()
	defer func() {
		elapsed := time.Since(start)

		if *verbose {
			log.Printf("\nRoot position:\n%s\n\n\t%s\n\n", position(p), fen)
			log.Printf("\t%d\t%d\t\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t",
				*depth,
				r.nodes,
				r.captures,
				r.epCaptures,
				r.castles,
				r.promotions,
				r.checks,
				r.doubleChecks,
				r.checkmates,
			)
			log.Printf("Elapsed time: %d ns", elapsed.Nanoseconds())
		} else {

			log.Printf("Nodes reached: %d", r.nodes)
			log.Printf("Elapsed time: %d ns", elapsed.Nanoseconds())
		}
	}()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		defer f.Close()
	}

	if *verbose {
		r.nodes = perftVerbose(p, *depth, r, true)
	} else {
		r.nodes = perft(p, *depth)
	}
}

// position formats a full chess position into a string.
func position(p chego.Position) string {
	var b strings.Builder

	for rank := 7; rank >= 0; rank-- {
		b.WriteByte(byte(rank) + 1 + '0')
		b.WriteString("  ")

		for file := range 8 {
			square := uint64(1 << (8*rank + file))

			symbol := byte('.')

			for i := chego.WPawn; i <= chego.BKing; i++ {
				if square&p.Bitboards[i] != 0 {
					symbol = chego.PieceSymbols[i]
					break
				}
			}

			b.WriteByte(symbol)
			b.WriteString("  ")
		}
		b.WriteByte('\n')
	}

	b.WriteString("   a  b  c  d  e  f  g  h\nActive color: ")

	if p.ActiveColor == chego.ColorWhite {
		b.WriteString("white\nEn passant: ")
	} else {
		b.WriteString("black\nEn passant: ")
	}

	if p.EPTarget == 0 {
		b.WriteString("none\nCastling rights: ")
	} else {
		b.WriteString(chego.Square2String[p.EPTarget])
		b.WriteString("\nCastling rights: ")
	}

	if p.CastlingRights&chego.CastlingWhiteShort != 0 {
		b.WriteByte('K')
	}
	if p.CastlingRights&chego.CastlingWhiteLong != 0 {
		b.WriteByte('Q')
	}
	if p.CastlingRights&chego.CastlingBlackShort != 0 {
		b.WriteByte('k')
	}
	if p.CastlingRights&chego.CastlingBlackLong != 0 {
		b.WriteByte('q')
	}

	return b.String()
}
