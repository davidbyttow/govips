package gimage

import (
	"io/ioutil"
	"math"
)

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func Round(x float64) float64 {
	v, frac := math.Modf(x)
	if x > 0.0 {
		if frac > 0.5 || (frac == 0.5 && uint64(v)%2 != 0) {
			v += 1.0
		}
	} else {
		if frac < -0.5 || (frac == -0.5 && uint64(v)%2 != 0) {
			v -= 1.0
		}
	}
	return v
}

func CopyBuffer(in []byte) []byte {
	out := make([]byte, len(in))
	copy(out, in)
	return out
}
