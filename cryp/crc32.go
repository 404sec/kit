package cryp

import (
	"hash/crc32"
	"strconv"
)

func CRC32(str string) string {
	return strconv.FormatUint(uint64(crc32.ChecksumIEEE([]byte(str))), 10)
}
