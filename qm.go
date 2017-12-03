package qm

import (
	"fmt"
	"math"
	"strings"
)

type Qm struct {
	vars []string
	size uint
}

func New(vars []string) *Qm {
	return &Qm{vars: vars, size: uint(len(vars))}
}

func (q *Qm) Solve(ones, dc []int) (complexity int, qm SetList) {
	if len(ones) == 0 {
		return 0, SetList([]Set{Set([]int{0})})
	}
	if len(ones)+len(dc) == 1<<q.size {
		return 0, SetList([]Set{Set([]int{1})})
	}
	cubes := make([]int, len(ones))
	copy(cubes, ones)
	primes := q.ComputePrimes(append(cubes, dc...))
	return q.solve(primes, ones)
}

func (q *Qm) solve(primes SetList, ones Set) (int, SetList) {
	var chart SetList
	for _, o := range ones {
		var cols Set
		for i := 0; i < len(primes); i++ {
			if (o & (^primes[i][1])) == primes[i][0] {
				cols.Add(i)
			}
		}
		chart = append(chart, cols)
	}
	var covers SetList
	if len(chart) > 0 {
		for _, i := range chart[0] {
			covers = append(covers, Set([]int{i}))
		}
	}
	for i := 1; i < len(chart); i++ {
		var newcovers SetList
		for _, cover := range covers {
			for _, pi := range chart[i] {
				x := cover[:]
				x.Add(pi)
				added := true
				for j := len(newcovers) - 1; j >= 0; j-- {
					if x.isSubset(newcovers[j]) {
						newcovers.DeleteAt(j)
					} else if newcovers[j].isSubset(x) && !newcovers[j].Equal(x) {
						added = false
					}
				}
				if added {
					newcovers = append(newcovers, x)
				}
			}
		}
		covers = newcovers
	}
	var result SetList
	min := math.MaxInt32
	for _, cover := range covers {
		var pc SetList
		for _, i := range cover {
			pc = append(pc, primes[i])
		}
		complex := q.CalculateComplexity(pc)
		if complex < min {
			min = complex
			result = pc
		}
	}
	return min, result
}

func parenteses(sep string, p []string) string {
	expr := strings.Join(p, sep)
	if len(p) > 1 {
		return "(" + expr + ")"
	}
	return expr
}

func (q *Qm) GetBoolFunc(terms SetList) string {
	var or []string
	for _, term := range terms {
		var and []string
		for i := 0; i < len(q.vars); i++ {
			j := uint(i)
			if term[0]&(1<<j) > 0 {
				and = append(and, q.vars[i])
			} else if (term[1] & (1 << j)) == 0 {
				and = append(and, fmt.Sprintf("(NOT %s)", q.vars[i]))
			}
		}
		or = append(or, parenteses(" AND ", and))
	}
	return parenteses(" OR ", or)
}

func zip(a, b []SetList) [][2]SetList {
	var r [][2]SetList
	for i, n := range a {
		r = append(r, [2]SetList{n, b[i]})
	}
	return r
}

func (q *Qm) ComputePrimes(cubes []int) SetList {
	var sigma []SetList
	for i, max := uint(0), (q.size + uint(1)); i < max; i++ {
		sigma = append(sigma, SetList{})
	}
	for _, cube := range cubes {
		i := hammingWeight(cube)
		if j := len(sigma); j <= i {
			for len(sigma) <= i {
				sigma = append(sigma, SetList{})
			}
		}
		sigma[i] = append(sigma[i], Set{cube, 0})
	}
	var primes SetList
	for len(sigma) > 0 {
		var (
			n         []SetList
			redundant SetList
		)
		for _, p := range zip(sigma[:len(sigma)-1], sigma[1:]) {
			var nc SetList
			k, j := p[0], p[1]
			for _, a := range k {
				for _, b := range j {
					combined, err := a.Combine(b)
					if err == nil {
						nc.Add(combined[:])
						redundant = redundant.Union(SetList([]Set{a, b}))
					}
				}
			}
			n = append(n, nc)
		}
		var p SetList
		for _, cs := range sigma {
			for _, cube := range cs {
				p.Add(cube)
			}
		}
		primes = p.Diff(redundant).Union(primes)
		sigma = n
	}
	return primes
}

// Calculatecomplexity caluculates complexity.
func (q *Qm) CalculateComplexity(terms SetList) int {
	c := len(terms)
	if c == 1 {
		c = 0
	}
	m := (1 << q.size) - 1
	for _, term := range terms {
		mask := ^term[1] & m
		t := hammingWeight(mask)
		if t == 1 {
			t = 0
		}
		c += t
		c += hammingWeight(^term[0] & mask)
	}
	return c
}

// hammingWeight counts set bits of the input.
func hammingWeight(i int) int {
	n := (i & 0x55555555) + ((i >> 1) & 0x55555555)
	n = (n & 0x33333333) + ((n >> 2) & 0x33333333)
	n = (n & 0x0f0f0f0f) + ((n >> 4) & 0x0f0f0f0f)
	n = (n & 0x00ff00ff) + ((n >> 8) & 0x00ff00ff)
	return ((n & 0x0000ffff) + ((n >> 16) & 0x0000ffff))
}
