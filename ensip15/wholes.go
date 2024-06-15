package ensip15

import (
	"fmt"
	"slices"
	"sort"

	"github.com/adraffy/ENSNormalize.go/util"
)

type Whole struct {
	valid       util.RuneSet
	confused    util.RuneSet
	complements map[rune][]int
}

func decodeWholes(d *util.Decoder, groups []*Group) (wholes []Whole, confusables map[rune]Whole) {
	type Extent struct {
		gs  map[*Group]bool
		cps map[rune]bool
	}
	confusables = make(map[rune]Whole)
	for {
		confused := util.NewRuneSetFromInts(d.ReadUnique())
		if confused.Size() == 0 {
			break
		}
		valid := util.NewRuneSetFromInts(d.ReadUnique())
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

func (l *ENSIP15) checkWhole(group *Group, cps []rune) error {
	var shared []rune
	var universe []int
	prev := 0
	for _, cp := range cps {
		w, ok := l.confusables[cp]
		if ok {
			comp := w.complements[cp]
			if prev == 0 {
				prev = len(comp)
				universe = make([]int, prev)
				copy(universe, comp)
			} else {
				next := 0
				for i := 0; i < prev; i++ {
					if _, ok := slices.BinarySearch(comp, universe[i]); ok {
						universe[next] = universe[i]
						next++
					}
				}
				prev = next
			}
			if prev == 0 {
				return nil
			}
		} else if l.uniqueNonConfusables.Contains(cp) {
			return nil
		} else {
			shared = append(shared, cp)
		}
	}
	if prev > 0 {
	next:
		for i := 0; i < prev; i++ {
			other := l.groups[universe[i]]
			for _, cp := range shared {
				if !other.Contains(cp) {
					continue next
				}
			}
			return fmt.Errorf("%v: %s/%s", ErrWholeConfusable, group, other)
		}
	}
	return nil
}
