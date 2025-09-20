# go-ens-normalize
0-dependancy [ENSIP-15](https://docs.ens.domains/ensip/15) in Go

* 🏛️ Reference Implementation: [adraffy/ens-normalize.js](https://github.com/adraffy/ens-normalize.js)
	* Unicode: `17.0.0`
	* Spec Hash: [`4febc8f5d285cbf80d2320fb0c1777ac25e378eb72910c34ec963d0a4e319c84`](https://github.com/adraffy/ens-normalize.js/blob/main/derive/output/spec.json)
* ✅️ Passes **100%** [ENSIP-15 Validation Tests](./ensip15/ensip15_test.go)
* ✅️ Passes **100%** [Unicode Normalization Tests](./nf/nf_test.go)

`go get github.com/adraffy/go-ens-normalize@v0.1.1`

### Primary API

```go
// global functions will panic on invalid names
ensip15.Normalize("RaFFY🚴‍♂️.eTh") // "raffy🚴‍♂.eth"

// works like Normalize()
ensip15.Beautify("1⃣2⃣.eth"); // "1️⃣2️⃣.eth"

// member functions return ("", err) on invalid names
ens := ensip15.Shared()
norm, err := ens.Normalize("a_")
```

#### Singleton
```go
ens := ensip15.Shared() // singleton
ens := ensip15.New() // new instance

nf := ensip15.Shared().NF() // singleton
nf := nf.New() // new instance
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
ensip15.Shared().NF().NFC([]rune{0x65, 0x300}) // [0xE8]
ensip15.Shared().NF().NFD([]rune{0xE8})        // [0x65, 0x300]
```

## Build

1. [Sync and Compress](./compress/)
1. `go test ./...`
