package main

import (
	"strings"
)

func countVowels(vowels string) int {
	var (
		a = 'a'
		e = 'e'
		i = 'i'
		o = 'o'
		u = 'u'
	)

	counter := 0

	vowels = strings.ToLower(vowels)

	for _, vowel := range vowels {
		if vowel == a || vowel == e || vowel == i || vowel == o || vowel == u {
			counter++
		}
	}

	return counter
}
