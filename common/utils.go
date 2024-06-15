package common

func RunesFromInts(v []int) []rune {
	runes := make([]rune, len(v))
	for i, x := range v {
		runes[i] = rune(x)
	}
	return runes;
}
