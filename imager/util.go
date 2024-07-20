package imager

func mapPrecedence[T comparable](s []T) map[T]int {
	m := map[T]int{}
	for i := range s {
		m[s[i]] = i * 10
	}
	return m
}

func sliceJoin[T any](s ...[]T) []T {
	n := 0
	for _, s := range s {
		n += len(s)
	}
	slc := make([]T, 0, n)
	for _, s := range s {
		slc = append(slc, s...)
	}
	return slc
}
