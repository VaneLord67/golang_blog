package common

import "strconv"

func FromUint64ToStr(input uint64) string {
	return strconv.FormatUint(input, 10)
}
