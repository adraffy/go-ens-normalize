package nf

import (
	_ "embed"

	"github.com/adraffy/ENSNormalize.go/util"
)

//go:embed nf.bin
var compressed []byte

const (
	SHIFT rune = 24
	MASK  rune = (1 << SHIFT) - 1
	NONE  rune = -1
)

const (
	S0      = 0xAC00
	L0      = 0x1100
	V0      = 0x1161
	T0      = 0x11A7
	L_COUNT = 19
	V_COUNT = 21
	T_COUNT = 28
	N_COUNT = V_COUNT * T_COUNT
	S_COUNT = L_COUNT * N_COUNT
	S1      = S0 + S_COUNT
	L1      = L0 + L_COUNT
	V1      = V0 + V_COUNT
	T1      = T0 + T_COUNT
)

func isHangul(cp rune) bool {
	return cp >= S0 && cp < S1
}
func unpackCC(packed rune) byte {
	return byte(packed >> SHIFT)
}
func unpackCP(packed rune) rune {
	return rune(packed & MASK)
}

type NF struct {
	unicodeVersion string
	exclusions     util.RuneSet
	quickCheck     util.RuneSet
	decomps        map[rune][]rune
	recomps        map[rune]map[rune]rune
	ranks          map[rune]byte
}

func New() *NF {
	d := util.NewDecoder(compressed)
	self := NF{}
	self.unicodeVersion = d.ReadString()
	self.exclusions = util.NewRuneSetFromInts(d.ReadUnique())
	self.quickCheck = util.NewRuneSetFromInts(d.ReadUnique())
	self.decomps = make(map[rune][]rune)
	self.recomps = make(map[rune]map[rune]rune)
	self.ranks = make(map[rune]byte)

	decomp1 := d.ReadSortedUnique()
	decomp1A := d.ReadUnsortedDeltas(len(decomp1))
	for i, cp := range decomp1 {
		self.decomps[rune(cp)] = []rune{rune(decomp1A[i])}
	}
	decomp2 := d.ReadSortedUnique()
	decomp2A := d.ReadUnsortedDeltas(len(decomp2))
	decomp2B := d.ReadUnsortedDeltas(len(decomp2))
	for i, cp := range decomp2 {
		cp := rune(cp)
		cpA := rune(decomp2A[i])
		cpB := rune(decomp2B[i])
		self.decomps[cp] = []rune{cpB, cpA}
		if !self.exclusions.Contains((cp)) {
			recomp := self.recomps[cpA]
			if recomp == nil {
				recomp = make(map[rune]rune)
				self.recomps[cpA] = recomp
			}
			recomp[cpB] = cp
		}
	}
	for i := 1; ; i++ {
		v := d.ReadUnique()
		if len(v) == 0 {
			break
		}
		for _, cp := range v {
			self.ranks[rune(cp)] = byte(i)
		}
	}
	d.AssertEOF()
	return &self
}

func (nf *NF) composePair(a, b rune) rune {
	if a >= L0 && a < L1 && b >= V0 && b < V1 {
		return S0 + (a-L0)*N_COUNT + (b-V0)*T_COUNT
	} else if isHangul(a) && b > T0 && b < T1 && (a-S0)%T_COUNT == 0 {
		return a + (b - T0)
	} else {
		if recomp, ok := nf.recomps[a]; ok {
			if cp, ok := recomp[b]; ok {
				return cp
			}
		}
		return NONE
	}
}

type Packer struct {
	nf    *NF
	buf   []rune
	check bool
}

func (p *Packer) add(cp rune) {
	if cc, ok := p.nf.ranks[cp]; ok {
		p.check = true
		cp |= rune(cc) << SHIFT
	}
	p.buf = append(p.buf, cp)
}

func (p *Packer) fixOrder() {
	if !p.check {
		return
	}
	v := p.buf
	prev := unpackCC(v[0])
	for i := 1; i < len(v); i++ {
		cc := unpackCC(v[i])
		if cc == 0 || prev <= cc {
			prev = cc
			continue
		}
		j := i - 1
		for {
			temp := v[j]
			v[j] = v[j+1]
			v[j+1] = temp
			if j == 0 {
				break
			}
			j = j - 1
			prev = unpackCC(v[j])
			if prev <= cc {
				break
			}
		}
		prev = unpackCC(v[i])
	}
}

func (nf *NF) decomposed(cps []rune) []rune {
	p := Packer{nf: nf}
	var buf []rune
	for _, cp0 := range cps {
		cp := cp0
		for {
			if cp < 0x80 {
				p.buf = append(p.buf, cp)
			} else if isHangul(cp) {
				sIndex := cp - S0
				lIndex := sIndex / N_COUNT
				vIndex := (sIndex % N_COUNT) / T_COUNT
				tIndex := sIndex % T_COUNT
				p.add(L0 + lIndex)
				p.add(V0 + vIndex)
				if tIndex > 0 {
					p.add(T0 + tIndex)
				}
			} else {
				if decomp, ok := nf.decomps[cp]; ok {
					buf = append(buf, decomp...)
				} else {
					p.add(cp)
				}
			}
			if len(buf) == 0 {
				break
			}
			last := len(buf) - 1
			cp = buf[last]
			buf = buf[:last]
		}
	}

	p.fixOrder()
	return p.buf
}

func (nf *NF) composedFromPacked(packed []rune) []rune {
	cps := make([]rune, 0, len(packed))
	var stack []rune
	prevCp := NONE
	var prevCc byte
	for _, p := range packed {
		cc := unpackCC(p)
		cp := unpackCP(p)
		if prevCp == NONE {
			if cc == 0 {
				prevCp = cp
			} else {
				cps = append(cps, cp)
			}
		} else if prevCc > 0 && prevCc >= cc {
			if cc == 0 {
				cps = append(cps, prevCp)
				cps = append(cps, stack...)
				stack = nil
				prevCp = cp
			} else {
				stack = append(stack, cp)
			}
			prevCc = cc
		} else {
			composed := nf.composePair(prevCp, cp)
			if composed != NONE {
				prevCp = composed
			} else if prevCc == 0 && cc == 0 {
				cps = append(cps, prevCp)
				prevCp = cp
			} else {
				stack = append(stack, cp)
				prevCc = cc
			}
		}
	}
	if prevCp != NONE {
		cps = append(cps, prevCp)
		cps = append(cps, stack...)
	}
	return cps
}

func (nf *NF) NFD(cps []rune) []rune {
	v := nf.decomposed(cps)
	for i, x := range v {
		v[i] = unpackCP(x)
	}
	return v
}
func (nf *NF) NFC(cps []rune) []rune {
	return nf.composedFromPacked(nf.decomposed(cps))
}

func (nf *NF) UnicodeVersion() string {
	return nf.unicodeVersion
}
