// internal/service/service.go
package service

import (
	"strings"
	"unicode"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

// Определяем возможные символы для детектирования
var morseSymbols = map[rune]bool{
	'.': true,
	'-': true,
	' ': true,
	'/': true,
}

// DetectAndConvert определяет, является ли входная строка текстом или кодом Морзе,
// и выполняет соответствующую конвертацию
func DetectAndConvert(input string) (string, error) {
	// Удаляем лишние пробелы по краям
	input = strings.TrimSpace(input)
	
	if input == "" {
		return "", nil
	}

	// Проверяем, является ли строка кодом Морзе
	if IsMorseCode(input) {
		// Конвертируем Морзе в текст
		return morse.ToText(input), nil
	}
	
	// Если это не Морзе, считаем что это текст и конвертируем в Морзе
	return morse.ToMorse(input), nil
}

// IsMorseCode проверяет, содержит ли строка только символы, характерные для кода Морзе
func IsMorseCode(s string) bool {
	// Если строка пустая, это не Морзе
	if len(s) == 0 {
		return false
	}

	// Проверяем каждый символ
	for _, r := range s {
		// Если символ не является допустимым символом Морзе и не пробел, это текст
		if !morseSymbols[r] && !unicode.IsSpace(r) {
			return false
		}
	}
	
	// Дополнительная проверка: если в строке есть хотя бы одна точка или тире,
	// и она не содержит букв/цифр
	hasMorseSymbols := strings.ContainsAny(s, ".-")
	
	// Проверяем наличие букв/цифр
	hasLetters := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			hasLetters = true
			break
		}
	}
	
	// Строка считается Морзе, если содержит символы Морзе и не содержит букв/цифр
	return hasMorseSymbols && !hasLetters
}
