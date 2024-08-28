package ensip15

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func readJSONFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	v, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return v
}

func TestNormalize(t *testing.T) {
	l := New()
	type Test struct {
		Name    string `json:"name"`
		Norm    string `json:"norm"`
		Error   bool   `json:"error"`
		Comment string `json:"comment"`
	}
	var tests []Test
	err := json.Unmarshal(readJSONFile("tests.json"), &tests)
	if err != nil {
		panic(err)
	}
	for _, test := range tests {
		if len(test.Norm) == 0 {
			test.Norm = test.Name
		}
		t.Run(ToHexSequence([]rune(test.Name)), func(t *testing.T) {
			norm, err := l.Normalize(test.Name)
			if test.Error {
				if err == nil {
					t.Errorf("expected error: %s", ToHexSequence([]rune(norm)))
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if norm != test.Norm {
				t.Errorf("wrong norm: %s vs %s", ToHexSequence([]rune(test.Norm)), ToHexSequence([]rune(norm)))
			}
		})
	}
}
