// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package ranger

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteRanger(t *testing.T) {
	for _, example := range []struct {
		data                 string
		size, offset, length int64
		substr               string
		fail                 bool
	}{
		{"", 0, 0, 0, "", false},
		{"abcdef", 6, 0, 0, "", false},
		{"abcdef", 6, 3, 0, "", false},
		{"abcdef", 6, 0, 6, "abcdef", false},
		{"abcdef", 6, 0, 5, "abcde", false},
		{"abcdef", 6, 0, 4, "abcd", false},
		{"abcdef", 6, 1, 4, "bcde", false},
		{"abcdef", 6, 2, 4, "cdef", false},
		{"abcdefg", 7, 1, 4, "bcde", false},
		{"abcdef", 6, 0, 7, "", true},
		{"abcdef", 6, -1, 7, "abcde", true},
		{"abcdef", 6, 0, -1, "abcde", true},
	} {
		rr := ByteRanger([]byte(example.data))
		if rr.Size() != example.size {
			t.Fatalf("invalid size: %v != %v", rr.Size(), example.size)
		}
		r, err := rr.Range(context.Background(), example.offset, example.length)
		if example.fail {
			if err == nil {
				t.Fatalf("expected error")
			}
			continue
		}
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		data, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if !bytes.Equal(data, []byte(example.substr)) {
			t.Fatalf("invalid subrange: %#v != %#v", string(data), example.substr)
		}
	}
}

func TestConcatReader(t *testing.T) {
	for _, example := range []struct {
		data                 []string
		size, offset, length int64
		substr               string
	}{
		{[]string{}, 0, 0, 0, ""},
		{[]string{""}, 0, 0, 0, ""},
		{[]string{"abcdefghijkl"}, 12, 1, 4, "bcde"},
		{[]string{"abcdef", "ghijkl"}, 12, 1, 4, "bcde"},
		{[]string{"abcdef", "ghijkl"}, 12, 1, 5, "bcdef"},
		{[]string{"abcdef", "ghijkl"}, 12, 1, 6, "bcdefg"},
		{[]string{"abcdef", "ghijkl"}, 12, 5, 4, "fghi"},
		{[]string{"abcdef", "ghijkl"}, 12, 6, 4, "ghij"},
		{[]string{"abcdef", "ghijkl"}, 12, 7, 4, "hijk"},
		{[]string{"abcdef", "ghijkl"}, 12, 7, 5, "hijkl"},
		{[]string{"abcdef", "ghijkl", "mnopqr"}, 18, 7, 7, "hijklmn"},
		{[]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"},
			12, 7, 3, "hij"},
	} {
		var readers []RangeCloser
		for _, data := range example.data {
			readers = append(readers, ByteRanger([]byte(data)))
		}
		rr := Concat(readers...)
		if rr.Size() != example.size {
			t.Fatalf("invalid size: %v != %v", rr.Size(), example.size)
		}
		r, err := rr.Range(context.Background(), example.offset, example.length)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		data, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if !bytes.Equal(data, []byte(example.substr)) {
			t.Fatalf("invalid subrange: %#v != %#v", string(data), example.substr)
		}
	}
}

func TestSubranger(t *testing.T) {
	for _, example := range []struct {
		data             string
		offset1, length1 int64
		offset2, length2 int64
		substr           string
	}{
		{"abcdefghijkl", 0, 4, 0, 4, "abcd"},
		{"abcdefghijkl", 0, 4, 0, 3, "abc"},
		{"abcdefghijkl", 0, 4, 1, 3, "bcd"},
		{"abcdefghijkl", 1, 4, 0, 4, "bcde"},
		{"abcdefghijkl", 1, 4, 0, 3, "bcd"},
		{"abcdefghijkl", 1, 4, 1, 3, "cde"},
		{"abcdefghijkl", 8, 4, 0, 4, "ijkl"},
		{"abcdefghijkl", 8, 4, 0, 3, "ijk"},
		{"abcdefghijkl", 8, 4, 1, 3, "jkl"},
	} {
		rr, err := Subrange(ByteRangeCloser([]byte(example.data)),
			example.offset1, example.length1)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if rr.Size() != example.length1 {
			t.Fatalf("invalid size: %v != %v", rr.Size(), example.length1)
		}
		r, err := rr.Range(context.Background(), example.offset2, example.length2)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		data, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if !bytes.Equal(data, []byte(example.substr)) {
			t.Fatalf("invalid subrange: %#v != %#v", string(data), example.substr)
		}
	}
}

func TestSubrangerError(t *testing.T) {
	for _, tt := range []struct {
		name           string
		data           string
		offset, length int64
	}{
		{name: "Negative offset", data: "abcd", offset: -1},
		{name: "Offset is bigger than DataSize", data: "abcd", offset: 5},
		{name: "Length and offset is bigger than DataSize", data: "abcd", offset: 4, length: 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rr, err := Subrange(ByteRangeCloser([]byte(tt.data)), tt.offset, tt.length)
			assert.Nil(t, rr)
			assert.NotNil(t, err)
		})

	}
}
