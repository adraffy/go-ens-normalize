package ensip15

import "fmt"

var (
	ErrInvalidLabelExtension = fmt.Errorf("invalid label extension")
	ErrIllegalMixture        = fmt.Errorf("illegal mixture")
	ErrWholeConfusable       = fmt.Errorf("whole-script confusable")
	ErrLeadingUnderscore     = fmt.Errorf("underscore allowed only at start")
	ErrFencedLeading         = fmt.Errorf("leading fenced")
	ErrFencedAdjacent        = fmt.Errorf("adjacent fenced")
	ErrFencedTrailing        = fmt.Errorf("trailing fenced")
	ErrDisallowedCharacter   = fmt.Errorf("disallowed character")
	ErrEmptyLabel            = fmt.Errorf("empty label")
	ErrCMLeading             = fmt.Errorf("leading combining mark")
	ErrCMAfterEmoji          = fmt.Errorf("emoji + combining mark")
	ErrNSMDuplicate          = fmt.Errorf("duplicate non-spacing marks")
	ErrNSMExcessive          = fmt.Errorf("excessive non-spacing marks")
)

func (l *ENSIP15) createMixtureError(group *Group, cp rune) error {
	conflict := l.SafeCodepoint(cp)
	var other *Group
	for _, g := range l.groups {
		if g.primary.Contains(cp) {
			other = g
			break
		}
	}
	if other != nil {
		conflict = fmt.Sprintf("%s %s", other, conflict)
	}
	return fmt.Errorf("%w: %s + %s", ErrIllegalMixture, group, conflict)
}
