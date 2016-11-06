package brewery

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestRecipeSteps(t *testing.T) {
	rec, err := readBrewFile()
	if err != nil {
		t.Fatal(err)
	}

	steps, err := rec.Steps()
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(steps)
}
