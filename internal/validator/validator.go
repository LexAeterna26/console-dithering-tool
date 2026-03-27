package validator

import (
	"strings"
	"unicode"
)

const (
	allowedSymbols = "_-."
	thresholdMin   = 0
	thresholdMax   = 255
	matrixSize2    = 2
	matrixSize4    = 4
	matrixSize8    = 8
	matrixSize16   = 16
)

func ValidateSuffix(suffix string) bool {
	for _, ch := range suffix {
		if !unicode.IsDigit(ch) && !unicode.IsLetter(ch) && !strings.ContainsRune(allowedSymbols, ch) {
			return false
		}
	}
	return true
}

func ValidateThreshold(threshold int) bool {
	if threshold > thresholdMax || threshold < thresholdMin {
		return false
	}
	return true
}

func ValidateMatrixSize(matrixSize int) bool {
	if !(matrixSize == matrixSize2 || matrixSize == matrixSize4 || matrixSize == matrixSize8 || matrixSize == matrixSize16) {
		return false
	}
	return true
}
