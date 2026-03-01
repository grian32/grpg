package util

func PopSlice[T any](slice []T) (T, []T) {
	elem := slice[len(slice)-1]
	return elem, slice[:len(slice)-1]
}
