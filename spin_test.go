package tripled

import (
	"bytes"
	"math/rand"
	"strconv"
	"testing"
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
			var buf bytes.Buffer
			_, _ = res.WriteTo(&buf)
			t.Logf("\n%s\n", buf.String())
			if res.Payout != test.exp {
				t.Errorf("expected: %d, got: %d", test.exp, res.Payout)
			}
		})
	}
}
