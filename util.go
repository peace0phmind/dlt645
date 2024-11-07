package dlt645

import (
	"fmt"
	"math"
)

func decimalDigits(value uint64) int {
	if value == 0 {
		return 1 // 特殊情况，0的位数是1
	}

	return int(math.Floor(math.Log10(float64(value)))) + 1
}

func uintToBcd(value uint64, bufSize int) []byte {
	valueDigits := decimalDigits(value)
	if bufSize*2 < valueDigits {
		panic(fmt.Errorf("buffer size is less than %d bytes", valueDigits))
	}

	buf := make([]byte, bufSize)

	for i := 0; i < valueDigits; i++ {
		nibble := byte(value % 10)
		value = value / 10
		if i%2 != 0 {
			buf[i/2] |= nibble << 4
		} else {
			buf[i/2] |= nibble
		}
	}

	return buf
}

func bcdToUint(data []byte, len int) (ret uint64) {
	base := uint64(0)

	for i := len*2 - 1; i >= 0; i-- {
		nibble := byte(0)
		if i%2 == 0 {
			nibble = data[i/2] & 0x0F
		} else {
			nibble = (data[i/2] >> 4) & 0x0F
		}

		// fix address compression prefix
		if nibble == 0x0A {
			nibble = 0
		}

		if base > 0 || nibble != 0 {
			ret = ret*base + uint64(nibble)
			if base == 0 {
				base = 10
			}
		}
	}

	return ret
}
