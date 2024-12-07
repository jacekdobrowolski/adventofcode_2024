package main

import (
	_ "embed"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed input
var input string

const iterations = 100

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
		input = strings.TrimSuffix(input, "\n")
		start := time.Now()
		results[i] = SumValid(input, Valid)
		took := time.Since(start)

		if took < bestTime1 {
			bestTime1 = took
		}
	}

	fmt.Printf("day 7 part 1 took %v\n", bestTime1)
	fmt.Printf("day 7 part 1 result: %d\n", results[0])

	for i := range iterations {
		start := time.Now()
		results[i] = SumValid(input, Valid2)
		took := time.Since(start)
		if took < bestTime2 {
			bestTime2 = took
		}
	}

	fmt.Printf("day 7 part 2 took %v\n", bestTime2)
	fmt.Printf("day 7 part 2 result: %d\n", results[0])
}

func SumValid(input string, validate func(int, int, []int) bool) int {
	result := 0
	equations := strings.Split(input, "\n")
	results := make(chan int)

	done := make(chan struct{})
	go func() {
		for r := range results {
			result += r
		}
		done <- struct{}{}
	}()

	wg := sync.WaitGroup{}

	for _, equation := range equations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			value, equation := ParseEquation(equation)
			if validate(value, equation[0], equation[1:]) {
				results <- value
			}
		}()
	}

	wg.Wait()
	close(results)
	<-done
	return result
}

func ParseEquation(equation string) (int, []int) {
	splitByColon := strings.SplitN(equation, ":", 2)

	valueStr := splitByColon[0]
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(err)
	}

	if len(splitByColon) < 2 {
		return value, []int{}
	}

	numbersStr := splitByColon[1]
	numbers := make([]int, 0, len(numbersStr))
	for _, numberStr := range strings.Split(strings.TrimLeft(numbersStr, " "), " ") {
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			panic(err)
		}

		numbers = append(numbers, number)
	}

	return value, numbers
}

func Valid(goal, value int, numbers []int) bool {
	if value > goal {
		return false
	}

	for _, num := range numbers {
		a := Valid(goal, value+num, numbers[1:])
		b := Valid(goal, value*num, numbers[1:])

		return a || b
	}

	return value == goal
}

func Valid2(goal, value int, numbers []int) bool {
	if value > goal {
		return false
	}

	for _, num := range numbers {
		a := Valid2(goal, value+num, numbers[1:])
		b := Valid2(goal, value*num, numbers[1:])
		cValue := value*int(math.Pow10(int(math.Floor(math.Log10(float64(num)))))+1) + num
		c := Valid2(goal, cValue, numbers[1:])

		return a || b || c
	}

	return value == goal
}
