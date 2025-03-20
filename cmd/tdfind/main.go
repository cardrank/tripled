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
	lines := flag.Int("lines", 9, "lines")
	flag.Parse()
	if err := run(*positions, *n, *count, *lines); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(positions string, n int64, count, lines int) error {
	pos, err := parsePositions(positions)
	if err != nil {
		return err
	}
	want := count
	for i := int64(0); i < n && count > 0; i++ {
		r := rand.New(rand.NewSource(i))
		v := []int{
			r.Intn(22),
			r.Intn(22),
			r.Intn(22),
		}
		if slices.Equal(pos, v) {
			fmt.Fprintf(os.Stdout, "%2d: %v: % 10d\n", want-count+1, v, i)
			symbols := lineRE.ReplaceAllString(tripled.NewResult(pos, lines).Symbols(), strings.Repeat(" ", 13))
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
