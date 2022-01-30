package protocol

import (
	"encoding/binary"
)

func IntToBytes(value int) []byte {
	return Int64toBytes(int64(value))
}

func Int64toBytes(value int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(value))
	return b
}

func BytesArrayToUint64(bytes []byte) int64 {
	if len(bytes) < 8 {
		return 0
	}
	return int64(binary.LittleEndian.Uint64(bytes))
}
