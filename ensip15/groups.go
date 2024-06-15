package ensip15

import (
	"fmt"
	"slices"

	"github.com/adraffy/ENSNormalize.go/util"
)

type Group struct {
	index         int
	name          string
	restricted    bool
	cmWhitelisted bool
	primary       util.RuneSet
	secondary     util.RuneSet
}

func (g *Group) Name() string {
	return g.name
}
func (g *Group) String() string {
	if g.restricted {
		return fmt.Sprintf("Restricted[%s]", g.name)
	} else {
		return g.name
	}
}
func (g *Group) IsRestricted() bool {
	return g.restricted
}
func (g *Group) Contains(cp rune) bool {
	return g.primary.Contains(cp) || g.secondary.Contains(cp)
}

func (l *ENSIP15) FindGroup(name string) *Group {
	i := slices.IndexFunc(l.groups, func(g *Group) bool {
		return g.name == name
	})
	return l.groups[i]
}

func decodeGroups(d *util.Decoder) (ret []*Group) {
	for {
		name := d.ReadString()
		if len(name) == 0 {
			break
		}
		bits := d.ReadUnsigned()
		ret = append(ret, &Group{
			index:         len(ret),
			name:          name,
			restricted:    (bits & 1) != 0,
			cmWhitelisted: (bits & 2) != 0,
			primary:       util.NewRuneSetFromInts(d.ReadUnique()),
			secondary:     util.NewRuneSetFromInts(d.ReadUnique()),
		})
	}
	return ret
}

func (l *ENSIP15) determineGroup(unique []rune) (*Group, error) {
	gs := slices.Clone(l.groups)
	prev := len(gs)
	for _, cp := range unique {
		next := 0
		for i := 0; i < prev; i++ {
			if gs[i].Contains(cp) {
				gs[next] = gs[i]
				next++
			}
		}
		if next == 0 {
			for _, g := range gs {
				if g.Contains(cp) {
					return nil, l.createMixtureError(gs[0], cp)
				}
			}
			return nil, fmt.Errorf("%w: %s", ErrDisallowedCharacter, l.SafeCodepoint(cp))
		}
		prev = next
		if prev == 1 {
			break
		}
	}
	return gs[0], nil
}

func (l *ENSIP15) checkGroup(group *Group, cps []rune) error {
	for _, cp := range cps {
		if !group.Contains(cp) {
			return l.createMixtureError(group, cp)
		}
	}
	if !group.cmWhitelisted {
		decomposed := l.nf.NFD(cps)
		e := len(decomposed)
		for i := 1; i < e; i++ {
			if l.nonSpacingMarks.Contains(decomposed[i]) {
				j := i + 1
				for ; j < e; j++ {
					cp := decomposed[j]
					if !l.nonSpacingMarks.Contains(cp) {
						break
					}
					for k := i; k < j; k++ {
						if decomposed[k] == cp {
							return fmt.Errorf("%w: %s", ErrNSMDuplicate, l.SafeCodepoint((cp)))
						}
					}
				}
				n := j - i
				if n > l.maxNonSpacingMarks {
					return fmt.Errorf("%w: %s (%d/%d)", ErrNSMExcessive, l.SafeImplode(decomposed[i-1:j]), n, l.maxNonSpacingMarks)
				}
				i = j
			}
		}
	}
	return nil
}
