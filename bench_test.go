package basen_test

import (
	"testing"

	"github.com/cmars/basen"
)

// These benchmarks mirror the ones in encoding/base64, and results
// should be comparable to those.

func BenchmarkBase58EncodeToString(b *testing.B) {
	data := make([]byte, 8192)
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		basen.Base58.EncodeToString(data)
	}
}

func BenchmarkBase58DecodeString(b *testing.B) {
	data := basen.Base58.EncodeToString(make([]byte, 8192))
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		basen.Base58.DecodeString(data)
	}
}

func BenchmarkBase62EncodeToString(b *testing.B) {
	data := make([]byte, 8192)
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		basen.Base62.EncodeToString(data)
	}
}

func BenchmarkBase62DecodeString(b *testing.B) {
	data := basen.Base62.EncodeToString(make([]byte, 8192))
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		basen.Base62.DecodeString(data)
	}
}
