package spongebob

import "unicode"

const (
	firstCapital = 65
	lastCapital  = 90
	firstLower   = 97
	lastLower    = 122
	caseDelta    = firstCapital - firstLower
)

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
		if v > unicode.MaxASCII {
			return s
		}
		if isLower(v) || isUpper(v) {
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