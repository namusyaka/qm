package qm

import (
	"errors"
	"reflect"
	"sort"
)

type Set []int

func (s Set) isSubset(c Set) bool {
	for _, n := range s {
		var exist bool
		for _, m := range c {
			if n == m {
				exist = true
				break
			}
		}
		if !exist {
			return false
		}
	}
	return true
}

func (s Set) Equal(c Set) bool {
	a, b := make([]int, len(s)), make([]int, len(c))
	copy(a, s)
	copy(b, c)
	sort.Ints(a)
	sort.Ints(b)
	return reflect.DeepEqual(a, b)
}

func (s Set) Combine(c Set) (r [2]int, err error) {
	if c[1] != s[1] {
		return r, errors.New("invalid member")
	}
	x := c[0] ^ s[0]
	if !((x & (^x + 1)) == x) {
		return r, errors.New("xor(first values) is zero or a powerr of two")
	}
	return [2]int{s[0] & c[0], s[1] | x}, nil
}

func (s Set) Contains(n int) bool {
	for _, m := range s {
		if n == m {
			return true
		}
	}
	return false
}

func (s *Set) Add(n int) {
	if !s.Contains(n) {
		*s = append(*s, n)
	}
}

type SetList []Set

func (sl SetList) Contains(s Set) bool {
	if len(sl) == 0 {
		return false
	}
	for _, set := range sl {
		if reflect.DeepEqual(set, s) {
			return true
		}
	}
	return false
}

func (sl *SetList) Add(s Set) {
	if !sl.Contains(s) {
		*sl = append(*sl, s)
	}
}

func (sl *SetList) DeleteAt(i int) {
	if i >= len(*sl) {
		return
	}
	nsl := make(SetList, len(*sl))
	copy(nsl, *sl)
	*sl = append(nsl[:i], nsl[i+1:]...)
}

func (sl SetList) Delete(sl2 SetList) SetList {
	var setlist SetList
	for _, n := range sl {
		var exist bool
		for _, m := range sl2 {
			if reflect.DeepEqual(n, m) {
				exist = true
			}
		}
		if !exist {
			setlist = append(setlist, n)
		}
	}
	return setlist
}

func (sl SetList) Diff(sl2 SetList) SetList {
	var r SetList
	for _, n := range sl {
		var cant bool
		for _, j := range sl2 {
			if reflect.DeepEqual(n, j) {
				cant = true
			}
		}
		if !cant {
			r = append(r, n)
		}
	}
	return r
}

func (sl SetList) Union(sl2 SetList) SetList {
	nsl := make(SetList, len(sl))
	copy(nsl, sl)
	for _, n := range sl2 {
		var exist bool
		for _, m := range sl {
			if n.Equal(m) {
				exist = true
			}
		}
		if !exist {
			nsl = append(nsl, n)
		}
	}
	return nsl
}
