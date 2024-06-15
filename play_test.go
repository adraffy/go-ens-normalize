package ensnormalize

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/adraffy/go-ensnormalize/common"
	"github.com/adraffy/go-ensnormalize/ensip15"
	"github.com/adraffy/go-ensnormalize/nf"
)

// func Test1(t *testing.T) {
// 	norm, err := Normalize("abc")
// 	if err != nil {
// 		t.Errorf("wtf");
// 	} else {
// 		t.Logf("Normalized: %s", norm)
// 	}
// }

func readJSONFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	v, err := io.ReadAll(file)
	if err != nil {
		return v
	}
	return v
}

func TestGo(t *testing.T) {

	// var v []rune
	// v = append(v, 1)
	// fmt.Println(v)

	l := ensip15.GetInstance()
	fmt.Println(l.Normalize("Chonk"))
	fmt.Println(l.Emojis())
}

func TestNF(t *testing.T) {
	nf := nf.New()
	var tests map[string]interface{}
	err := json.Unmarshal(readJSONFile("compress/data/nf-tests.json"), &tests)
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
				if common.CompareRunes(nfd, nfd0) != 0 {
					t.Errorf("NFD[%d]: expect %v, got %v", i, nfd0, nfd)
				}
				if common.CompareRunes(nfc, nfc0) != 0 {
					t.Errorf("NFC[%d]: expect %v, got %v", i, nfc0, nfc)
				}
			}
		})
	}
}
