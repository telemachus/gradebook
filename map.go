package gradebook

// See https://stackoverflow.com/a/71905349 for the use of M ~map[K]V.

// keys returns the keys from a map[string]any.
func keys[M ~map[K]V, K comparable, V any](m M) []K {
	ks := make([]K, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}

	return ks
}

// vals returns the values from a map[comparable]string.
func vals[M ~map[K]V, K comparable, V any](m M) []V {
	vs := make([]V, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}

	return vs
}
