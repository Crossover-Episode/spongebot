package spongebob

import (
	"strings"
	"unicode"
)

const (
	firstCapital = 65
	lastCapital  = 90
	firstLower   = 97
	lastLower    = 122
	caseDelta    = firstCapital - firstLower
)

//StripText removes tabs, and line breaks
func StripText(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	return s
}

func ToText(s string, startCapital bool) string {
	isLower := func(char int32) bool {
		return char >= firstLower && char <= lastLower
	}

	isUpper := func(char int32) bool {
		return char >= firstCapital && char <= lastCapital
	}

	toCapitol := startCapital
	var sbString string
	for _, v := range s {
		if (isLower(v) || isUpper(v)) && v <= unicode.MaxASCII {
			if toCapitol && isLower(v) {
				sbString += string(v + caseDelta)
			} else if !toCapitol && isUpper(v) {
				sbString += string(v - caseDelta)
			} else {
				sbString += string(v)
			}
			toCapitol = !toCapitol
		} else {
			sbString += string(v)
		}
	}
	return sbString
}
