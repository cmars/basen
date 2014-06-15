// Copyright (c) 2014 Casey Marshall. See LICENSE file for details.

package basen

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"unicode/utf8"
)

// Encoding represents a given base-N encoding.
type Encoding struct {
	alphabet string
	index    map[byte]int
	base     *big.Int
}

const stdBase62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// StdBase62 represents bytes as a base-62 number (0-9A-Za-z).
var StdBase62 = NewEncoding(stdBase62Alphabet)

// NewEncoding creates a new base-N representation from the given alphabet.
// Panics if the alphabet is not unique. Only ASCII characters are supported.
func NewEncoding(alphabet string) *Encoding {
	return &Encoding{
		alphabet: alphabet,
		index:    newAlphabetMap(alphabet),
		base:     big.NewInt(int64(len(alphabet))),
	}
}

func newAlphabetMap(s string) map[byte]int {
	if utf8.RuneCountInString(s) != len(s) {
		panic("multi-byte characters not supported")
	}
	result := make(map[byte]int)
	for i := range s {
		result[s[i]] = i
	}
	if len(result) != len(s) {
		panic("alphabet contains non-unique characters")
	}
	return result
}

var zero = big.NewInt(int64(0))

// Random returns the base-encoded representation of n random bytes.
func (enc *Encoding) Random(n int) (string, error) {
	buf := make([]byte, n)
	_, err := rand.Reader.Read(buf)
	if err != nil {
		return "", err
	}
	return enc.EncodeToString(buf), nil
}

// MustRandom returns the base-encoded representation of n random bytes,
// panicking in the unlikely event of a read error from the random source.
func (enc *Encoding) MustRandom(n int) string {
	s, err := enc.Random(n)
	if err != nil {
		panic(err)
	}
	return s
}

// Base returns the number base of the encoding.
func (enc *Encoding) Base() int {
	return len(enc.alphabet)
}

// EncodeToString returns the base-encoded string representation
// of the given bytes.
func (enc *Encoding) EncodeToString(b []byte) string {
	n := new(big.Int)
	r := new(big.Int)
	n.SetBytes(b)
	var result []byte
	for n.Cmp(zero) > 0 {
		n, r = n.DivMod(n, enc.base, r)
		result = append([]byte{enc.alphabet[r.Int64()]}, result...)
	}
	return string(result)
}

// DecodeString returns the bytes for the given base-encoded string.
func (enc *Encoding) DecodeString(s string) ([]byte, error) {
	result := big.NewInt(0)
	for i := range s {
		n, ok := enc.index[s[i]]
		if !ok {
			return nil, fmt.Errorf("invalid character %q at index %d", s[i], i)
		}
		result = big.NewInt(0).Add(big.NewInt(0).Mul(result, enc.base), big.NewInt(int64(n)))
	}
	return result.Bytes(), nil
}
