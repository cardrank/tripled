package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/cardrank/tripled"
)

func main() {
	seed := flag.Int64("seed", 0, "seed")
	lines := flag.Int("lines", 9, "lines")
	flag.Parse()
	if err := run(*seed, *lines); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(seed int64, lines int) error {
	t := seed
	if seed == 0 {
		t = time.Now().UnixNano()
	}
	if seed == 0 {
		fmt.Fprintln(os.Stdout, "seed:", t)
	}
	r := rand.New(rand.NewSource(t))
	res, err := tripled.Spin(r, lines)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%c\n", res)
	return nil
}
