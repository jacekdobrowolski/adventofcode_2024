package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"iter"
	"time"
)

//go:embed input
var input []byte
var lineLength int

func init() {
	lineLength = bytes.IndexByte(input, '\n')
}

func main() {
	start := time.Now()

	antenas := make(map[byte][]int)
	for i, char := range input {
		if char == '\n' || char == '.' {
			continue
		}
		antenas[char] = append(antenas[char], i)
	}
	fmt.Println(antenas)

	antinodes := make(map[int]struct{}, 0)

	for _, locations := range antenas {
		for a := range Permutations(locations) {
			for _, antinode := range Harmonics(a[0], a[1]) {
				antinodes[antinode] = struct{}{}
			}
		}
	}

	fmt.Println(len(antinodes))
	fmt.Println(time.Since(start))
}

func Permutations[E any, S [2]E](slice []E) iter.Seq[S] {
	return func(yield func(S) bool) {
		for i, a := range slice[:len(slice)-1] {
			for _, b := range slice[i+1:] {
				yield([...]E{a, b})
			}
		}
	}
}

func Antinodes(antenaA, antenaB int) []int {
	antinodes := make([]int, 0, 2)

	ax, ay := ToXY(antenaA)
	bx, by := ToXY(antenaB)

	diffX := ax - bx
	diffY := ay - by

	antinodeCx := ax + diffX
	antinodeCy := ay + diffY
	if CheckBounds(antinodeCx, antinodeCy) {
		antinodes = append(antinodes, ToOffset(antinodeCx, antinodeCy))
	}

	antinodeDx := bx - diffX
	antinodeDy := by - diffY
	if CheckBounds(antinodeDx, antinodeDy) {
		antinodes = append(antinodes, ToOffset(antinodeDx, antinodeDy))
	}

	return antinodes
}

func Harmonics(antenaA, antenaB int) []int {
	antinodes := make([]int, 0, 2)

	ax, ay := ToXY(antenaA)
	bx, by := ToXY(antenaB)

	diffX := ax - bx
	diffY := ay - by

	antinodeCx := ax
	antinodeCy := ay
	for CheckBounds(antinodeCx, antinodeCy) {
		antinodes = append(antinodes, ToOffset(antinodeCx, antinodeCy))

		antinodeCx += diffX
		antinodeCy += diffY
	}

	antinodeDx := bx
	antinodeDy := by

	for CheckBounds(antinodeDx, antinodeDy) {
		antinodes = append(antinodes, ToOffset(antinodeDx, antinodeDy))

		antinodeDx -= diffX
		antinodeDy -= diffY
	}

	return antinodes
}

func ToXY(offset int) (X, Y int) {
	Y = offset / (lineLength + 1)
	X = offset % (lineLength + 1)

	return X, Y
}

func ToOffset(X, Y int) int {
	return Y*(lineLength+1) + X
}

func Abs(x int) int {
	if x < 0 {
		return x * -1
	}

	return x
}

func CheckBounds(x, y int) bool {
	maxY := len(input)/lineLength - 1
	maxX := lineLength

	return 0 <= x && x < maxX && 0 <= y && y < maxY
}
