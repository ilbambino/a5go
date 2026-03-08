package core

import "strconv"

func panicString(message string) {
	panic(message)
}

func itoa(v int) string {
	return strconv.Itoa(v)
}

func uitoa(v uint64) string {
	return strconv.FormatUint(v, 10)
}
