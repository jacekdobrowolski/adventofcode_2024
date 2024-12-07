package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"time"
)

//go:embed input
var input []byte

const iterations = 100000

func main() {

	if cpuProfile := os.Getenv("CPU_PROFILE"); cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}

		defer pprof.StopCPUProfile()
	}
	bestTime1 := time.Second
	bestTime2 := time.Second

	results := make([]int, iterations)

	for i := range iterations {
		start := time.Now()
		results[i] = Mul(input)
		took := time.Since(start)

		if took < bestTime1 {
			bestTime1 = took
		}
	}

	fmt.Printf("day 3 part 1 took %v\n", bestTime1)
	fmt.Printf("day 3 part 1 result: %d\n", results[0])

	for i := range iterations {
		start := time.Now()
		results[i] = MulConditional(input)
		took := time.Since(start)

		if took < bestTime2 {
			bestTime2 = took
		}
	}
	fmt.Printf("day 3 part 2 took %v\n", bestTime2)
	fmt.Printf("day 3 part 2 result: %d\n", results[0])

}

var mulToken = []byte("mul(")

var doToken = []byte("do()")

var dontToken = []byte("don't()")

func Mul(input []byte) int {
	var result = 0

	for len(input) > 8 {
		idx := bytes.Index(input, mulToken)
		if idx < 0 {
			break
		}

		idx += len(mulToken)
		input = input[idx:]

		commaIdx := bytes.IndexRune(input[:4], ',')
		if commaIdx < 0 {
			continue
		}

		parenthesisIdx := bytes.IndexRune(input[:9], ')')
		if commaIdx < 0 {
			continue
		}

		if parenthesisIdx < commaIdx {
			continue
		}

		a, ok := Parse(input[:commaIdx])
		if !ok {
			continue
		}

		b, ok := Parse(input[commaIdx+1 : parenthesisIdx])
		if !ok {
			continue
		}

		result = result + a*b
	}

	return result
}

func MulConditional(input []byte) int {
	type cacheStruct struct {
		length  int
		doIdx   int
		dontIdx int
	}

	cache := cacheStruct{length: math.MaxInt}

	findEnableCommands := func(input []byte) (int, int) {
		doIdx := cache.doIdx - (cache.length - len(input))
		if doIdx <= 0 {
			doIdx = bytes.Index(input, doToken)
		}

		dontIdx := cache.dontIdx - (cache.length - len(input))
		if dontIdx <= 0 {
			dontIdx = bytes.Index(input, dontToken)
		}

		cache = cacheStruct{
			length:  len(input),
			doIdx:   doIdx,
			dontIdx: dontIdx,
		}

		return doIdx, dontIdx
	}

	var mulConditional func(input []byte, enabled bool) int = func(input []byte, enabled bool) int {
		result := 0
		for len(input) > 8 {
			mulIdx := bytes.Index(input, mulToken)
			if mulIdx < 0 {
				break
			}

			doIdx, dontIdx := findEnableCommands(input)

			if doIdx < mulIdx && (doIdx < dontIdx || dontIdx < 0) && doIdx >= 0 {
				doIdx += len(doToken)
				input = input[doIdx:]
				enabled = true
				continue
			}

			if dontIdx < mulIdx && (dontIdx < doIdx || doIdx < 0) && dontIdx >= 0 {
				dontIdx += len(dontToken)
				input = input[dontIdx:]

				enabled = false
				continue
			}

			mulIdx += len(mulToken)
			input = input[mulIdx:]

			commaIdx := bytes.IndexRune(input[:4], ',')
			if commaIdx < 0 {
				continue
			}

			parenthesisIdx := bytes.IndexRune(input[:9], ')')
			if commaIdx < 0 {
				continue
			}

			if parenthesisIdx < commaIdx {
				continue
			}

			a, ok := Parse(input[:commaIdx])
			if !ok {
				continue
			}

			b, ok := Parse(input[commaIdx+1 : parenthesisIdx])
			if !ok {
				continue
			}

			input = input[parenthesisIdx:]

			if !enabled {
				continue
			}

			result += a * b
		}
		return result
	}
	return mulConditional(input, true)
}

func Parse(x []byte) (int, bool) {
	result := 0
	for _, b := range x {
		if b < '0' || b > '9' {
			return result, false
		}
		result = result*10 + int(b-byte('0'))
	}

	return result, true
}
