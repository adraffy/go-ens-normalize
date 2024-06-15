package ensip15

import (
	"github.com/adraffy/go-ensnormalize/common"
	"github.com/adraffy/go-ensnormalize/decoder"
)

// const (
// 	ASCII      = 0
// 	SCRIPT     = 1
// 	RESTRICTED = 2
// 	EMOJI      = 4
// )

type Group struct {
	index         int
	name          string
	restricted    bool
	cmWhitelisted bool
	primary       common.RuneSet
	secondary     common.RuneSet
}

func (g *Group) Name() string {
	return g.name
}
func (g *Group) IsRestricted() bool {
	return g.restricted
}
func (g *Group) Contains(cp rune) bool {
	return g.primary.Contains(cp) || g.secondary.Contains(cp)
}

func decodeGroups(d *decoder.Decoder) (ret []*Group) {
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
			primary:       common.RuneSetFromInts(d.ReadUnique()),
			secondary:     common.RuneSetFromInts(d.ReadUnique()),
		})
	}
	return ret
}
