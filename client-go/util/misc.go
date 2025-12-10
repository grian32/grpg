package util

func MinI(val, min int32) int32 {
	if val > min {
		return val
	} else {
		return min
	}
}

func StringPtr(s string) *string {
	return &s
}
