package service

import (
	"strings"
	"unicode"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)


var morseSymbols = map[rune]bool{
	'.': true,
	'-': true,
	' ': true,
	'/': true,
}


func DetectAndConvert(input string) (string, error) {
	
	input = strings.TrimSpace(input)

	if input == "" {
		return "", nil
	}

	if IsMorseCode(input) {
		
		return morse.ToText(input), nil
	}

	
	return morse.ToMorse(input), nil
}


func IsMorseCode(s string) bool {
	
	if len(s) == 0 {
		return false
	}

	
	for _, r := range s {
		
		if !morseSymbols[r] && !unicode.IsSpace(r) {
			return false
		}
	}

	
	hasMorseSymbols := strings.ContainsAny(s, ".-")

	hasLetters := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			hasLetters = true
			break
		}
	}

	return hasMorseSymbols && !hasLetters
}
