package qm

import (
	"reflect"
	"testing"
)

func TestGetBoolFunc(t *testing.T) {
	tests := [...]struct {
		have SetList
		want string
	}{
		{have: SetList([]Set{Set{3, 0}}), want: "(A AND B)"},
		{have: SetList([]Set{Set{0, 0}}), want: "((NOT A) AND (NOT B))"},
		{have: SetList([]Set{Set{1, 2}}), want: "A"},
		{have: SetList([]Set{Set{1, 2}}), want: "A"},
		{have: SetList([]Set{Set{2, 1}}), want: "B"},
		{have: SetList([]Set{Set{0, 2}}), want: "(NOT A)"},
		{have: SetList([]Set{Set{0, 1}}), want: "(NOT B)"},
		{have: SetList([]Set{Set{1, 2}, Set{2, 1}}), want: "(A OR B)"},
		{have: SetList([]Set{Set{0, 1}, Set{0, 2}}), want: "((NOT B) OR (NOT A))"},
	}
	qm := New([]string{"A", "B"})
	for i, tt := range tests {
		if have, want := qm.GetBoolFunc(tt.have), tt.want; want != have {
			t.Errorf("#%d wrong bool function, want = '%s', have = '%s'", i, want, have)
		}
	}
}

func TestSolve(t *testing.T) {
	type (
		ht struct {
			ones []int
			dc   []int
		}
		wt struct {
			qm      SetList
			complex int
		}
	)
	tests := [...]struct {
		have ht
		want wt
	}{
		{
			have: ht{ones: []int{}, dc: []int{}},
			want: wt{qm: SetList([]Set{Set{0}}), complex: 0},
		},
		{
			have: ht{ones: []int{1, 3}, dc: []int{0, 2}},
			want: wt{qm: SetList([]Set{Set{1}}), complex: 0},
		},
		{
			have: ht{ones: []int{0, 1, 2, 3}, dc: []int{}},
			want: wt{qm: SetList([]Set{Set{1}}), complex: 0},
		},
		{
			have: ht{ones: []int{3}, dc: []int{}},
			want: wt{qm: SetList([]Set{Set{3, 0}}), complex: 2},
		},
		{
			have: ht{ones: []int{0}, dc: []int{}},
			want: wt{qm: SetList([]Set{Set{0, 0}}), complex: 4},
		},
		{
			have: ht{ones: []int{1, 3}, dc: []int{}},
			want: wt{qm: SetList([]Set{Set{1, 2}}), complex: 0},
		},
		{
			have: ht{ones: []int{1}, dc: []int{3}},
			want: wt{qm: SetList([]Set{Set{1, 2}}), complex: 0},
		},
		{
			have: ht{ones: []int{2, 3}, dc: []int{}},
			want: wt{qm: SetList([]Set{Set{2, 1}}), complex: 0},
		},
		{
			have: ht{ones: []int{0, 2}, dc: []int{}},
			want: wt{qm: SetList([]Set{Set{0, 2}}), complex: 1},
		},
		{
			have: ht{ones: []int{0, 1}, dc: []int{}},
			want: wt{qm: SetList([]Set{Set{0, 1}}), complex: 1},
		},
		{
			have: ht{ones: []int{1, 2, 3}, dc: []int{}},
			want: wt{qm: SetList([]Set{Set{1, 2}, Set{2, 1}}), complex: 2},
		},
		{
			have: ht{ones: []int{0, 1, 2}, dc: []int{}},
			want: wt{qm: SetList([]Set{Set{0, 1}, Set{0, 2}}), complex: 4},
		},
	}
	qm := New([]string{"A", "B"})
	for i, tt := range tests {
		complex, out := qm.Solve(tt.have.ones, tt.have.dc)
		if have, want := complex, tt.want.complex; want != have {
			t.Errorf("#%d.0 wrong complexity, want = %d, have = %d", i, want, have)
		}
		if have, want := out, tt.want.qm; !reflect.DeepEqual(want, have) {
			t.Errorf("#%d.1 wrong output, want = %v, have = %v", i, want, have)
		}
	}
}

func TestComputePrimes(t *testing.T) {
	tests := [...]struct {
		have []int
		want SetList
	}{
		{have: []int{0, 1, 2}, want: SetList([]Set{Set{0, 1}, Set{0, 2}})},
		{have: []int{3}, want: SetList([]Set{Set{3, 0}})},
		{have: []int{0}, want: SetList([]Set{Set{0, 0}})},
		{have: []int{1, 3}, want: SetList([]Set{Set{1, 2}})},
		{have: []int{2, 3}, want: SetList([]Set{Set{2, 1}})},
		{have: []int{0, 2}, want: SetList([]Set{Set{0, 2}})},
		{have: []int{0, 1}, want: SetList([]Set{Set{0, 1}})},
		{have: []int{1, 2, 3}, want: SetList([]Set{Set{1, 2}, Set{2, 1}})},
		{have: []int{4, 8, 10, 11, 12, 15, 9, 14}, want: SetList([]Set{Set{8, 6}, Set{8, 3}, Set{10, 5}, Set{4, 8}})},
	}
	q := New([]string{"A", "B"})
	for i, tt := range tests {
		if have, want := q.ComputePrimes(tt.have), tt.want; !reflect.DeepEqual(want, have) {
			t.Errorf("#%d wrong primes, want = %v, have = %v", i, want, have)
		}
	}
}

func TestCalculateComplexity(t *testing.T) {
	tests := [...]struct {
		have [][]int
		want int
	}{
		{
			have: [][]int{{1, 6}},
			want: 0,
		},
		{
			have: [][]int{{0, 6}},
			want: 1,
		},
		{
			have: [][]int{{3, 4}},
			want: 2,
		},
		{
			have: [][]int{{7, 0}},
			want: 3,
		},
		{
			have: [][]int{{1, 6}, {2, 5}, {4, 3}},
			want: 3,
		},
		{
			have: [][]int{{0, 6}, {2, 5}, {4, 3}},
			want: 4,
		},
		{
			have: [][]int{{0, 6}, {0, 5}, {4, 3}},
			want: 5,
		},
		{
			have: [][]int{{0, 6}, {0, 5}, {0, 3}},
			want: 6,
		},
		{
			have: [][]int{{3, 4}, {7, 0}, {5, 2}},
			want: 10,
		},
		{
			have: [][]int{{1, 4}, {7, 0}, {5, 2}},
			want: 11,
		},
		{
			have: [][]int{{2, 4}, {7, 0}, {5, 2}},
			want: 11,
		},
		{
			have: [][]int{{0, 4}, {7, 0}, {5, 2}},
			want: 12,
		},
		{
			have: [][]int{{0, 4}, {0, 0}, {5, 2}},
			want: 15,
		},
		{
			have: [][]int{{0, 4}, {0, 0}, {0, 2}},
			want: 17,
		},
	}
	q := New([]string{"A", "B", "C"})
	for i, tt := range tests {
		var sl SetList
		for _, s := range tt.have {
			sl = append(sl, Set(s))
		}
		if have, want := q.CalculateComplexity(sl), tt.want; want != have {
			t.Errorf("#%d wrong complexity caluculation, want = %d, have = %d", i, want, have)
		}
	}
}
