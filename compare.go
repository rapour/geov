package geov

func sameMap[T comparable](m1, m2 map[T]bool) bool {

	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok || v1 != v2 {
			return false
		}
	}
	return true
}

func subsetMap[T comparable](set, subset map[T]bool) bool {

	for k, v1 := range subset {
		v2, ok := set[k]
		if !ok || v1 != v2 {
			return false
		}
	}
	return true

}
