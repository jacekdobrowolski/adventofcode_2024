package main

import (
	"errors"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"testing"
)

var _ uint32

func BenchmarkTask1(b *testing.B) {
	filename := "input"

	b.Run("Task1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Task1(filename)
		}
	})

	for i := 10; i <= 50; i += 2 {
		b.Run("Quicksort "+strconv.Itoa(i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Task1V2(filename, i)
			}
		})
	}
}

func Task1_Parse(filename string) uint32 {
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
			listA = append(listA, Parsev2(([5]byte)(buf[i:i+5])))
			listB = append(listB, Parsev2(([5]byte)(buf[i+8:i+13])))
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

func Parsev2(input [5]byte) uint32 {
	return uint32(input[4]-byte('0'))*10000 +
		uint32(input[3]-byte('0'))*1000 +
		uint32(input[2]-byte('0'))*100 +
		uint32(input[1]-byte('0'))*10 +
		uint32(input[0]-byte('0'))*1
}

func FuzzParse(f *testing.F) {
	f.Add(uint32(10000))
	f.Add(uint32(99999))
	f.Add(uint32(93341))
	f.Fuzz(func(t *testing.T, i uint32) {
		if i >= 10000 && i <= 99999 {
			input := ([5]byte)([]byte(strconv.Itoa(int(i))))
			if i != Parse(input) {
				t.Errorf("%d got %d", i, Parse(input))
			}
		}
	})
}
