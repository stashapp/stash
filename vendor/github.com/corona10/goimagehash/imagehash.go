// Copyright 2017 The goimagehash Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goimagehash

import (
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// Kind describes the kinds of hash.
type Kind int

// ImageHash is a struct of hash computation.
type ImageHash struct {
	hash uint64
	kind Kind
}

// ExtImageHash is a struct of big hash computation.
type ExtImageHash struct {
	hash []uint64
	kind Kind
	bits int
}

const (
	// Unknown is a enum value of the unknown hash.
	Unknown Kind = iota
	// AHash is a enum value of the average hash.
	AHash
	//PHash is a enum value of the perceptual hash.
	PHash
	// DHash is a enum value of the difference hash.
	DHash
	// WHash is a enum value of the wavelet hash.
	WHash
)

// NewImageHash function creates a new image hash.
func NewImageHash(hash uint64, kind Kind) *ImageHash {
	return &ImageHash{hash: hash, kind: kind}
}

// Bits method returns an actual hash bit size
func (h *ImageHash) Bits() int {
	return 64
}

// Distance method returns a distance between two hashes.
func (h *ImageHash) Distance(other *ImageHash) (int, error) {
	if h.GetKind() != other.GetKind() {
		return -1, errors.New("Image hashes's kind should be identical")
	}

	lhash := h.GetHash()
	rhash := other.GetHash()

	hamming := lhash ^ rhash
	return popcnt(hamming), nil
}

// GetHash method returns a 64bits hash value.
func (h *ImageHash) GetHash() uint64 {
	return h.hash
}

// GetKind method returns a kind of image hash.
func (h *ImageHash) GetKind() Kind {
	return h.kind
}

func (h *ImageHash) leftShiftSet(idx int) {
	h.hash |= 1 << uint(idx)
}

const strFmt = "%1s:%016x"

// Dump method writes a binary serialization into w io.Writer.
func (h *ImageHash) Dump(w io.Writer) error {
	type D struct {
		Hash uint64
		Kind Kind
	}
	enc := gob.NewEncoder(w)
	err := enc.Encode(D{Hash: h.hash, Kind: h.kind})
	if err != nil {
		return err
	}
	return nil
}

// LoadImageHash method loads a ImageHash from io.Reader.
func LoadImageHash(b io.Reader) (*ImageHash, error) {
	type E struct {
		Hash uint64
		Kind Kind
	}
	var e E
	dec := gob.NewDecoder(b)
	err := dec.Decode(&e)
	if err != nil {
		return nil, err
	}
	return &ImageHash{hash: e.Hash, kind: e.Kind}, nil
}

// ImageHashFromString returns an image hash from a hex representation
//
// Deprecated: Use goimagehash.LoadImageHash instead.
func ImageHashFromString(s string) (*ImageHash, error) {
	var kindStr string
	var hash uint64
	_, err := fmt.Sscanf(s, strFmt, &kindStr, &hash)
	if err != nil {
		return nil, errors.New("Couldn't parse string " + s)
	}

	kind := Unknown
	switch kindStr {
	case "a":
		kind = AHash
	case "p":
		kind = PHash
	case "d":
		kind = DHash
	case "w":
		kind = WHash
	}
	return NewImageHash(hash, kind), nil
}

// ToString returns a hex representation of the hash
func (h *ImageHash) ToString() string {
	kindStr := ""
	switch h.kind {
	case AHash:
		kindStr = "a"
	case PHash:
		kindStr = "p"
	case DHash:
		kindStr = "d"
	case WHash:
		kindStr = "w"
	}
	return fmt.Sprintf(strFmt, kindStr, h.hash)
}

// NewExtImageHash function creates a new big hash
func NewExtImageHash(hash []uint64, kind Kind, bits int) *ExtImageHash {
	return &ExtImageHash{hash: hash, kind: kind, bits: bits}
}

// Bits method returns an actual hash bit size
func (h *ExtImageHash) Bits() int {
	return h.bits
}

// Distance method returns a distance between two big hashes
func (h *ExtImageHash) Distance(other *ExtImageHash) (int, error) {
	if h.GetKind() != other.GetKind() {
		return -1, errors.New("Extended Image hashes's kind should be identical")
	}

	if h.Bits() != other.Bits() {
		msg := fmt.Sprintf("Extended image hash should has an identical bit size but got %v vs %v", h.Bits(), other.Bits())
		return -1, errors.New(msg)
	}

	lHash := h.GetHash()
	rHash := other.GetHash()
	if len(lHash) != len(rHash) {
		return -1, errors.New("Extended Image hashes's size should be identical")
	}

	distance := 0
	for idx, lh := range lHash {
		rh := rHash[idx]
		hamming := lh ^ rh
		distance += popcnt(hamming)
	}
	return distance, nil
}

// GetHash method returns a big hash value
func (h *ExtImageHash) GetHash() []uint64 {
	return h.hash
}

// GetKind method returns a kind of big hash
func (h *ExtImageHash) GetKind() Kind {
	return h.kind
}

// Dump method writes a binary serialization into w io.Writer.
func (h *ExtImageHash) Dump(w io.Writer) error {
	type D struct {
		Hash []uint64
		Kind Kind
		Bits int
	}
	enc := gob.NewEncoder(w)
	err := enc.Encode(D{Hash: h.hash, Kind: h.kind, Bits: h.bits})
	if err != nil {
		return err
	}
	return nil
}

// LoadExtImageHash method loads a ExtImageHash from io.Reader.
func LoadExtImageHash(b io.Reader) (*ExtImageHash, error) {
	type E struct {
		Hash []uint64
		Kind Kind
		Bits int
	}
	var e E
	dec := gob.NewDecoder(b)
	err := dec.Decode(&e)
	if err != nil {
		return nil, err
	}
	return &ExtImageHash{hash: e.Hash, kind: e.Kind, bits: e.Bits}, nil
}

const extStrFmt = "%1s:%s"

// ExtImageHashFromString returns a big hash from a hex representation
//
// Deprecated: Use goimagehash.LoadExtImageHash instead.
func ExtImageHashFromString(s string) (*ExtImageHash, error) {
	var kindStr string
	var hashStr string
	_, err := fmt.Sscanf(s, extStrFmt, &kindStr, &hashStr)
	if err != nil {
		return nil, errors.New("Couldn't parse string " + s)
	}

	hexBytes, err := hex.DecodeString(hashStr)
	if err != nil {
		return nil, err
	}

	var hash []uint64
	lenOfByte := 8
	for i := 0; i < len(hexBytes)/lenOfByte; i++ {
		startIndex := i * lenOfByte
		endIndex := startIndex + lenOfByte
		hashUint64 := binary.BigEndian.Uint64(hexBytes[startIndex:endIndex])
		hash = append(hash, hashUint64)
	}

	kind := Unknown
	switch kindStr {
	case "a":
		kind = AHash
	case "p":
		kind = PHash
	case "d":
		kind = DHash
	case "w":
		kind = WHash
	}
	return NewExtImageHash(hash, kind, len(hash)*64), nil
}

// ToString returns a hex representation of big hash
func (h *ExtImageHash) ToString() string {
	var hexBytes []byte
	for _, hash := range h.hash {
		hashBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(hashBytes, hash)
		hexBytes = append(hexBytes, hashBytes...)
	}
	hexStr := hex.EncodeToString(hexBytes)

	kindStr := ""
	switch h.kind {
	case AHash:
		kindStr = "a"
	case PHash:
		kindStr = "p"
	case DHash:
		kindStr = "d"
	case WHash:
		kindStr = "w"
	}
	return fmt.Sprintf(extStrFmt, kindStr, hexStr)
}
