package core

import "strconv"

func HexToU64(hex string) (uint64, error) {
	return strconv.ParseUint(hex, 16, 64)
}

func U64ToHex(index uint64) string {
	return strconv.FormatUint(index, 16)
}
