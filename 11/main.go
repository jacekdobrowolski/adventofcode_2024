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
	stones := make([]int, 0, len(stonesStr))

	for _, stoneStr := range stonesStr {
		stone, err := strconv.Atoi(stoneStr)
		if err != nil {
			log.Fatal(err)
		}
		stones = append(stones, stone)
	}

	var newStones []int

	for i := range 75 {

		newStones = make([]int, 0, len(stones))
		for _, stone := range stones {
			result := Blink(stone)
			if result[0] == 0 {
				newStones = append(newStones, result[1])
			} else {
				newStones = append(newStones, result[0], result[1])
			}
		}
		fmt.Println(i)
		stones = newStones
	}

	fmt.Println(len(newStones))
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
