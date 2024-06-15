package ensip15

import (
	"fmt"
	"strings"
)

func Join(labels []string) string {
	return strings.Join(labels, ".")
}

func Split(name string) []string {
	if len(name) == 0 {
		return nil // empty name allowance
	}
	return strings.Split(name, ".")
}

func ToHexSequence(cps []rune) string {
	var sb strings.Builder
	for i, cp := range cps {
		if i > 0 {
			sb.WriteRune(' ')
		}
		appendHex(&sb, cp)
	}
	return sb.String()
}

func appendHex(sb *strings.Builder, cp rune) {
	sb.WriteString(fmt.Sprintf("%02X", cp))
}

func appendHexEscape(sb *strings.Builder, cp rune) {
	sb.WriteRune('{')
	appendHex(sb, cp)
	sb.WriteRune('}')
}

func isASCII(cps []rune) bool {
	for _, cp := range cps {
		if cp >= 0x80 {
			return false
		}
	}
	return true
}

func uniqueRunes(cps []rune) []rune {
	set := make(map[rune]bool)
	v := make([]rune, 0, len(cps))
	for _, cp := range cps {
		if !set[cp] {
			set[cp] = true
			v = append(v, cp)
		}
	}
	return v
}

func compareRunes(a, b []rune) int {
	c := len(a) - len(b)
	if c != 0 {
		return c
	}
	for i, aa := range a {
		switch {
		case aa < b[i]:
			return -1
		case aa > b[i]:
			return 1
		}
	}
	return 0
}

func (l *ENSIP15) SafeCodepoint(cp rune) string {
	var sb strings.Builder
	if !l.shouldEscape.Contains(cp) {
		sb.WriteRune('"')
		l.safeImplode(&sb, []rune{cp})
		sb.WriteRune('"')
		sb.WriteRune(' ')
	}
	appendHexEscape(&sb, cp)
	return sb.String()
}

func (l *ENSIP15) safeImplode(sb *strings.Builder, cps []rune) {
	if len(cps) == 0 {
		return
	}
	if l.combiningMarks.Contains(cps[0]) {
		sb.WriteRune(0x25CC)
	}
	for _, cp := range cps {
		if l.shouldEscape.Contains(cp) {
			appendHexEscape(sb, cp)
		} else {
			sb.WriteRune(cp)
		}
	}
	// some messages can be mixed-directional and result in spillover
	// use 200E after a input string to reset the bidi direction
	// https://www.w3.org/International/questions/qa-bidi-unicode-controls#exceptions
	sb.WriteRune(0x200E)
}

func (l *ENSIP15) SafeImplode(cps []rune) string {
	var sb strings.Builder
	l.safeImplode(&sb, cps)
	return sb.String()
}
