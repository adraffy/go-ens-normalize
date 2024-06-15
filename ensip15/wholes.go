package ensip15

import (
	"sort"

	"github.com/adraffy/go-ensnormalize/common"
	"github.com/adraffy/go-ensnormalize/decoder"
)

type Whole struct {
	valid       common.RuneSet
	confused    common.RuneSet
	complements map[rune][]int
}

var UNIQUE_PH = Whole{}

func decodeWholes(d *decoder.Decoder, groups []*Group) (wholes []Whole, confusables map[rune]Whole) {
	type Extent struct {
		gs  map[*Group]bool
		cps map[rune]bool
	}
	confusables = make(map[rune]Whole)
	for {
		confused := d.ReadUniqueRuneSet()
		if confused.Size() == 0 {
			break
		}
		valid := d.ReadUniqueRuneSet()
		whole := Whole{
			valid:       valid,
			confused:    confused,
			complements: make(map[rune][]int),
		}
		wholes = append(wholes, whole)
		for _, cp := range confused.ToArray() {
			confusables[rune(cp)] = whole
		}
		cover := make(map[*Group]bool)
		var extents []*Extent
		for _, cp := range append(valid.ToArray(), confused.ToArray()...) {
			gs := make(map[*Group]bool)
			for _, g := range groups {
				if g.Contains(cp) {
					gs[g] = true
				}
			}
			var ext *Extent
		outer:
			for _, x := range extents {
				for g := range gs {
					if _, ok := x.gs[g]; ok {
						ext = x
						break outer
					}
				}
			}
			if ext == nil {
				ext = &Extent{
					gs:  make(map[*Group]bool),
					cps: make(map[rune]bool),
				}
				extents = append(extents, ext)
			}
			for g := range gs {
				ext.gs[g] = true
				cover[g] = true
			}
			ext.cps[cp] = true
		}
		for _, x := range extents {
			var comps []int
			for g := range cover {
				if _, ok := x.gs[g]; !ok {
					comps = append(comps, g.index)
				}
			}
			sort.Ints(comps)
			for cp := range x.cps {
				whole.complements[cp] = comps
			}
		}
	}
	return wholes, confusables
}
