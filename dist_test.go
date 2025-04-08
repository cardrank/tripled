package tripled

import (
	"maps"
	"math/rand"
	"slices"
	"strconv"
	"testing"
	"time"
)

func TestDist(t *testing.T) {
	const bet = 1
	const iters = 1_000_000
	for n := range 9 {
		lines := n + 1
		t.Run(strconv.Itoa(lines), func(t *testing.T) {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			wagred, payout, wins := 0, 0, 0
			for range iters {
				res, err := DefaultDist.Spin(r, lines)
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				wagred += bet * lines
				payout += bet * res.Payout
				if res.Payout != 0 {
					wins++
				}
			}
			t.Logf(
				"wins: %d/%d (%0.2f%%) payout: %d/%d (%0.2f%%)",
				wins, iters,
				100.0*float64(wins)/float64(iters),
				payout, wagred,
				100.0*float64(payout)/float64(wagred),
			)
		})
	}
}

func TestDistPulls(t *testing.T) {
	const bet = 100
	const starting = 10000
	for n := range 9 {
		lines := n + 1
		t.Run(strconv.Itoa(lines), func(t *testing.T) {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			running, pulls := starting, 0
			m := make(map[int]int)
			for running > 0 && running >= bet*lines {
				res, err := DefaultDist.Spin(r, lines)
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				running += -bet*lines + bet*res.Payout
				m[res.Payout]++
				pulls++
			}
			t.Logf("pulls: %d", pulls)
			keys := slices.Sorted(maps.Keys(m))
			slices.Reverse(keys)
			for _, k := range keys {
				t.Logf("% 5d: %d", k, m[k])
			}
		})
	}
}
