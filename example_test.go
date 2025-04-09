package tripled_test

import (
	"fmt"
	"math/rand"

	"github.com/cardrank/tripled"
)

func Example() {
	for i, seed := range []int64{
		507,
		2931,
		4981,
		7689,
		9085,
		12415,
		16614,
		18623,
	} {
		if i != 0 {
			fmt.Println("---------")
		}
		fmt.Println("seed:", seed)
		// note: use a real rng source
		r := rand.New(rand.NewSource(seed))
		res, err := tripled.DefaultDist.Spin(r, 9)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}
	// Output:
	// seed: 507
	// pos: 1 12 17
	// ..=
	// -=.
	// ..-
	// payout: 0x
	// ---------
	// seed: 2931
	// pos: 12 19 17
	// ===
	// ...
	// -.-
	// lines:
	//  2 payouts 20x
	// payout: 20x
	// ---------
	// seed: 4981
	// pos: 21 15 15
	// .=-
	// ◆..
	// .≡=
	// lines:
	//  1 payouts 2x
	//  6 payouts 2x
	//  7 payouts 2x
	// payout: 6x
	// ---------
	// seed: 7689
	// pos: 7 0 8
	// .◆.
	// ≡.≡
	// ...
	// lines:
	//  2 payouts 2x
	//  6 payouts 120x
	// payout: 122x
	// ---------
	// seed: 9085
	// pos: 16 11 11
	// 777
	// ...
	// -=-
	// lines:
	//  2 payouts 100x
	//  3 payouts 5x
	// payout: 105x
	// ---------
	// seed: 12415
	// pos: 0 21 0
	// ◆.◆
	// .◆.
	// -..
	// lines:
	//  1 payouts 2x
	//  2 payouts 10x
	//  4 payouts 10x
	//  5 payouts 90x
	//  8 payouts 2x
	//  9 payouts 1199x
	// payout: 1313x
	// ---------
	// seed: 16614
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
	// ---------
	// seed: 18623
	// pos: 7 9 6
	// .≡.
	// ≡.-
	// .7.
	// lines:
	//  6 payouts 5x
	// payout: 5x
}
