package util

func ClampMin(curr uint32, min uint32) uint32 {
	if curr < min {
		return min
	}

	return curr
}

func ClampMax(curr uint32, max uint32) uint32 {
	if curr > max {
		return max
	}

	return curr
}
