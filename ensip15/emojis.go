package ensip15

import (
	"github.com/adraffy/go-ens-normalize/util"
)

const (
	FE0F = 0xFE0F
	ZWJ  = 0x200D
)

type EmojiSequence struct {
	normalized []rune
	beautified []rune
}

func (seq EmojiSequence) Normalized() string {
	return string(seq.normalized)
}
func (seq EmojiSequence) Beautified() string {
	return string(seq.beautified)
}
func (seq EmojiSequence) String() string {
	return seq.Beautified()
}
func (seq EmojiSequence) IsMangled() bool {
	return len(seq.normalized) < len(seq.beautified)
}
func (seq EmojiSequence) HasZWJ() bool {
	for _, x := range seq.beautified {
		if x == ZWJ {
			return true
		}
	}
	return false
}

func decodeEmojis(d *util.Decoder, prev []rune) (v []EmojiSequence) {
	for _, cp := range d.ReadSortedAscending(d.ReadUnsigned()) {
		beautified := make([]rune, 0, len(prev)+1)
		beautified = append(beautified, prev...)
		beautified = append(beautified, rune(cp))
		normalized := make([]rune, 0, len(beautified))
		for _, x := range beautified {
			if x != FE0F {
				normalized = append(normalized, x)
			}
		}
		if len(normalized) == len(beautified) {
			normalized = beautified
		}
		v = append(v, EmojiSequence{
			normalized,
			beautified,
		})
	}
	for _, cp := range d.ReadSortedAscending(d.ReadUnsigned()) {
		v = append(v, decodeEmojis(d, append(prev, rune(cp)))...)
	}
	return v
}

type EmojiNode struct {
	emoji    *EmojiSequence
	children map[rune]*EmojiNode
}

func (node *EmojiNode) Child(cp rune) *EmojiNode {
	if node.children == nil {
		node.children = make(map[rune]*EmojiNode)
	}
	child, ok := node.children[cp]
	if !ok {
		child = &EmojiNode{}
		node.children[cp] = child
	}
	return child
}

func makeEmojiTree(all []EmojiSequence) *EmojiNode {
	root := &EmojiNode{}
	for _, emoji := range all {
		v := []*EmojiNode{root}
		for _, cp := range emoji.beautified {
			if cp == FE0F {
				for _, node := range v {
					v = append(v, node.Child(cp))
				}
			} else {
				for i, node := range v {
					v[i] = node.Child(cp)
				}
			}
		}
		for _, node := range v {
			node.emoji = &emoji
		}
	}
	return root
}

func (l *ENSIP15) ParseEmojiAt(cps []rune, pos int) (emoji *EmojiSequence, end int) {
	end = -1
	node := l.emojiRoot
	for pos < len(cps) {
		if node.children == nil {
			break
		}
		node = node.children[cps[pos]]
		if node == nil {
			break
		}
		pos++
		if node.emoji != nil {
			emoji = node.emoji
			end = pos
		}
	}
	return emoji, end
}
