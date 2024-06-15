package ensip15

import (
	"fmt"

	"github.com/adraffy/go-ensnormalize/decoder"
)

const (
	FE0F = 0xFE0F
	ZWJ  = 0x200D
)

type EmojiSequence struct {
	normalized []rune
	beautified []rune
}

func NewEmojiSequence(cps []rune) EmojiSequence {
	v := make([]rune, 0, len(cps))
	for _, x := range cps {
		if x != FE0F {
			v = append(v, x)
		}
	}
	if len(v) == len(cps) {
		v = cps
	}
	return EmojiSequence{
		normalized: v,
		beautified: cps,
	}
}

func (seq EmojiSequence) String() string {
	return fmt.Sprint(seq.beautified)
	//return string(seq.beautified)
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

func decodeEmojis(d *decoder.Decoder, path []rune) (v []EmojiSequence) {
	for _, cp := range d.ReadSortedAscending(d.ReadUnsigned()) {
		v = append(v, NewEmojiSequence(append(path, rune(cp))))
	}
	for _, cp := range d.ReadSortedAscending(d.ReadUnsigned()) {
		v = append(v, decodeEmojis(d, append(path, rune(cp)))...)
	}
	return v
}

type EmojiNode struct {
	emoji    EmojiSequence
	children map[rune]EmojiNode
}

func (node *EmojiNode) Child(cp rune) EmojiNode {
	if node.children == nil {
		node.children = make(map[rune]EmojiNode)
	}
	child, ok := node.children[cp]
	if !ok {
		child = EmojiNode{}
		node.children[cp] = child
	}
	return child
}

func makeEmojiTree(v []EmojiSequence) EmojiNode {
	root := EmojiNode{}
	for _, emoji := range v {
		var v []EmojiNode
		v = append(v, root)
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
			node.emoji = emoji
		}
	}
	return root
}
