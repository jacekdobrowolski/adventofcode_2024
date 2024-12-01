package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"slices"
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
	return uint32(input[0]-byte('0'))*10000 +
		uint32(input[1]-byte('0'))*1000 +
		uint32(input[2]-byte('0'))*100 +
		uint32(input[3]-byte('0'))*10 +
		uint32(input[4]-byte('0'))*1
}

// InsertionSort function to sort small subarrays
func InsertionSort(arr []uint32) {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1
		// Shift elements of arr[0..i-1] that are greater than key
		// to one position ahead of their current position
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

// QuickSort function to sort the array
func QuickSort(arr []uint32) {
	if len(arr) <= 1 {
		return
	}

	// If array length is 10 or less, use InsertionSort
	if len(arr) <= 26 {
		InsertionSort(arr)
		return
	}

	// Choose pivot (here we use the last element as pivot)
	pivot := arr[len(arr)-1]

	// Partitioning step
	i := -1
	for j := 0; j < len(arr)-1; j++ {
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i] // Swap values
		}
	}

	// Place pivot in the correct position
	arr[i+1], arr[len(arr)-1] = arr[len(arr)-1], arr[i+1]

	// Recursively apply QuickSort on left and right partitions
	QuickSort(arr[:i+1])
	QuickSort(arr[i+2:])
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

	QuickSort(listA)
	QuickSort(listB)

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
		length, err := input.Read(buf)
		lines := length / lineLength

		for i := range lines {
			i *= lineLength
			listA[([5]byte)(buf[i+0:i+5])] += 1
			listB[([5]byte)(buf[i+8:i+13])] += 1
		}

		if errors.Is(err, io.EOF) {
			break
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
