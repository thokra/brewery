package brewery

import (
	"strings"
	"testing"
)

func TestDependencyGraph(t *testing.T) {
	g := graph{}
	g.add("A", "B")
	g.add("C", "A")
	g.add("D", "B")
	g.add("J", "C", "B")
	g.add("O", "N", "K")
	g.add("N", "B")
	g.add("B")
	g.add("K", "N")

	_, err := g.resolve()
	if err != nil {
		t.Fatal(err)
	}
}

var strSliceTests = []struct {
	in, ex, out string
}{
	{"1,2,3,4,5", "1,3,5", "2,4"},
	{"1,2,3,4,5", "1,3", "2,4,5"},
	{"1,2,3,4,5", "1", "2,3,4,5"},
	{"1,2,3,4,5", "1,2,3,4,5", ""},
	{"1,2,3,4,5", "", "1,2,3,4,5"},
	{"3,2,1", "2,1", "3"},
}

func TestStrSliceExclude(t *testing.T) {
	for ti, tt := range strSliceTests {
		in := strings.Split(tt.in, ",")
		exclude := strings.Split(tt.ex, ",")
		out := strings.Split(tt.out, ",")
		if tt.ex == "" {
			exclude = []string{}
		}
		if tt.out == "" {
			out = []string{}
		}

		ns := strSliceExclude(in, exclude...)
		if len(ns) != len(out) {
			t.Errorf("Testcase %v with invalid length. Expected %v got %v", ti, len(out), len(ns))
			continue
		}

		for i, v := range ns {
			if out[i] != v {
				t.Errorf("Test case %v, index %v had wrong value. Expected %v got %v", ti, i, out[i], v)
			}
		}
	}
}
