package isbinary

import (
	"bytes"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func toUTF8(s string) []byte {
	out := []byte{}
	for _, char := range s {
		var curr [utf8.UTFMax]byte
		n := utf8.EncodeRune(curr[:], char)
		out = append(out, curr[:n]...)
	}

	return out
}

var (
	testCases = []struct {
		Input  []byte
		Result bool
	}{
		{[]byte("some text"), false},
		{[]byte("some text with a \x00 char"), true},
		{[]byte("\xEF\xBB\xBF text with utf-8 BOM"), false},
		{[]byte("text with suspicious \xFF\xFF\xFF\xFF\xFF\xFF\xFF"), true},
		{toUTF8("utf8 text  世界"), false},
		{toUTF8("utf8 世界 with null \x00"), true},
		{toUTF8("utf8 世界 with suspicious \x01\x01\x01\x01\x01"), true},
	}
)

func TestIsBinary(t *testing.T) {
	for i, tcase := range testCases {
		assert.Equal(t, tcase.Result, Test(tcase.Input),
			"test case %d did not have expected result", i)
	}
}

func TestIsBinaryReader(t *testing.T) {
	for i, tcase := range testCases {
		r := bytes.NewReader(tcase.Input)
		res, err := TestReader(r)
		assert.NoError(t, err)
		assert.Equal(t, tcase.Result, res,
			"test case %d did not have expected result", i)
	}
}
