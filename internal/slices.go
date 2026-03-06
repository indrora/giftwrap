package internal

func SliceDice[E any, R any](S []E, F func(E) R) []R {
	ret := make([]R, len(S))
	for i := range S {
		ret[i] = F(S[i])
	}
	return ret
}

func FilterDown[E any, R any](S []E, F func(E) *R) []R {
	ret := make([]R, 0, len(S))
	for i := range S {
		v := F(S[i])
		if v != nil {
			ret = append(ret, *v)
		}
	}

	return ret
}
