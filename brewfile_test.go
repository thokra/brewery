package brewery

import (
	"encoding/json"
	"os"
	"testing"
)

func readBrewFile() (*Recipe, error) {
	env := map[string]string{
		"BUILD_NUMBER":     "1",
		"REG_D09_PASSWORD": "password",
		"GC_SERVICE_KEY":   "service_key",
	}

	f, err := os.OpenFile("testdata/brew.1.yaml", os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	return Read(f, env)
}

func TestBrewfileRead(t *testing.T) {
	rec, err := readBrewFile()
	if err != nil {
		t.Fatal(err)
	}
	_ = rec

	// spew.Dump(rec.dependencyGraph)

	il, err := rec.Ingredients()
	if err != nil {
		t.Fatal(err)
	}
	json.NewEncoder(os.Stdout).Encode(il)
}
