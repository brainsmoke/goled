package main

import (
    "post6.net/goled/polyhedron"
    "fmt"
)

const (
    TopLeft     = 0
    BottomLeft  = 1
    BottomRight = 2
    TopRight    = 3
)


var traversal = polyhedron.RemapRoute{

    BottomLeft, TopLeft,
    TopLeft,
    BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
    TopLeft,
    BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
    TopLeft,
    BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
    TopLeft,
    BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
    TopLeft,
    BottomLeft, BottomLeft, TopLeft,

    BottomLeft, TopLeft,
    TopLeft,
    BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
    TopLeft,
    BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
    TopLeft,
    BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
    TopLeft,
    BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
    TopLeft,
    BottomLeft, BottomLeft, TopLeft,
}

var pattern = [][]uint{

	{BottomLeft, BottomLeft, BottomLeft},
	{BottomLeft, BottomLeft, TopLeft},
	{BottomLeft, TopLeft, BottomLeft},
	{BottomLeft, TopLeft, TopLeft},
	{TopLeft, BottomLeft, BottomLeft},
	{TopLeft, BottomLeft, TopLeft},
	{TopLeft, TopLeft, BottomLeft},

	{BottomLeft, BottomLeft},
	{BottomLeft, TopLeft},
	{TopLeft, BottomLeft},
	{TopLeft, TopLeft},

	{BottomLeft},
	{TopLeft},

//	{},
}


const full = (1<<30) - 1
const bad = (1<<30)

func find_half_mappings(start, depth, used int, collision uint64, possible, bitmask []uint64, chains *[][]int, log *[30]int, ch chan<- []int) {

	if depth >= 8 || used+2 < depth*4 {
		return
	}

	for i:=start; i<len(bitmask); i++ {

		if bitmask[i] == bad {
			continue
		}

		if (bitmask[i] & collision) == 0 && (possible[i] | collision) == full {

			c := collision | bitmask[i]
			log[depth] = i

			if c == full {
				ch <- append([]int(nil),log[:depth+1]...)
			} else {
				find_half_mappings( i+1, depth+1, used+len((*chains)[i]), c, possible, bitmask, chains, log, ch)
			}
		}
	}

	if depth == 0 {
		close(ch)
	}
}


func find_mappings(sp []int,ch chan<- [][]int) {

    faces := polyhedron.RemapFaces(polyhedron.DeltoidalHexecontahedronFaces(), 0, traversal)
	n := len(pattern)
	half := 30*n
	var bitmask = make([]uint64, half)
	var chains = make([][]int, 60*n)

	for i:=0; i<30; i++ {
		for j:=0; j<n; j++ {

			ix := i
			bitmask[i*n+j] = 1<<uint(ix)
			chains[i*n+j] = append(chains[i*n+j], ix)
			chains[i*n+j+half] = append(chains[i*n+j+half], ix+30)

			for k:=0; k<len(pattern[j]); k++ {
				ix = faces[ix].Neighbours[pattern[j][k]]
				bitmask[i*n+j] |= 1<<uint(ix)
				chains[i*n+j] = append(chains[i*13+j], ix)
				chains[i*n+j+half] = append(chains[i*13+j+half], ix+30)
			}

			if bitmask[i*n+j] &^ full != 0 {
				bitmask[i*n+j] = bad
			}
		}
	}

	var possible []uint64 = make([]uint64, half)

	for i := range possible {
		for j:=i; j<half; j++ {
			if bitmask[j] != bad {
				possible[i] |= bitmask[j]
			}
		}
	}

	var log [30]int
	ch_in := make(chan []int)
	go find_half_mappings(0, 0, 0, 0, possible, bitmask, &chains, &log, ch_in)

	var a_list, b_list [][][]int

	a_best, b_best := 9001, 9001
	for mapping := range ch_in {
		chain_a := make([][]int, 8)
		chain_b := make([][]int, 8)
		for i:=range mapping {
			chain_a[i] = chains[mapping[i]]
			chain_b[i] = chains[390+mapping[i]]
		}
		a_score := score(chain_a, sp)
		b_score := score(chain_b, sp)
		if a_score < a_best {
			a_best = a_score
			a_list = [][][]int{}
		}
		if b_score < b_best {
			b_best = b_score
			b_list = [][][]int{}
		}
		if a_score == a_best {
			a_list = append(a_list, chain_a)
		}
		if b_score == b_best {
			b_list = append(b_list, chain_b)
		}
	}

	for a := range a_list {
		for b := range b_list {
			full_list := append(append([][]int{}, a_list[a]...), b_list[b]...)
			ch <- full_list
		}
	}
	close(ch)
}

func score(m [][]int, sp []int) int {
	s := 0

	for i := range m {
		s += sp[m[i][0]]
	}
	return s
}


func max_traces(needed [60][3]int) ([3]int, int) {

	var max [3]int
	var sum int = 0

	for j := 0; j<3; j++ {
		for i := range needed {
			if max[j] < needed[i][j] {
				max[j] = needed[i][j]
			}
		}
		sum += max[j]
	}

	return max, sum
}


func calc_traces_x(depth int, mapping [][]int, origin *[60][][][2]int, needed *[60][3]int, log *[][][2]int) {

	if depth == len(mapping) {
		max, smax := max_traces(*needed)
		if smax == 8 {
			fmt.Printf("  %+v %+v\n", max, smax)
			fmt.Printf("  %+v\n", *log)
		}
		return
	}

	choices := (*origin)[mapping[depth][0]]
	for i := range choices {

		for j := 0; j <= len(choices[i]); j++ {
			var face, from, to int

			if j == len(choices[i]) {
				face = mapping[depth][0]
				to = 3
			} else {
				face = choices[i][j][0]
				to =  choices[i][j][1]
			}

			if j == 0 {
				if (to|1) == 1 {
					needed[face][1]+=1
				} else {
					needed[face][2]+=1
				}
			} else {

				from = 3-choices[i][j-1][1]
				if (from|1) != (to|1) {
					needed[face][0]+=1
				}
			}
		}

		(*log)[depth] = choices[i]
		calc_traces_x(depth+1, mapping, origin, needed, log)

		for j := 0; j <= len(choices[i]); j++ {
			var face, from, to int

			if j == len(choices[i]) {
				face = mapping[depth][0]
				to = 3
			} else {
				face = choices[i][j][0]
				to =  choices[i][j][1]
			}

			if j == 0 {
				if (to|1) == 1 {
					needed[face][1]-=1
				} else {
					needed[face][2]-=1
				}
			} else {

				from = 3-choices[i][j-1][1]
				if (from|1) != (to|1) {
					needed[face][0]-=1
				}
			}
		}


	}
}


func calc_traces(ch <-chan [][]int, origin *[60][][][2]int) {

	//var needed [60][4]int

	for mapping := range ch {
		fmt.Printf("%+v\n", mapping)
		var needed [60][3]int
		log := make([][][2]int, 16)
		calc_traces_x(0, mapping, origin, &needed, &log)
	}
}


func main() {

    faces := polyhedron.RemapFaces(polyhedron.DeltoidalHexecontahedronFaces(), 0, traversal)

	var sp [60]int
	var origin [60][][][2]int

	for i := range sp {
		sp[i] = -1
	}
	sp[1] = 0
	sp[2] = 0
	sp[3] = 0
	origin[1] = [][][2]int{[][2]int{}}
	origin[2] = [][][2]int{[][2]int{}}
	origin[3] = [][][2]int{[][2]int{}}

	for left, i := 60-3, 0; left > 0 ; i++ {
		for ix := range faces {
			if sp[ix] == i {
				for move, n := range faces[ix].Neighbours {
					if sp[n] == -1 {
						sp[n] = i+1
						left--
					}
					if sp[n] == i+1 {
						for q := range origin[ix] {
							path := append([][2]int{}, append(origin[ix][q], [2]int{ix, move})...)
							origin[n] = append(origin[n], path)
						}
					}
				}
			}
		}
	}

	ch := make(chan [][]int)
	go find_mappings(sp[:], ch)
	calc_traces(ch, &origin)
}
