package tripled

import (
	"fmt"
	"io"
	"maps"
	"slices"
)

// Spin spins the reels, calculating the results for the spin.
func Spin(r Rand, lines int) (Result, error) {
	// validate
	if lines <= 0 || len(Lines) < lines {
		return Result{}, ErrInvalidLines
	}
	// randomize reel positions
	pos := make([]int, len(Reels))
	for i := range len(Reels) {
		pos[i] = r.Intn(len(Reels[i]))
	}
	return NewResult(lines, pos...), nil
}

// Result is a spin result.
type Result struct {
	Pos    []int       `json:"pos"`
	Lines  map[int]int `json:"lines"`
	Payout int         `json:"payout"`
}

// NewResult creates a result with the specified reel positions.
func NewResult(lines int, pos ...int) Result {
	if len(pos) != len(Reels) {
		panic(fmt.Sprintf("must have %d pos", len(Reels)))
	}
	res := Result{
		Pos:   pos,
		Lines: make(map[int]int),
	}
	// determine payout
	symbols := Symbols(res.Pos...)
	for i := range lines {
		if d := Payout(Lines[i], symbols...); d != 0 {
			res.Lines[i] = d
			res.Payout += d
		}
	}
	return res
}

// Format satisfies the [fmt.Formatter] interface.
func (res Result) Format(f fmt.State, verb rune) {
	_, _ = res.WriteTo(f)
}

// WriteTo writes result to the writer.
func (res Result) WriteTo(w io.Writer) (int64, error) {
	fmt.Fprintf(w, "pos: %d %d %d\n", res.Pos[0], res.Pos[1], res.Pos[2])
	fmt.Fprintf(w, "%s\n", res.Symbols())
	if len(res.Lines) > 0 {
		fmt.Fprintln(w, "lines:")
		for _, k := range slices.Sorted(maps.Keys(res.Lines)) {
			fmt.Fprintf(w, "% 2d payouts %dx\n", k+1, res.Lines[k])
		}
	}
	fmt.Fprintf(w, "payout: %dx", res.Payout)
	return 0, nil
}

// Symbols produces a string representing the final view of the result.
func (res Result) Symbols() string {
	return SymbolsString(res.Pos...)
}
