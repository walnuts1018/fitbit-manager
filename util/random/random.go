package random

import (
	"crypto/rand"
	"fmt"
)

const UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const LowerLetters = "abcdefghijklmnopqrstuvwxyz"
const Numbers = "0123456789"
const Symbols = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
const Alphabets = UpperLetters + LowerLetters
const Alphanumeric = Alphabets + Numbers
const AlphanumericSymbols = Alphanumeric + Symbols

var testValue []byte = nil

func SetTestValue(value []byte) {
	testValue = value
}

func randRead(b []byte) (int, error) {
	if testValue != nil {
		copy(b, testValue)
		return len(b), nil
	}
	return rand.Read(b)
}

func String(length uint, base string) (string, error) {
	b := make([]byte, length)
	if _, err := randRead(b); err != nil {
		return "", fmt.Errorf("failed to read random: %w", err)
	}

	var result string
	for _, v := range b {
		result += string(base[int(v)%len(base)])
	}
	return result, nil
}

func Byte(length int) ([]byte, error) {
	b := make([]byte, length)
	if _, err := randRead(b); err != nil {
		return nil, fmt.Errorf("failed to read random: %w", err)
	}
	return b, nil
}
