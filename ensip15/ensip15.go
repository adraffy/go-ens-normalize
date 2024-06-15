package ensip15

import (
	_ "embed"
	"slices"
	"sort"
	"sync"

	"github.com/adraffy/go-ensnormalize/common"
	"github.com/adraffy/go-ensnormalize/decoder"
	"github.com/adraffy/go-ensnormalize/nf"
)

//go:embed spec.bin
var compressed []byte

type ENSIP15 struct {
	nf                   *nf.NF
	shouldEscape         common.RuneSet
	ignored              common.RuneSet
	combiningMarks       common.RuneSet
	nonSpacingMarks      common.RuneSet
	maxNonSpacingMarks   int
	nfcCheck             common.RuneSet
	fenced               map[rune]string
	mapped               map[rune][]rune
	groups               []*Group
	emojis               []EmojiSequence
	emojiRoot            EmojiNode
	possiblyValid        common.RuneSet
	wholes               []Whole
	confusables          map[rune]Whole
	uniqueNonConfusables common.RuneSet
	_LATIN               *Group
	_GREEK               *Group
	_ASCII               *Group
	_EMOJI               *Group
}

var instance *ENSIP15
var once sync.Once

func GetInstance() *ENSIP15 {
	once.Do(func() {
		instance = New()
	})
	return instance
}

func decodeNamedCodepoints(d *decoder.Decoder) map[rune]string {
	ret := make(map[rune]string)
	for _, cp := range d.ReadSortedAscending(d.ReadUnsigned()) {
		ret[rune(cp)] = d.ReadString()
	}
	return ret
}

func decodeMapped(d *decoder.Decoder) map[rune][]rune {
	ret := make(map[rune][]rune)
	for {
		w := d.ReadUnsigned()
		if w == 0 {
			break
		}
		keys := d.ReadSortedUnique()
		n := len(keys)
		m := make([][]rune, n)
		for i := 0; i < n; i++ {
			m[i] = make([]rune, w)
		}
		for j := 0; j < w; j++ {
			v := d.ReadUnsortedDeltas(n)
			for i := 0; i < n; i++ {
				m[i][j] = rune(v[i])
			}
		}
		for i := 0; i < n; i++ {
			ret[rune(keys[i])] = m[i]
		}
	}
	return ret
}

func New() *ENSIP15 {
	d := decoder.New(compressed)
	l := ENSIP15{}
	l.nf = nf.New()
	l.shouldEscape = d.ReadUniqueRuneSet()
	l.ignored = d.ReadUniqueRuneSet()
	l.combiningMarks = d.ReadUniqueRuneSet()
	l.maxNonSpacingMarks = d.ReadUnsigned()
	l.nonSpacingMarks = d.ReadUniqueRuneSet()
	l.nfcCheck = d.ReadUniqueRuneSet()
	l.fenced = decodeNamedCodepoints(d)
	l.mapped = decodeMapped(d)
	l.groups = decodeGroups(d)
	l.emojis = decodeEmojis(d, nil)
	l.wholes, l.confusables = decodeWholes(d, l.groups)

	sort.Slice(l.emojis, func(i, j int) bool {
		return common.CompareRunes(l.emojis[i].normalized, l.emojis[j].normalized) < 0
	})

	l.emojiRoot = makeEmojiTree(l.emojis)

	union := make(map[rune]bool)
	multi := make(map[rune]bool)
	for _, g := range l.groups {
		for _, cp := range append(g.primary.ToArray(), g.secondary.ToArray()...) {
			if _, ok := union[cp]; ok {
				multi[cp] = true
			} else {
				union[cp] = true
			}
		}
	}

	possiblyValid := make(map[rune]bool)
	for cp := range union {
		possiblyValid[cp] = true
		for _, cp := range l.nf.NFD([]rune{cp}) {
			possiblyValid[cp] = true
		}
	}
	l.possiblyValid = common.RuneSetFromKeys(possiblyValid)

	for cp := range multi {
		delete(union, cp)
	}
	for cp := range l.confusables {
		delete(union, cp)
	}
	l.uniqueNonConfusables = common.RuneSetFromKeys(union)

	// direct group references
	l._LATIN = l.FindGroup("Latin")
	l._GREEK = l.FindGroup("Greek")
	l._ASCII = &Group{
		index:         -1,
		restricted:    false,
		name:          "ASCII",
		cmWhitelisted: false,
		primary:       l.possiblyValid.Filter(func(cp rune) bool { return cp < 0x80 }),
		secondary:     common.RuneSet{},
	}
	l._EMOJI = &Group{
		index:         -1,
		restricted:    false,
		cmWhitelisted: false,
		primary:       common.RuneSet{},
		secondary:     common.RuneSet{},
	}
	return &l
}

func (l *ENSIP15) FindGroup(name string) *Group {
	i := slices.IndexFunc(l.groups, func(g *Group) bool {
		return g.name == name
	})
	return l.groups[i]
}

func (l *ENSIP15) Normalize(name string) (string, error) {
	return "yo", nil
}

func (l *ENSIP15) ShouldEscape() common.RuneSet {
	return l.shouldEscape
}

func (l *ENSIP15) Emojis() (v []EmojiSequence) {
	v = make([]EmojiSequence, len(l.emojis))
	copy(v, l.emojis)
	return v
}

func (l *ENSIP15) GroupASCII() *Group {
	return l._ASCII
}
func (l *ENSIP15) GroupEmoji() *Group {
	return l._EMOJI
}
