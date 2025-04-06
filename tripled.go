package tripled

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Rand is the shared random interface.
type Rand interface {
	Intn(n int) int
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
func Symbols(pos ...int) []Symbol {
	symbols := make([]Symbol, 9)
	for j := range 3 {
		for i := range 3 {
			n := len(Reels[i])
			symbols[i+j*3] = Reels[i][(((pos[i]+j-1)%n)+n)%n]
		}
	}
	return symbols
}

// SymbolsString returns the symbols as a string.
func SymbolsString(pos ...int) string {
	var sb strings.Builder
	for i, s := range Symbols(pos...) {
		if i != 0 && i%3 == 0 {
			_ = sb.WriteByte('\n')
		}
		_, _ = sb.WriteRune(s.Rune())
	}
	return sb.String()
}

// Payout determines the payout for matching line symbols based on the bit
// mask.
func Payout(mask int, line ...Symbol) int {
	d, s, b3, b2, b1, n := 0, 0, 0, 0, 0, 0
	for i := range len(line) {
		if line[i] == Blank || mask&(1<<i) == 0 {
			continue
		}
		switch line[i] {
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
	case d == 3:
		return 1199 // 3x diamond -- 1199x
	case n == 3: // diamond multiplier
		switch mlt := int(math.Pow(3, float64(d))); n - d {
		case s:
			return mlt * 100 // 3x Seven -- 100x
		case b3:
			return mlt * 40 // 3x Bar3 -- 40x
		case b2:
			return mlt * 20 // 3x Bar2 -- 20x
		case b1:
			return mlt * 10 // 3x Bar1 -- 10x
		case b3 + b2 + b1:
			return mlt * 5 // Any Bar -- 5x
		}
	}
	// 2 or more diamonds
	switch d {
	case 2:
		return 10
	case 1:
		return 2
	}
	return 0
}

// Reels are reels of slot symbols.
var Reels = [3][22]Symbol{
	{
		Blank,   // 0
		Bar1,    // 1
		Blank,   // 2
		Bar2,    // 3
		Blank,   // 4
		Bar1,    // 5
		Blank,   // 6
		Bar3,    // 7
		Blank,   // 8
		Bar1,    // 9
		Blank,   // 10
		Bar2,    // 11
		Blank,   // 12
		Bar1,    // 13
		Blank,   // 14
		Seven,   // 15
		Blank,   // 16
		Bar1,    // 17
		Blank,   // 18
		Bar1,    // 19
		Blank,   // 20
		Diamond, // 21
	},
	{
		Blank,   // 0
		Blank,   // 1
		Bar2,    // 2
		Blank,   // 3
		Bar1,    // 4
		Blank,   // 5
		Bar2,    // 6
		Blank,   // 7
		Bar3,    // 8
		Blank,   // 9
		Seven,   // 10
		Blank,   // 11
		Bar2,    // 12
		Blank,   // 13
		Bar2,    // 14
		Blank,   // 15
		Bar3,    // 16
		Blank,   // 17
		Bar2,    // 18
		Blank,   // 19
		Blank,   // 20
		Diamond, // 21
	},
	{
		Blank,   // 0
		Blank,   // 1
		Bar1,    // 2
		Blank,   // 3
		Bar2,    // 4
		Blank,   // 5
		Bar1,    // 6
		Blank,   // 7
		Bar3,    // 8
		Blank,   // 9
		Seven,   // 10
		Blank,   // 11
		Bar1,    // 12
		Blank,   // 13
		Bar1,    // 14
		Blank,   // 15
		Bar2,    // 16
		Blank,   // 17
		Bar1,    // 18
		Blank,   // 19
		Blank,   // 20
		Diamond, // 21
	},
}

// Lines are pay line masks.
var Lines = [9]int{
	0: cw | cc | ce,
	1: nw | nc | ne,
	2: sw | sc | se,
	3: nw | cc | se,
	4: sw | cc | ne,
	5: cw | nc | ce,
	6: cw | sc | ce,
	7: sw | cc | se,
	8: nw | cc | ne,
}

// coordinates.
const (
	nw = 1 << iota
	nc
	ne
	cw
	cc
	ce
	sw
	sc
	se
)
