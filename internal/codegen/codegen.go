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
	// sanEx expression is needed to extract clean SANs from the dirty
	// PGN movetext.
	sanEx = regexp.MustCompile(`([NBRQK]?[a-h]?[1-8]?x?[a-h][1-8](=[NBRQ])?[+#]?)|(O-O(-O)?[+#]?)`)
	// Annotations sometimes occur in movetext and must be trimmed.
	annotationEx = regexp.MustCompile(`[!?]{1,2}`)
)

// workerPool manages the execution of a set of jobs by concurrent workers.
type workerPool struct {
	sync.Mutex
	wg      sync.WaitGroup
	results [218]int
	jobs    chan string
}

func newWorkerPool() *workerPool {
	return &workerPool{
		wg:   sync.WaitGroup{},
		jobs: make(chan string),
	}
}

// processGame processes a single movetext at a time by sequentially parsing and
// applying the specified moves and counting which index of strictly legal move
// list was actually played.
func (p *workerPool) processGame() {
	for {
		movetext, ok := <-p.jobs
		if !ok {
			p.wg.Done()
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
					p.Lock()
					p.results[i]++
					hasMatched = true
					p.Unlock()
					break
				}
			}
			if !hasMatched {
				fmt.Printf("no match: %s in movetext %s\n", token, movetext)
			}
		}
	}
}

// clean reads from the specified reader line by line and extracts valid SAN
// move encodings into the writer.  Each output line will contain a sequence of
// SAN moves, separated by a single whitespace. This allows each game to be
// analyzed quickly and independently.
//
// Note: SAN moves appearing inside comments are also recognized as valid moves.
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

// generate generates Huffman codes for indices of legal moves in a MoveList of
// strictly legal moves.  Based on data from the provided reader.  The input
// data must be clean and formed by the [clean] function.
func generate(r *bufio.Reader, output *os.File, workers int) {
	numGames := 0
	p := newWorkerPool()

	for range workers {
		go p.processGame()
		p.wg.Add(1)
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
			p.jobs <- chunk[:len(chunk)-2]
			numGames++

			if err == io.EOF {
				break
			}
		}
		close(p.jobs)
	}()

	p.wg.Wait()

	codes := huffmanTree(&p.results)

	numMoves := 0
	for i := range 218 {
		fmt.Fprintf(output, "{0b%s, %d}\t\t\t// index %d | played %d times\n",
			codes[i], len(codes[i]), i, p.results[i])
		numMoves += p.results[i]
	}
	fmt.Fprintf(output, "%d games analyzed\n", numGames)
	fmt.Fprintf(output, "%d moves in tree\n", numMoves)
}

// huffmanTree sorts the input array and builds the Huffman coding tree.
func huffmanTree(results *[218]int) (codes [218]string) {
	sorted := make([]*chego.Node, 0)
	for i := range results {
		// Manually assign 1 frequency to moves that haven't been played to
		// build a valid Huffman tree.
		if results[i] == 0 {
			results[i]++
		}

		n := chego.NewNode(nil, nil, i, results[i])

		wasAppended := false
		for j := range sorted {
			if sorted[j].Freq < results[i] {
				sorted = append(sorted[:j], append(
					[]*chego.Node{n},
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

		n := chego.NewNode(left, right, -1, left.Freq+right.Freq)

		wasAppended := false
		for i := range sorted {
			if sorted[i].Freq < n.Freq {
				sorted = append(sorted[:i], append(
					[]*chego.Node{n},
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
	chego.TraversePreOrder(sorted[0], &codes, "")
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
