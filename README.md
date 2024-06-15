# ENSNormalize.go
0-dependancy [ENSIP-15](https://docs.ens.domains/ensip/15) in Go

* Reference Implementation: [adraffy/ens-normalize.js](https://github.com/adraffy/ens-normalize.js)
	* Unicode: `15.1.0`
	* Spec Hash: [`1f6d3bdb7a724fe3b91f6d73ab14defcb719e0f4ab79022089c940e7e9c56b9c`](https://raw.githubusercontent.com/adraffy/ens-normalize.js/main/derive/output/spec.json)
* Passes **100%** [ENSIP-15 Validation Tests](./ensip15/ensip15_test.go)
* Passes **100%** [Unicode Normalization Tests](./nf/nf_test.go)

### Primary API

```go
import "github.com/adraffy/ENSNormalize.go/ensip15"

// panics on invalid names
ensip15.Normalize("RaFFY🚴‍♂️.eTh") // "raffy🚴‍♂.eth"
// same as: ensip15.Shared().Normalize(...)

// works like Normalize()
ensip15.Beautify("1⃣2⃣.eth"); // "1️⃣2️⃣.eth"
```

### Error Handling

All [errors](./ensip15/errors.go) are safe to print.

### Utilities

Normalize name fragments for substring search:

```go
ensip15.Shared().NormalizeFragment("AB--", false) // "ab--"
ensip15.Shared().NormalizeFragment("..\u0300", false) // "..̀"
ensip15.Shared().NormalizeFragment("\u03BF\u043E", false) // "οо"
// note: Normalize() errors on these inputs
```

Construct safe strings:

```go
ensip15.Shared().SafeCodepoint(0x303) // "◌̃ {303}"
ensip15.Shared().SafeCodepoint(0xFE0F) // "{FE0F}"
ensip15.Shared().SafeImplode([]rune{0x303, 0xFE0F}) // "◌̃{FE0F}"
```

### Unicode Normalization Forms

```go
// use shared instance
nf := ensip15.Shared().NF()
nf.NFC([]rune{0x65, 0x300}) // [0xE8]
nf.NFD([]rune{0xE8})        // [0x65, 0x300]

// create new instance
import "github.com/adraffy/ENSNormalize.go/nf"
nf := nf.New() 
```