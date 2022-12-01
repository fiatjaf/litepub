package litepub

func mapSlice[B any, A any](slice []B, fn func(item B) A) []A {
	newSlice := make([]A, len(slice))
	for i, item := range slice {
		newSlice[i] = fn(item)
	}
	return newSlice
}
