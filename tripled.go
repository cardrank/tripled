package tripled

import (
	"fmt"
	"io"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"
)

// Rand is the shared random interface.
type Rand interface {
	Intn(n int) int
}

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
	return NewResult(pos, lines), nil
}

// Result is a spin result.
type Result struct {
	Pos    []int       `json:"pos"`
	Lines  map[int]int `json:"lines"`
	Payout int         `json:"payout"`
}

// NewResult creates a result with the specified reel positions.
func NewResult(pos []int, lines int) Result {
	res := Result{
		Pos:   pos,
		Lines: make(map[int]int),
	}
	// determine payout
	symbols := Symbols(res.Pos)
	for i := range lines {
		if d := Payout(symbols, Lines[i]); d != 0 {
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
	var sb strings.Builder
	for i, s := range Symbols(res.Pos) {
		if i != 0 && i%3 == 0 {
			_ = sb.WriteByte('\n')
		}
		_, _ = sb.WriteRune(s.Rune())
	}
	return sb.String()
}

// Symbol is a slot symbol.
type Symbol int

// Symbols.
const (
	Blank Symbol = iota
	Bar1
	Bar2
	Bar3
	Seven
	Diamond
)

// Format satisfies the [fmt.Formatter] interface.
func (s Symbol) Format(f fmt.State, verb rune) {
	var buf []byte
	switch verb {
	case 's', 'v':
		buf = append(buf, s.Name()...)
	case 'c':
		buf = append(buf, string(s.Rune())...)
	case 'd':
		buf = append(buf, strconv.Itoa(int(s))...)
	default:
		buf = append(buf, fmt.Sprintf("%%!%c(ERROR=unknown verb, symbol: %d)", verb, int(s))...)
	}
	_, _ = f.Write(buf)
}

// Name returns the symbol name.
func (s Symbol) Name() string {
	switch s {
	case Bar1:
		return "Bar1"
	case Bar2:
		return "Bar2"
	case Bar3:
		return "Bar3"
	case Seven:
		return "Seven"
	case Diamond:
		return "Diamond"
	}
	return "Blank"
}

// Rune returns the symbol rune.
func (s Symbol) Rune() rune {
	switch s {
	case Bar1:
		return '-'
	case Bar2:
		return '='
	case Bar3:
		return '≡'
	case Seven:
		return '7'
	case Diamond:
		return '◆'
	}
	return '.'
}

// Error is an error.
type Error string

// Error satisfies the [error] interface.
func (err Error) Error() string {
	return string(err)
}

// Errors.
const (
	// ErrInvalidLines is the invalid lines error.
	ErrInvalidLines Error = "invalid lines"
)

// Symbols returns the symbols in the positions.
func Symbols(pos []int) []Symbol {
	v := make([]Symbol, 9)
	for j := range 3 {
		for i := range 3 {
			n := len(Reels[i])
			v[i+j*3] = Reels[i][(((pos[i]+j-1)%n)+n)%n]
		}
	}
	return v
}

// Payout determines the payout for the mask in symbols.
func Payout(symbols []Symbol, mask int) int {
	d, s, b3, b2, b1, n := 0, 0, 0, 0, 0, 0
	for i := range 9 {
		if symbols[i] == Blank || mask&(1<<i) == 0 {
			continue
		}
		switch symbols[i] {
		case Diamond:
			d++
		case Seven:
			s++
		case Bar3:
			b3++
		case Bar2:
			b2++
		case Bar1:
			b1++
		}
		n++
	}
	switch {
	case d == 3: // 3x diamond
		return 1199
	case n == 3:
		// diamond multiplier
		switch mlt := int(math.Pow(3, float64(d))); n - d {
		case s:
			return mlt * 100 // 3x Seven -- 20x
		case b3:
			return mlt * 40 // 3x Bar3 -- 8x
		case b2:
			return mlt * 20 // 3x Bar2 -- 4x
		case b1:
			return mlt * 10 // 3x Bar1 -- 2x
		case b3 + b2 + b1:
			return mlt * 5 // Any Bar -- 5x
		}
		// 1 diamond + 2 other
		return 2
	case d == 2: // 2x diamond
		return 10
	case d == 1: // 1x diamond
		return 2
	}
	return 0
}

// Reels are the slot reels.
var Reels = [3][22]Symbol{
	{
		Blank,
		Bar1,
		Blank,
		Bar2,
		Blank,
		Bar1,
		Blank,
		Bar3, // 8
		Blank,
		Bar1,
		Blank,
		Bar2,
		Blank,
		Bar1,
		Blank,
		Seven,
		Blank,
		Bar1,
		Blank,
		Bar1,
		Blank,
		Diamond,
	},
	{
		Blank,
		Blank,
		Bar2,
		Blank,
		Bar1,
		Blank,
		Bar2,
		Blank,
		Bar3,
		Blank,
		Seven,
		Blank,
		Bar2,
		Blank,
		Bar2,
		Blank,
		Bar3,
		Blank,
		Bar2,
		Blank,
		Blank,
		Diamond,
	},
	{
		Blank,
		Blank,
		Bar1,
		Blank,
		Bar2,
		Blank,
		Bar1,
		Blank,
		Bar3,
		Blank,
		Seven,
		Blank,
		Bar1,
		Blank,
		Bar1,
		Blank,
		Bar2,
		Blank,
		Bar1,
		Blank,
		Blank,
		Diamond,
	},
}

// Lines are the payout lines.
var Lines = [9]int{
	0: CW | CC | CE,
	1: NW | NC | NE,
	2: SW | SC | SE,
	3: NW | CC | SE,
	4: SW | CC | NE,
	5: CW | NC | CE,
	6: CW | SC | CE,
	7: SW | CC | SE,
	8: NW | CC | NE,
}

// Coordinates.
const (
	NW = 1 << iota
	NC
	NE
	CW
	CC
	CE
	SW
	SC
	SE
)
