package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"slices"
	"sync"
	"time"
)

const ITERATIONS = 1

//go:embed input
var input []byte

type Num = [2]byte
type Manual = []Num

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

	minTook := time.Hour
	var result = 0

	for range ITERATIONS {
		start := time.Now()

		rulebook, manuals := ParseInput(input)
		result = 0
		results := make(chan int)

		go func() {
			for res := range results {
				result += res
			}
		}()

		wg := sync.WaitGroup{}

		for _, manual := range manuals {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if Valid(manual, rulebook) {
					results <- ParseNum(manual[len(manual)/2])
				}
			}()
		}
		wg.Wait()
		close(results)

		took := time.Since(start)
		if took < minTook {
			minTook = took
		}
	}
	fmt.Printf("day 5 part 1 took %v\n", minTook)
	fmt.Printf("day 5 part 1 result %d\n", result)

	for range ITERATIONS {
		start := time.Now()

		rulebook, manuals := ParseInput(input)
		result = 0
		results := make(chan int)

		go func() {
			for res := range results {
				result += res
			}
		}()

		// wg := sync.WaitGroup{}

		for _, manual := range manuals {
			// wg.Add(1)
			// go func() {
			// 	defer wg.Done()

			if !Valid(manual, rulebook) {
				PrintManual(manual)
			}

			if Fixed(manual, rulebook) {
				PrintManual(manual)
				results <- ParseNum(manual[len(manual)/2])
			}
			// }()
		}
		// wg.Wait()
		close(results)

		took := time.Since(start)
		if took < minTook {
			minTook = took
		}
	}
	fmt.Printf("day 5 part 2 took %v\n", minTook)
	fmt.Printf("day 5 part 2 result %d\n", result)
}
func ParseNum(num Num) int {
	return 10*int(num[0]-'0') + int(num[1]-'0')
}
func Valid(m Manual, rulebook map[Num][]Num) bool {
	forbidenPages := make([]Num, 0, 29*24)
	for _, page := range m {
		if slices.Contains(forbidenPages, page) {
			return false
		}

		if pages, ok := rulebook[page]; ok {
			forbidenPages = append(forbidenPages, pages...)
		}
	}

	return true
}

func Fixed(m Manual, rulebook map[Num][]Num) bool {
	forbidenPages := make([]Num, 0, 29*24)
	wrongOrder := false

	for _, page := range m {
		if slices.Contains(forbidenPages, page) {
			wrongOrder = true
		}

		if pages, ok := rulebook[page]; ok {
			forbidenPages = append(forbidenPages, pages...)
		}
	}

	if !wrongOrder {
		return false
	}

	for !Valid(m, rulebook) {
		forbidenPages := make([][]Num, 0, 29*24)
	INNER:
		for i := 1; i < len(m); i += 1 {
			fmt.Println(i)
			page := m[i]
			for _, fp := range forbidenPages {
				if slices.Contains(fp, page) {
					m[i-1], m[i] = m[i], m[i-1]
					i -= 2
					forbidenPages = forbidenPages[:len(forbidenPages)-1]
					PrintManual(m)
					continue INNER
				}
			}

			if pages, ok := rulebook[page]; ok {
				forbidenPages = append(forbidenPages, pages)
			}
		}
		PrintManual(m)
	}

	return true
}

func ParseInput(input []byte) (map[Num][]Num, []Manual) {
	rulebook := make(map[Num][]Num, 81)
	manuals := make([]Manual, 0, 210)

	var firstPage, secondPage Num
	var pipe byte
	i := 0

	for {
		firstPage, pipe, secondPage = (Num)(input[i:2+i]), input[2+i], (Num)(input[3+i:5+i])

		if pipe != byte('|') {
			break
		}

		i += 6
		rulebook[secondPage] = append(rulebook[secondPage], firstPage)

	}

	i += 1 // move over the empy line
	for i < len(input) {
		if input[i] == '\n' {
			break
		}

		manual := make(Manual, 0, 29)
	MANUAL:
		for i+2 < len(input) {

			firstPage = (Num)(input[i : i+2])

			manual = append(manual, firstPage)

			if input[i+2] == byte('\n') {
				i += 3
				break MANUAL
			}

			i += 3
		}

		manuals = append(manuals, manual)
	}
	return rulebook, manuals
}

func PrintRuleBook(r map[Num][]Num) {
	for key, items := range r {
		fmt.Printf("%s: {", key)
		for _, item := range items {
			fmt.Printf("%s, ", item)
		}
		fmt.Printf("}\n")
	}
}

func PrintManuals(m []Manual) {
	for i, manual := range m {
		fmt.Printf("%d - ", i)
		PrintManual(manual)
	}
}

func PrintManual(manual Manual) {
	for _, item := range manual {
		fmt.Printf("%s, ", item)
	}
	fmt.Printf("\n")
}
