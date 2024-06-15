package ensip15

import (
	"sync"
)

var shared *ENSIP15
var once sync.Once

func Shared() *ENSIP15 {
	once.Do(func() {
		shared = New()
	})
	return shared
}

func Normalize(name string) string {
	s, err := Shared().Normalize(name)
	if err != nil {
		panic(err)
	}
	return s
}

func Beautify(name string) string {
	s, err := Shared().Beautify(name)
	if err != nil {
		panic(err)
	}
	return s
}
