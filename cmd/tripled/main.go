// Command tripled is a command-line implementation of Triple Diamond slots.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/cardrank/tripled"
)

func main() {
	seed := flag.Int64("seed", 0, "seed")
	lines := flag.Int("lines", 9, "lines")
	pulls := flag.Int("pulls", 1, "pulls")
	flag.Parse()
	if err := run(*seed, *lines, *pulls); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(seed int64, lines, pulls int) error {
	t := seed
	if seed == 0 {
		t = time.Now().UnixNano()
	}
	if seed == 0 {
		fmt.Fprintln(os.Stdout, "seed:", t)
	}
	r := rand.New(rand.NewSource(t))
	won, bet := 0, 0
	for range pulls {
		res, err := tripled.DefaultDist.Spin(r, lines)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "%c\n", res)
		won += res.Payout
		bet += lines
	}
	fmt.Fprintf(os.Stdout, "won: %d/%d (%0.2f%%)\n", won, bet, 100.0*float64(won)/float64(bet))
	return nil
}
