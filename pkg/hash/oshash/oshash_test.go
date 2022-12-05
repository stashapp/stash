package oshash

import (
	"bytes"
	"math/rand"
	"testing"
)

func BenchmarkOsHash(b *testing.B) {
	src := rand.NewSource(9999)
	r := rand.New(src)

	size := int64(1234567890)

	head := make([]byte, 1024*64)
	_, err := r.Read(head)
	if err != nil {
		b.Errorf("unable to generate head array: %v", err)
	}

	tail := make([]byte, 1024*64)
	_, err = r.Read(tail)
	if err != nil {
		b.Errorf("unable to generate tail array: %v", err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, err := oshash(size, head, tail)
		if err != nil {
			b.Errorf("unexpected error: %v", err)
		}
	}
}

func TestFromReader(t *testing.T) {
	makeByteArray := func(base []byte, mag int) []byte {
		ret := base
		for i := 0; i < mag; i++ {
			ret = append(ret, ret...)
		}
		return ret
	}

	makeTailArray := func(base []byte, tail []byte) []byte {
		ret := base
		t := make([]byte, chunkSize)
		copy(t[len(t)-len(tail):], tail)
		ret = append(ret, t...)
		return ret
	}

	tests := []struct {
		name    string
		data    []byte
		want    string
		wantErr bool
	}{
		{
			"empty",
			[]byte{},
			"",
			true,
		},
		{
			"regular",
			makeByteArray([]byte("this is a test"), 15),
			"6a0eba04654d0b9b",
			false,
		},
		{
			"< chunk size",
			[]byte("hello world"),
			"d3e392dee38cd4df",
			false,
		},
		{
			"< 8",
			[]byte("hello"),
			"",
			true,
		},
		{
			"identical #1",
			makeTailArray(make([]byte, chunkSize), []byte("this is dumb")),
			"d5d6ddd820756920",
			false,
		},
		{
			"identical #2",
			makeTailArray(make([]byte, chunkSize), []byte("dumb is this")),
			"d5d6ddd820756920",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.data)

			got, err := FromReader(r, int64(len(tt.data)))
			if (err != nil) != tt.wantErr {
				t.Errorf("FromReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FromReader() = %v, want %v", got, tt.want)
			}
		})
	}
}
