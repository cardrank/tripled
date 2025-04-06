package tripled

import (
	"math"
	"math/rand"
)

// DefaultShuffles is the default shuffles.
var DefaultShuffles = 3

// DefaultPayout is the default payout.
var DefaultPayout = 0.97

// DefaultDist is the default dist.
var DefaultDist = NewDist(DefaultPayout)

// Dist is a dist.
type Dist struct {
	strip [9][][]int
}

// NewDist creates a dist.
func NewDist(f float64) *Dist {
	if f <= 0.0 || 1.0 <= f {
		panic("invalid f")
	}
	d := new(Dist)
	for n := range 9 {
		var payout int
		var w, l [][]int
		for i := range 22 {
			for j := range 22 {
				for k := range 22 {
					res := NewResult(n+1, i, j, k)
					if res.Payout != 0 {
						w = append(w, res.Pos)
						payout += res.Payout // delta
					} else {
						l = append(l, res.Pos)
					}
				}
			}
		}
		// fill
		d.strip[n] = make([][]int, int(math.Ceil(float64(payout)/float64(n+1)/f)))
		copy(d.strip[n], w)
		for i := len(w); i < len(d.strip[n]); i += len(l) {
			copy(d.strip[n][i:], l)
		}
		// shuffle
		r := rand.New(rand.NewSource(1788975))
		for range DefaultShuffles {
			r.Shuffle(len(d.strip[n]), func(i, j int) {
				d.strip[n][i], d.strip[n][j] = d.strip[n][j], d.strip[n][i]
			})
		}
	}
	return d
}

// Spin spins the reels.
func (d *Dist) Spin(r Rand, lines int) (Result, error) {
	// validate
	if lines < 1 || len(Lines) < lines {
		return Result{}, ErrInvalidLines
	}
	return NewResult(lines, d.strip[lines-1][r.Intn(len(d.strip[lines-1]))]...), nil
}

// PayoutKey determines the payout key for the payout.
func PayoutKey(payout int) int {
	switch {
	case payout >= 900:
		return 900
	case payout >= 300:
		return 300
	case payout >= 100:
		return 100
	case payout >= 40:
		return 40
	case payout >= 20:
		return 20
	case payout >= 10:
		return 10
	case payout >= 2:
		return 2
	}
	return 0
}
