// codegen.go implements Huffman code generation.
//
// It is internal, as it is only used to precalculate Huffman codes.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/treepeck/chego"
)

var (
	// PGN movetext may include move numbers and other tokens that are
	// not useful for Huffman encoding. sanEx is used to extract only
	// Standard Algebraic Notation tokens from the movetext.
	sanEx = regexp.MustCompile(`([NBRQK]?[a-h]?[1-8]?x?[a-h][1-8](=[NBRQ])?[+#]?)|(O-O(-O)?[+#]?)`)
	// PGN movetext may include annotations that need to be detected and removed.
	// sanEx cannot remove annotations because they may be attached directly to
	// the SAN move (e.g., e2??).
	annotationEx = regexp.MustCompile(`[!?]{1,2}`)
)

// clean reads from the specified reader line by line and writes valid SAN tokens
// into the resulting file.  r reader must read a valid PGN database file.
// Each output line will contain a sequence of SAN tokens, separated by a single
// whitespace.  All other PGN data will be trimmed.
//
// Note: SAN tokens appearing inside comments are also recognized as valid moves.
// Ensure that the input PGN file doesn't contain SAN moves within comments.
func clean(r *bufio.Reader, output *os.File) {
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}

		// Tag pairs are separated from the movetext by a single empty line.
		if line == "\n" {
			// Read movetext section.
			var b strings.Builder

			hasMoves := false
			for {
				movetext, err := r.ReadString('\n')
				if err != nil || movetext == "\n" {
					break
				}

				for token := range strings.SplitSeq(movetext, " ") {
					// Trim '??', '!!', '?!', and '!?' annotations.
					san := annotationEx.ReplaceAll([]byte(token),
						[]byte(""))
					if sanEx.Match(san) {
						hasMoves = true
						b.WriteString(string(san))
						b.WriteByte(' ')
					}
				}
			}

			// If the game doesn't contain a single move, skip it.
			if !hasMoves {
				continue
			}

			// Append new line to separate movetexts.
			b.WriteByte('\n')
			if _, err := output.WriteString(b.String()); err != nil {
				panic(err)
			}
		}
	}
}

// generate writes the array of generated Huffman codes to the output file.  r
// reader must read a file produced by the [clean] function.
func generate(r *bufio.Reader, output *os.File, workers int) {
	numGames := 0
	g := newGenerator()

	for range workers {
		go g.processGame()
		g.wg.Add(1)
	}

	go func() {
		for {
			chunk, err := r.ReadString('\n')
			if err != nil && err != io.EOF {
				panic(err)
			}

			if len(chunk) < 4 {
				fmt.Printf("Breaking: '%s'\nThe results will be written to the file %s\n.",
					chunk, output.Name())
				break
			}

			// ReadString returns the chunk with the ' \n' suffix, which needs
			// to be trimmed.
			g.jobs <- chunk[:len(chunk)-2]
			numGames++

			if err == io.EOF {
				break
			}
		}
		close(g.jobs)
	}()

	g.wg.Wait()

	codes := g.encode()

	numMoves := 0
	for i := range 218 {
		fmt.Fprintf(output, "{0b%s, %d}\t\t\t// index %d | played %d times\n",
			codes[i], len(codes[i]), i, g.results[i])
		numMoves += g.results[i]
	}
	fmt.Fprintf(output, "%d games analyzed\n", numGames)
	fmt.Fprintf(output, "%d moves in tree\n", numMoves)
}

// node represents the Huffman tree node.
type node struct {
	left  *node
	right *node
	index int // Index of the legal move (in move list).
	freq  int // Number of played times.
}

func newNode(left, right *node, ind, freq int) *node {
	return &node{left: left, right: right, index: ind, freq: freq}
}

// generator manages concurrent Huffman code generation.  It also protects the
// resulting Huffman code array from concurrent writes using a mutex.
type generator struct {
	sync.Mutex
	wg      sync.WaitGroup
	results [218]int
	jobs    chan string
}

func newGenerator() *generator {
	return &generator{
		wg:   sync.WaitGroup{},
		jobs: make(chan string),
	}
}

// processGame processes a single movetext line by sequentially parsing
// and applying the specified moves, tracking which index in the
// strictly legal move list was actually played.
func (g *generator) processGame() {
	for {
		movetext, ok := <-g.jobs
		if !ok {
			g.wg.Done()
			return
		}

		pos := chego.ParseFEN(chego.InitialPos)
		ml := chego.MoveList{}
		chego.GenLegalMoves(pos, &ml)

		for token := range strings.SplitSeq(movetext, " ") {
			prev := pos
			cml := ml
			hasMatched := false

			for i := range cml.LastMoveIndex {
				pos = prev
				ml = cml

				if chego.Move2SAN(ml.Moves[i], &pos, &ml) == token {
					g.Lock()
					g.results[i]++
					hasMatched = true
					g.Unlock()
					break
				}
			}
			if !hasMatched {
				fmt.Printf("no match: %s in movetext %s\n", token, movetext)
			}
		}
	}
}

// traversePreOrder traverses the subtree in pre-order and builds the array of
// resulting encoded strings.
func (n *node) traversePreOrder(codes *[218]string, current string) {
	if n == nil {
		return
	}

	if n.left == nil && n.right == nil {
		(*codes)[n.index] = current
		return
	}

	n.left.traversePreOrder(codes, current+"1")
	n.right.traversePreOrder(codes, current+"0")
}

// encode generates an array of encoded strings. Each string represents
// the legal move at the corresponding index.
func (g *generator) encode() [218]string {
	// TODO: Slowest sorting code ever.  I was lazy to apply quicksort here.
	// It's bad.
	sorted := make([]*node, 0)
	for i := range g.results {
		// Manually assign 1 frequency to moves that haven't been played to
		// build a valid Huffman tree.
		if g.results[i] == 0 {
			g.results[i]++
		}

		n := newNode(nil, nil, i, g.results[i])

		wasAppended := false
		for j := range sorted {
			if sorted[j].freq < g.results[i] {
				sorted = append(sorted[:j], append(
					[]*node{n},
					sorted[j:]...,
				)...)
				wasAppended = true
				break
			}
		}
		if !wasAppended {
			sorted = append(sorted, n)
		}
	}

	// Build the Huffman coding tree.
	for len(sorted) > 1 {
		// Pop left node.
		left := sorted[len(sorted)-1]
		sorted = sorted[:len(sorted)-1]

		// Pop right node.
		right := sorted[len(sorted)-1]
		sorted = sorted[:len(sorted)-1]

		n := newNode(left, right, -1, left.freq+right.freq)

		wasAppended := false
		for i := range sorted {
			if sorted[i].freq < n.freq {
				sorted = append(sorted[:i], append(
					[]*node{n},
					sorted[i:]...,
				)...)
				wasAppended = true
				break
			}
		}
		if !wasAppended {
			sorted = append(sorted, n)
		}
	}

	// Form codes.
	var codes [218]string
	sorted[0].traversePreOrder(&codes, "")
	return codes
}

func main() {
	input := flag.String("input", "clean.txt", "Path to the input file")
	workers := flag.Int("workers", 1, "Number of concurrent routines which will perform the task")
	task := flag.String("task", "gen", "clean task will prepare the PGN file for code generation, generate task will generate the Huffman codes based of frequency of played moves")
	output := flag.String("output", "output.pgn", "Path to the output file")

	flag.Parse()

	inputFile, err := os.Open(*input)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(*output)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	if *task == "clean" {
		clean(bufio.NewReader(inputFile), outputFile)
	} else {
		generate(bufio.NewReader(inputFile), outputFile, *workers)
	}
}
