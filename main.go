package main

import (
	"fmt"

	"github.com/adraffy/ENSNormalize.go/ensip15"
)

func main() {
	fmt.Println(ensip15.Normalize("RaFFY🚴‍♂️.eTh"))
	fmt.Println(ensip15.Beautify("1⃣2⃣.eth"))

	fmt.Println(ensip15.Shared().NormalizeFragment("AB--", false))
	fmt.Println(ensip15.Shared().NormalizeFragment("..\u0300", false))
	fmt.Println(ensip15.Shared().NormalizeFragment("\u03BF\u043E", false))

	fmt.Println(ensip15.Shared().SafeCodepoint(0x303))               // "◌̃ {303}"
	fmt.Println(ensip15.Shared().SafeCodepoint(0xFE0F))              // "{FE0F}"
	fmt.Println(ensip15.Shared().SafeImplode([]rune{0x303, 0xFE0F})) // "◌̃{FE0F}"

	nf := ensip15.Shared().NF()
	fmt.Println(nf.NFC([]rune{0x65, 0x300}))
	fmt.Println(nf.NFD([]rune{0xE8}))
}
