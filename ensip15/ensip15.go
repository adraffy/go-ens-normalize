package ensip15

import (
	_ "embed"
	"fmt"
	"sort"

	"github.com/adraffy/ENSNormalize.go/nf"
	"github.com/adraffy/ENSNormalize.go/util"
)

//go:embed spec.bin
var compressed []byte

type ENSIP15 struct {
	nf                   *nf.NF
	shouldEscape         util.RuneSet
	ignored              util.RuneSet
	combiningMarks       util.RuneSet
	nonSpacingMarks      util.RuneSet
	maxNonSpacingMarks   int
	nfcCheck             util.RuneSet
	fenced               map[rune]string
	mapped               map[rune][]rune
	groups               []*Group
	emojis               []EmojiSequence
	emojiRoot            *EmojiNode
	possiblyValid        util.RuneSet
	wholes               []Whole
	confusables          map[rune]Whole
	uniqueNonConfusables util.RuneSet
	_LATIN               *Group
	_GREEK               *Group
	_ASCII               *Group
	_EMOJI               *Group
}

func decodeNamedCodepoints(d *util.Decoder) map[rune]string {
	ret := make(map[rune]string)
	for _, cp := range d.ReadSortedAscending(d.ReadUnsigned()) {
		ret[rune(cp)] = d.ReadString()
	}
	return ret
}

func decodeMapped(d *util.Decoder) map[rune][]rune {
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
	d := util.NewDecoder(compressed)
	l := ENSIP15{}
	l.nf = nf.New()
	l.shouldEscape = util.NewRuneSetFromInts(d.ReadUnique())
	l.ignored = util.NewRuneSetFromInts(d.ReadUnique())
	l.combiningMarks = util.NewRuneSetFromInts(d.ReadUnique())
	l.maxNonSpacingMarks = d.ReadUnsigned()
	l.nonSpacingMarks = util.NewRuneSetFromInts(d.ReadUnique())
	l.nfcCheck = util.NewRuneSetFromInts(d.ReadUnique())
	l.fenced = decodeNamedCodepoints(d)
	l.mapped = decodeMapped(d)
	l.groups = decodeGroups(d)
	l.emojis = decodeEmojis(d, nil)
	l.wholes, l.confusables = decodeWholes(d, l.groups)

	sort.Slice(l.emojis, func(i, j int) bool {
		return compareRunes(l.emojis[i].normalized, l.emojis[j].normalized) < 0
	})

	l.emojiRoot = makeEmojiTree(l.emojis)

	union := make(map[rune]bool)
	multi := make(map[rune]bool)
	for _, g := range l.groups {
		for _, cp := range append(g.primary.ToArray(), g.secondary.ToArray()...) {
			if union[cp] {
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
	l.possiblyValid = util.NewRuneSetFromKeys(possiblyValid)

	for cp := range multi {
		delete(union, cp)
	}
	for cp := range l.confusables {
		delete(union, cp)
	}
	l.uniqueNonConfusables = util.NewRuneSetFromKeys(union)

	// direct group references
	l._LATIN = l.FindGroup("Latin")
	l._GREEK = l.FindGroup("Greek")
	l._ASCII = &Group{
		index:         -1,
		restricted:    false,
		name:          "ASCII",
		cmWhitelisted: false,
		primary:       l.possiblyValid.Filter(func(cp rune) bool { return cp < 0x80 }),
	}
	l._EMOJI = &Group{
		index:         -1,
		restricted:    false,
		cmWhitelisted: false,
	}
	return &l
}

func (l *ENSIP15) Normalize(name string) (string, error) {
	return l.transform(
		name,
		l.nf.NFC,
		func(e EmojiSequence) []rune { return e.normalized },
		func(tokens []OutputToken) (string, error) {
			cps := FlattenTokens(tokens)
			_, err := l.checkValidLabel(cps, tokens)
			if err != nil {
				return "", err
			}
			return string(cps), nil
		},
	)
}

func (l *ENSIP15) Beautify(name string) (string, error) {
	return l.transform(
		name,
		l.nf.NFC,
		func(e EmojiSequence) []rune { return e.beautified },
		func(tokens []OutputToken) (string, error) {
			cps := FlattenTokens(tokens)
			_, err := l.checkValidLabel(cps, tokens)
			if err != nil {
				return "", nil
			}
			return string(cps), nil
		},
	)
}

func (l *ENSIP15) NormalizeFragment(frag string, decompose bool) (string, error) {
	nf := l.nf.NFC
	if decompose {
		nf = l.nf.NFD
	}
	return l.transform(
		frag,
		nf,
		func(e EmojiSequence) []rune { return e.normalized },
		func(tokens []OutputToken) (string, error) {
			return string(FlattenTokens(tokens)), nil
		},
	)
}

func (l *ENSIP15) transform(
	name string,
	nf func([]rune) []rune,
	ef func(EmojiSequence) []rune,
	normalizer func(tokens []OutputToken) (string, error),
) (string, error) {
	labels := Split(name)
	for i, label := range labels {
		cps := []rune(label)
		tokens, err := l.outputTokenize(cps, nf, ef)
		if err == nil {
			var norm string
			norm, err = normalizer(tokens)
			if err == nil {
				labels[i] = norm
				continue
			}
		}
		if len(labels) > 0 {
			err = fmt.Errorf("invalid label \"%s\": %w", l.SafeImplode(cps), err)
		}
		return "", err
	}
	return Join(labels), nil
}

func checkLeadingUnderscore(cps []rune) error {
	const UNDERSCORE = 0x5F
	allowed := true
	for _, cp := range cps {
		if allowed {
			if cp != UNDERSCORE {
				allowed = false
			}
		} else {
			if cp == UNDERSCORE {
				return ErrLeadingUnderscore
			}
		}
	}
	return nil
}

func checkLabelExtension(cps []rune) error {
	const HYPHEN = 0x2D
	if len(cps) >= 4 && cps[2] == HYPHEN && cps[3] == HYPHEN {
		return fmt.Errorf("%w: %s", ErrInvalidLabelExtension, string(cps[:4]))
	}
	return nil
}

func (l *ENSIP15) checkCombiningMarks(tokens []OutputToken) error {
	for i, x := range tokens {
		if x.Emoji == nil {
			cp := x.Codepoints[0]
			if l.combiningMarks.Contains(cp) {
				if i == 0 {
					return fmt.Errorf("%v: %s", ErrCMLeading, l.SafeCodepoint(cp))
				} else {
					return fmt.Errorf("%v: %s + %s", ErrCMAfterEmoji, tokens[i-1].Emoji.Beautified(), l.SafeCodepoint(cp))
				}
			}
		}
	}
	return nil
}

func (l *ENSIP15) checkFenced(cps []rune) error {
	name, ok := l.fenced[cps[0]]
	if ok {
		return fmt.Errorf("%w: %s", ErrFencedLeading, name)
	}
	n := len(cps)
	lastPos := -1
	var lastName string
	for i := 1; i < n; i++ {
		name, ok := l.fenced[cps[i]]
		if ok {
			if lastPos == i {
				return fmt.Errorf("%w: %s + %s", ErrFencedAdjacent, lastName, name)
			}
			lastPos = i + 1
			lastName = name
		}
	}
	if lastPos == n {
		return fmt.Errorf("%w: %s", ErrFencedTrailing, lastName)
	}
	return nil
}

func (l *ENSIP15) checkValidLabel(cps []rune, tokens []OutputToken) (*Group, error) {
	if len(cps) == 0 {
		return nil, ErrEmptyLabel
	}
	if err := checkLeadingUnderscore(cps); err != nil {
		return nil, err
	}
	hasEmoji := len(tokens) > 1 || tokens[0].Emoji != nil
	if !hasEmoji && isASCII(cps) {
		if err := checkLabelExtension(cps); err != nil {
			return nil, err
		}
		return l._ASCII, nil
	}
	chars := make([]rune, 0, len(cps))
	for _, t := range tokens {
		if t.Emoji == nil {
			chars = append(chars, t.Codepoints...)
		}
	}
	if hasEmoji && len(chars) == 0 {
		return l._EMOJI, nil
	}
	if err := l.checkCombiningMarks(tokens); err != nil {
		return nil, err
	}
	if err := l.checkFenced(cps); err != nil {
		return nil, err
	}
	unique := uniqueRunes(chars)
	group, err := l.determineGroup(unique)
	if err != nil {
		return nil, err
	}
	if err := l.checkGroup(group, chars); err != nil {
		return nil, err
	}
	if err := l.checkWhole(group, unique); err != nil {
		return nil, err
	}
	return group, nil
}
