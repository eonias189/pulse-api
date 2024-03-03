package utils

func Map[S ~[]T, T any, R any](s S, m func(T) R) []R {
	res := make([]R, len(s))
	for i, t := range s {
		res[i] = m(t)
	}
	return res
}

func Filter[S ~[]T, T any](s S, check func(T) bool) []T {
	res := []T{}
	for _, i := range s {
		if check(i) {
			res = append(res, i)
		}
	}
	return res
}

func All[S ~[]T, T any](s S, check func(T) bool) bool {
	for _, i := range s {
		if !check(i) {
			return false
		}
	}
	return true
}
