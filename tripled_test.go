package tripled

import (
	"maps"
	"math/rand"
	"slices"
	"strconv"
	"testing"
	"time"
)

func TestSpin(t *testing.T) {
	tests := []struct {
		seed int64
		exp  int
	}{
		{1742259393976371944, 6},
		{3702, 20},   // ===
		{3554, 1227}, // ◆◆◆
		{1382, 105},  // 777
		{14073, 100}, // 777
		{11032, 136},
		{562, 122},
		{17395, 35},
		{1742328304162502961, 15},
		{1742330639429664398, 6},
		{2772, 10},
		{9589, 10},
		{1742283085356624423, 20},
		{1742286288251353083, 12},
		{1742286440499332297, 12},
		{1742286608108168119, 19},
		{1742329474192859291, 69},
		{1742329476823589538, 15},
		{9277, 10},
		{15475, 308},
		{27998, 4},
		{29256, 10},
		{5813, 23},
		{26646, 19},
		{498, 6},
		{7035, 0},
		{8472, 0},
		{10472, 69},
		{10222, 10},
		{21296, 0},
		{525, 1313},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Logf("seed: %d", test.seed)
			r := rand.New(rand.NewSource(test.seed))
			res, err := Spin(r, 9)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			t.Logf("\n%s\n", res)
			if res.Payout != test.exp {
				t.Errorf("expected: %d, got: %d", test.exp, res.Payout)
			}
		})
	}
}

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
			t.Logf("pulls: %d/%d (%0.2f%%)", pulls-m[0], pulls, 100.0*float64(pulls-m[0])/float64(pulls))
			keys := slices.Sorted(maps.Keys(m))
			slices.Reverse(keys)
			for _, k := range keys {
				t.Logf("% 5d: %d", k, m[k])
			}
		})
	}
}
