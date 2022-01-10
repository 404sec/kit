package typex

import "fmt"

type Float64 float64

func (n Float64) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.2f", n)), nil
}

type Float32 float64

func (n Float32) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.2f", n)), nil
}
