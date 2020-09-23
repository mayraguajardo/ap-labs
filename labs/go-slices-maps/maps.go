package main

import (
	"strings"

	"golang.org/x/tour/wc"
)

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	mapa := make(map[string]int)
	for i := range words{
		mapa[words[i]]++
	}
	return mapa
}

func main() {
	wc.Test(WordCount)
}