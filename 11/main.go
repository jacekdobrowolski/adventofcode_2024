package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

func main() {
	input := "92 0 286041 8034 34394 795 8 2051489"
	stonesStr := strings.Split(input, " ")
	stones := make(map[int]int)

	for _, stoneStr := range stonesStr {
		stone, err := strconv.Atoi(stoneStr)
		if err != nil {
			log.Fatal(err)
		}
		stones[stone]++
	}

	for i := range 75 {
		newStones := make(map[int]int)
		for stone, count := range stones {
			result := Blink(stone)
			if result[0] == 0 {
				newStones[result[1]] += count
			} else {
				newStones[result[0]] += count
				newStones[result[1]] += count
			}
		}
		fmt.Println(i)
		stones = newStones
	}
	result := 0

	for _, sum := range stones {
		result += sum
	}

	fmt.Println(result)
}

func Blink(stone int) [2]int {
	if stone == 0 {
		return [2]int{0, 1}
	}
	digitCount := int(math.Floor(math.Log10(float64(stone)))) + 1
	if digitCount%2 == 0 {
		firstStone := stone / int(math.Pow10(digitCount/2))
		secondStone := stone - firstStone*int(math.Pow10(digitCount/2))

		return [2]int{
			firstStone,
			secondStone,
		}
	}

	return [2]int{0, stone * 2024}

}
