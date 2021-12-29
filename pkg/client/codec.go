package client

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

const (
	CodecNone   = "none"
	CodecHex    = "hex"
	CodecBase64 = "base64"
)

var (
	// BinaryEncodingFormat specifies the default serialization format of binary data
	BinaryCodec = CodecBase64
)

func SetBinaryCodec(codec string) error {
	switch codec {
	case CodecNone, CodecHex, CodecBase64:
		BinaryCodec = codec
	default:
		return fmt.Errorf("invalid binary codec: %v", codec)
	}

	return nil
}

func encodeBinaryData(data []byte) string {
	switch BinaryCodec {
	case CodecHex:
		return hex.EncodeToString(data)
	case CodecBase64:
		return base64.StdEncoding.EncodeToString(data)
	default:
		return string(data)
	}
}
