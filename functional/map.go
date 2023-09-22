package functional

type Pairs[K comparable, V any] interface {
	Len() int
	Get(int) (K, V)
}

func Tabulate[K comparable, V any](s Pairs[K, V]) map[K]V {
	r := make(map[K]V, s.Len())
	for i := 0; i < s.Len(); i++ {
		k, v := s.Get(i)
		r[k] = v
	}

	return r
}

func MapSlice[U, V any](f func(U) V, s []U) []V {
	r := make([]V, len(s))
	for i, v := range s {
		r[i] = f(v)
	}

	return r
}

func MapValues[K comparable, U, V any](f func(U) V, m map[K]U) map[K]V {
	r := make(map[K]V, len(m))
	for k, v := range m {
		r[k] = f(v)
	}

	return r
}

func Values[K comparable, V any](m map[K]V) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}

	return r
}
