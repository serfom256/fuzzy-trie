package core

import (
	"github.com/agnivade/levenshtein"
)

const (
	suffixRegex = '*'
)

func absInt(x int) int {
	y := x >> 31
	return (x ^ y) - y
}
func shouldCollectSuffix(toSearch string, pos int) bool {
	return pos < len(toSearch) && toSearch[pos] == suffixRegex
}

func checkHashExistence(data *SearchData, pos int, distance int, node *TNode) bool {
	hash := ((pos + 1) << data.Typos) | distance
	if data.canMoveNext(node, distance, hash) {
		return true
	}
	data.addNodeToCache(node, hash)
	return false
}

func isSame(s1 string, s2 string, distance int) bool {
	return levenshtein.ComputeDistance(s1, s2) <= distance
}
