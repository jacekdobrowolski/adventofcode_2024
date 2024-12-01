package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"slices"
	"strconv"
	"strings"
	"time"
)

const linesAtOnce = 1000
const lineLength = 14

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

	start := time.Now()

	result := Task1("input")

	fmt.Printf("1 took: %v\n", time.Since(start))
	fmt.Printf("1 ans: %d\n", result)

	start = time.Now()

	result = Task2("input")

	fmt.Printf("2 took: %v\n", time.Since(start))
	fmt.Printf("2 ans: %d\n", result)
}

func Task1(filename string) uint32 {
	input, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer input.Close()

	maxCapacity := lineLength * linesAtOnce
	buf := make([]byte, maxCapacity)
	listA := make([]uint32, 0, 1000)
	listB := make([]uint32, 0, 1000)

	for {
		length, err := input.Read(buf)

		linesAtOnce := length / lineLength

		for i := range linesAtOnce {
			i *= lineLength
			listA = append(listA, Parse(([5]byte)(buf[i:i+5])))
			listB = append(listB, Parse(([5]byte)(buf[i+8:i+13])))
		}

		if errors.Is(err, io.EOF) {
			break
		}
	}

	slices.Sort(listA)
	slices.Sort(listB)

	var result uint32

	for i := range listA {
		var distance uint32
		if listA[i] > listB[i] {
			distance = listA[i] - listB[i]
		} else {
			distance = listB[i] - listA[i]
		}

		result += distance
	}
	return result
}

func Parse(input [5]byte) uint32 {
	return uint32(input[4]-byte('0'))*10000 +
		uint32(input[3]-byte('0'))*1000 +
		uint32(input[2]-byte('0'))*100 +
		uint32(input[1]-byte('0'))*10 +
		uint32(input[0]-byte('0'))*1
}

func Task1V2(filename string) uint32 {
	input, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer input.Close()

	maxCapacity := lineLength * linesAtOnce
	buf := make([]byte, maxCapacity)
	listA := make([]uint32, 0, 1000)
	listB := make([]uint32, 0, 1000)

	for {
		_, err = input.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}

		for i := range linesAtOnce {
			split := strings.Split(string(buf[i:i+lineLength]), " ")
			a, _ := strconv.ParseUint(split[0], 10, 32)
			b, _ := strconv.ParseUint(split[1], 10, 32)
			listA = append(listA, uint32(a))
			listB = append(listB, uint32(b))
		}
	}

	slices.Sort(listA)
	slices.Sort(listB)

	var result uint32
	var distance uint32

	for i := range listA {
		if listA[i] > listB[i] {
			distance = listA[i] - listB[i]
		} else {
			distance = listB[i] - listA[i]
		}

		result += distance
	}
	return result
}

func Task2(filename string) uint32 {
	input, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer input.Close()

	var result uint32

	maxCapacity := lineLength * linesAtOnce
	buf := make([]byte, maxCapacity)
	listA := make(map[[5]byte]uint16, 1000)
	listB := make(map[[5]byte]uint16, 1000)

	for {
		_, err := input.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}

		for i := range linesAtOnce {
			i *= lineLength
			listA[([5]byte)(buf[i+0:i+5])] += 1
			listB[([5]byte)(buf[i+8:i+13])] += 1
		}
	}

	for num, aCount := range listA {

		bCount, ok := listB[num]
		if !ok {
			bCount = 0
		}
		result += Parse(num) * uint32(aCount) * uint32(bCount)
	}

	return result
}
