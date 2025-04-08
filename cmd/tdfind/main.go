// Command tdfind finds triple diamond seeds. Useful for finding example seeds.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/cardrank/tripled"
)

func main() {
	positions := flag.String("pos", "21,21,21", "positions")
	n := flag.Int64("n", 2_000_000, "iterations")
	count := flag.Int("count", 1, "count")
	flag.Parse()
	if err := run(*positions, *n, *count); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(positions string, n int64, count int) error {
	pos, err := parsePositions(positions)
	if err != nil {
		return err
	}
	want := count
	for i := int64(0); i < n && count > 0; i++ {
		r := rand.New(rand.NewSource(i))
		res, err := tripled.DefaultDist.Spin(r, 9)
		if err != nil {
			return err
		}
		if slices.Equal(pos, res.Pos) {
			fmt.Fprintf(os.Stdout, "%2d: %v: % 10d\n", want-count+1, res.Pos, i)
			symbols := lineRE.ReplaceAllString(tripled.SymbolsString(pos...), strings.Repeat(" ", 13))
			fmt.Fprintln(os.Stdout, symbols)
			count--
		}
	}
	return nil
}

var lineRE = regexp.MustCompile(`(?im)^`)

func parsePositions(s string) ([]int, error) {
	v := strings.SplitN(s, ",", 3)
	if len(v) != 3 {
		return nil, fmt.Errorf("invalid positions %q", s)
	}
	pos := make([]int, 3)
	for i := range 3 {
		j, err := strconv.ParseInt(v[i], 10, 64)
		switch {
		case err != nil:
			return nil, fmt.Errorf("invalid position %q: %w", v[i], err)
		case j < 0, 21 < j:
			return nil, fmt.Errorf("invalid position %q", v[i])
		}
		pos[i] = int(j)
	}
	return pos, nil
}
