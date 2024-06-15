package ensip15

import (
	"fmt"
)

type OutputToken struct {
	Codepoints []rune
	Emoji      *EmojiSequence
}

func (ot OutputToken) String() string {
	if ot.Emoji != nil {
		return fmt.Sprintf("Emoji[%s]", ToHexSequence(ot.Emoji.normalized))
	} else {
		return fmt.Sprintf("Text[%s]", ToHexSequence(ot.Codepoints))
	}
}

func FlattenTokens(tokens []OutputToken) []rune {
	n := 0
	for _, x := range tokens {
		n += len(x.Codepoints)
	}
	cps := make([]rune, 0, n)
	for _, x := range tokens {
		cps = append(cps, x.Codepoints...)
	}
	return cps
}

func (l *ENSIP15) outputTokenize(
	cps []rune,
	nf func([]rune) []rune,
	ef func(EmojiSequence) []rune,
) (tokens []OutputToken, err error) {
	var buf []rune
	for i := 0; i < len(cps); {
		emoji, end := l.ParseEmojiAt(cps, i)
		if emoji != nil {
			if len(buf) > 0 {
				tokens = append(tokens, OutputToken{
					Codepoints: nf(buf),
				})
				buf = nil
			}
			tokens = append(tokens, OutputToken{
				Codepoints: ef(*emoji),
				Emoji:      emoji,
			})
			i = end
		} else {
			cp := cps[i]
			if l.possiblyValid.Contains(cp) {
				buf = append(buf, cp)
			} else if mapped, ok := l.mapped[cp]; ok {
				buf = append(buf, mapped...)
			} else if !l.ignored.Contains(cp) {
				return nil, fmt.Errorf("%w: %s", ErrDisallowedCharacter, l.SafeCodepoint(cp))
			}
			i++
		}
	}
	if len(buf) > 0 {
		tokens = append(tokens, OutputToken{
			Codepoints: nf(buf),
		})
	}
	return tokens, nil
}
