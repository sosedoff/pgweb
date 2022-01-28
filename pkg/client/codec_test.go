package client

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetBinaryCodec(t *testing.T) {
	examples := []struct {
		input string
		err   error
	}{
		{input: CodecNone, err: nil},
		{input: CodecBase58, err: nil},
		{input: CodecBase64, err: nil},
		{input: CodecHex, err: nil},
		{input: "foobar", err: errors.New("invalid binary codec: foobar")},
	}

	for _, ex := range examples {
		t.Run(ex.input, func(t *testing.T) {
			val := BinaryCodec
			defer func() {
				BinaryCodec = val
			}()

			assert.Equal(t, ex.err, SetBinaryCodec(ex.input))
		})
	}
}

func Test_encodeBinaryData(t *testing.T) {
	examples := []struct {
		input    string
		expected string
		encoding string
	}{
		{input: "hello world", expected: "hello world", encoding: CodecNone},
		{input: "hello world", expected: "StV1DL6CwTryKyV", encoding: CodecBase58},
		{input: "hello world", expected: "aGVsbG8gd29ybGQ=", encoding: CodecBase64},
		{input: "hello world", expected: "68656c6c6f20776f726c64", encoding: CodecHex},
	}

	for _, ex := range examples {
		t.Run(ex.input, func(t *testing.T) {
			assert.Equal(t, ex.expected, encodeBinaryData([]byte(ex.input), ex.encoding))
		})
	}
}
