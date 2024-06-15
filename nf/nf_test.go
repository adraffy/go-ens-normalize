package nf

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
	v, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return v
}

func TestNF(t *testing.T) {
	nf := New()
	var tests map[string]interface{}
	err := json.Unmarshal(readJSONFile("nf-tests.json"), &tests)
	if err != nil {
		panic(err)
	}
	for name, value := range tests {
		list, ok := value.([]interface{})
		if !ok {
			continue
		}
		t.Run(name, func(t *testing.T) {
			for i, x := range list {
				v := x.([]interface{})
				input := []rune(v[0].(string))
				nfd0 := []rune(v[1].(string))
				nfc0 := []rune(v[2].(string))
				nfd := nf.NFD(input)
				nfc := nf.NFC(input)
				if string(nfd) != string(nfd0) {
					t.Errorf("NFD[%d]: expect %v, got %v", i, nfd0, nfd)
				}
				if string(nfc) != string(nfc0) {
					t.Errorf("NFC[%d]: expect %v, got %v", i, nfc0, nfc)
				}
			}
		})
	}
}
