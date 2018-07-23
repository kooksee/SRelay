package types

import (
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func BytesTrimSpace(bs []byte) []byte {
	for i, b := range bs {
		if b != 0 {
			bs = bs[i:]
			break
		}
	}

	for i, b := range bs {
		if b == 0 {
			bs = bs[:i]
			break
		}
	}

	return bs
}
