package ensip15

import (
	"slices"

	"github.com/adraffy/ENSNormalize.go/nf"
	"github.com/adraffy/ENSNormalize.go/util"
)

func (l *ENSIP15) ShouldEscape() util.RuneSet {
	return l.shouldEscape
}
func (l *ENSIP15) Ignored() util.RuneSet {
	return l.ignored
}

func (l *ENSIP15) Groups() (v []*Group) {
	return slices.Clone(l.groups)
}
func (l *ENSIP15) Emojis() (v []EmojiSequence) {
	return slices.Clone(l.emojis)
}

func (l *ENSIP15) ASCIIGroup() *Group {
	return l._ASCII
}
func (l *ENSIP15) EmojiGroup() *Group {
	return l._EMOJI
}
func (l *ENSIP15) NF() *nf.NF {
	return l.nf
}
