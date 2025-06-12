package main

import "testing"

func TestFinder(t *testing.T) {
	type testCase struct {
		target int
		slc    []int
	}
	cases := []testCase{
		{target: 18, slc: []int{7, 11, 9}},
		{target: 11, slc: []int{3, 5, 6}},
	}

	for _, v := range cases {
		s := Finder(v.slc, v.target)
		if s != v.target {
			t.Errorf("got %d , want %d", s, v.target)
		}
	}

}
