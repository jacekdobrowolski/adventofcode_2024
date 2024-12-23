package main

import (
	_ "embed"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"
)

// var input = []byte(`kh-tc
// qp-kh
// de-cg
// ka-co
// yn-aq
// qp-ub
// cg-tb
// vc-aq
// tb-ka
// wh-tc
// yn-cg
// kh-ub
// ta-co
// de-co
// tc-td
// tb-wq
// wh-td
// ta-ka
// td-qp
// aq-cg
// wq-ub
// ub-vc
// de-ta
// wq-aq
// wq-vc
// wh-yn
// ka-de
// kh-ta
// co-tc
// wh-qp
// tb-vc
// td-yn`)

//go:embed input
var input []byte

func main() {
	start := time.Now()
	network := make(map[[2]byte][][2]byte)
	hist := make(map[[2]byte]int)

	for i := 0; i < len(input); i += 6 {
		a, b := input[i:i+2], input[i+3:i+5]
		network[([2]byte)(a)] = append(network[([2]byte)(a)], ([2]byte)(b))
		network[([2]byte)(b)] = append(network[([2]byte)(b)], ([2]byte)(a))
		hist[([2]byte)(a)]++
		hist[([2]byte)(b)]++
	}

	result := 0
	cycles := make(map[string]struct{}, 0)
	for computer, connected := range network {
		for _, c := range connected {
			for _, d := range network[c] {
				if d == computer {
					continue
				}
				for _, e := range network[d] {
					if e == d {
						continue
					}
					if e == computer {
						// found loop
						cycle := []string{string(c[:]), string(d[:]), string(e[:])}
						slices.Sort(cycle)
						cycleStr := strings.Join(cycle, "")
						cycles[cycleStr] = struct{}{}
						result++
					}
				}
			}
		}
	}

	keys := slices.Collect(maps.Keys(network))
	maxCliques := make([][]node, 0)

	maxCliques = BronKerbosch([]node{}, keys, []node{}, network, maxCliques)
	slices.SortFunc(maxCliques, func(a, b []node) int { return len(b) - len(a) })

	fmt.Println(repr(maxCliques[0]))
	fmt.Println(time.Since(start))
}

func BFS(graph map[[2]byte][][2]byte, root [2]byte) map[[2]byte]struct{} {
	reached := make(map[[2]byte]struct{}, 0)
	q := make([][2]byte, 0)
	q = append(q, root)

	for len(q) > 0 {
		v := q[len(q)-1]
		q = q[:len(q)-1]

		for _, w := range graph[v] {
			if _, ok := reached[w]; ok {
				continue
			}

			reached[w] = struct{}{}
			q = append(q, w)
		}
	}

	return reached
}

type node = [2]byte

func BronKerbosch(R, P, X []node, network map[node][]node, maxCliques [][]node) [][]node {
	if len(P) <= 0 && len(X) <= 0 {
		maxCliques = append(maxCliques, R)
	}

	for _, v := range P {
		maxCliques = BronKerbosch(Union(R, []node{v}), Intersection(P, network[v]), Intersection(X, network[v]), network, maxCliques)
		P = Diff(P, []node{v})
		X = Union(X, []node{v})
	}

	return maxCliques
}

func Intersection(a, b []node) []node {
	result := make([]node, 0, max(len(a), len(b)))

	for _, elemA := range a {
		if slices.Contains(b, elemA) {
			result = append(result, elemA)
		}
	}

	return result
}

func Union(a, b []node) []node {
	result := make([]node, 0, len(a)+len(b))

	for _, elemA := range a {
		if !slices.Contains(result, elemA) {
			result = append(result, elemA)
		}
	}

	for _, elemB := range b {
		if !slices.Contains(result, elemB) {
			result = append(result, elemB)
		}
	}

	return result
}

func Diff(a, b []node) []node {
	result := make([]node, 0, len(a))

	for _, elemA := range a {
		if !slices.Contains(b, elemA) {
			result = append(result, elemA)
		}
	}

	return result
}

func repr(s []node) string {
	result := make([]string, 0, len(s))
	for _, elem := range s {
		result = append(result, string(elem[:]))
	}

	slices.Sort(result)

	return strings.Join(result, ",")
}
