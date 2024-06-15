package ensnormalize

import (
	"fmt"

	"github.com/adraffy/go-ensnormalize/nf"
)

func main() {
	nf := nf.New()
	fmt.Printf("Version: %s\n", nf.UnicodeVersion())
	//nf.Dump()
}
