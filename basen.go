// Copyright (c) 2014 Casey Marshall. See LICENSE file for details.

package basen

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

type Encoding struct {
	alphabet string
	index    map[byte]int
	base     *big.Int
}

const stdBase62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var StdBase62 = NewEncoding(stdBase62Alphabet)

func NewEncoding(alphabet string) *Encoding {
	return &Encoding{
		alphabet: alphabet,
		index:    newAlphabetMap(alphabet),
		base:     big.NewInt(int64(len(alphabet))),
	}
}

func newAlphabetMap(s string) map[byte]int {
	result := make(map[byte]int)
	for i := range s {
		result[s[i]] = i
	}
	return result
}

var zero = big.NewInt(int64(0))

func (enc *Encoding) Random(r io.Reader, n int) (string, error) {
	if r == nil {
		r = rand.Reader
	}
	buf := make([]byte, n)
	_, err := r.Read(buf)
	if err != nil {
		return "", err
	}
	return enc.EncodeToString(buf), nil
}

func (enc *Encoding) MustRandom(r io.Reader, n int) string {
	s, err := enc.Random(r, n)
	if err != nil {
		panic(err)
	}
	return s
}

func (enc *Encoding) Base() int {
	return len(enc.alphabet)
}

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
