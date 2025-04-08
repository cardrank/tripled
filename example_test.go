package tripled_test

import (
	"fmt"
	"math/rand"

	"github.com/cardrank/tripled"
)

func Example() {
	for i, seed := range []int64{
		507,
		7689,
		9085,
		18623,
		2931,
		16614,
	} {
		if i != 0 {
			fmt.Println()
		}
		// note: use a real rng source
		r := rand.New(rand.NewSource(seed))
		res, err := tripled.DefaultDist.Spin(r, 9)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}
	// Output:
	// pos: 1 12 17
	// ..=
	// -=.
	// ..-
	// payout: 0x
	//
	// pos: 7 0 8
	// .◆.
	// ≡.≡
	// ...
	// lines:
	//  2 payouts 2x
	//  6 payouts 120x
	// payout: 122x
	//
	// pos: 16 11 11
	// 777
	// ...
	// -=-
	// lines:
	//  2 payouts 100x
	//  3 payouts 5x
	// payout: 105x
	//
	// pos: 7 9 6
	// .≡.
	// ≡.-
	// .7.
	// lines:
	//  6 payouts 5x
	// payout: 5x
	//
	// pos: 12 19 17
	// ===
	// ...
	// -.-
	// lines:
	//  2 payouts 20x
	// payout: 20x
	//
	// pos: 0 21 1
	// ◆..
	// .◆.
	// -.-
	// lines:
	//  1 payouts 2x
	//  2 payouts 2x
	//  4 payouts 90x
	//  5 payouts 2x
	//  8 payouts 30x
	//  9 payouts 10x
	// payout: 136x
}
