// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package common

import "strconv"

// Float64 is a float64 that marshals to decimal notation (e.g. 0.0000001) instead of scientific notation (e.g. 1e-7).
type Float64 float64

func (f Float64) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(f), 'f', -1, 64)), nil
}

func (f *Float64) UnmarshalJSON(data []byte) error {
	v, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return err
	}
	*f = Float64(v)
	return nil
}

// Return a pointer to the given Float64 value.
func Float64Ptr(f Float64) *Float64 {
	return &f
}

// Return a pointer to the given int value.
func IntPtr(i int) *int {
	return &i
}
