// Command tripled is a command-line implementation of Triple Diamond slots.
package main

import (
	"flag"
	"fmt"
	"maps"
	"math/rand"
	"os"
	"slices"
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
	win, won, bet, dist := 0, 0, 0, make(map[int]int)
	for range pulls {
		res, err := tripled.DefaultDist.Spin(r, lines)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "%c\n", res)
		won += res.Payout
		bet += lines
		dist[res.Payout]++
		if res.Payout != 0 {
			win++
		}
	}
	keys := slices.Sorted(maps.Keys(dist))
	slices.Reverse(keys)
	for _, k := range keys {
		fmt.Fprintf(os.Stdout, "%4d: %d\n", k, dist[k])
	}
	fmt.Fprintf(
		os.Stdout,
		"win: %d/%d (%0.2f%%) won: %d/%d (%0.2f%%)\n",
		win, pulls, 100.0*float64(win)/float64(pulls),
		won, bet, 100.0*float64(won)/float64(bet),
	)
	return nil
}
