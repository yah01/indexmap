package indexmap

type Set[T comparable] map[T]struct{}

func (set Set[T]) Contain(elems ...T) bool {
	for i := range elems {
		_, ok := set[elems[i]]
		if !ok {
			return false
		}
	}

	return true
}

func (set Set[T]) Insert(elems ...T) {
	for i := range elems {
		set[elems[i]] = struct{}{}
	}
}

func (set Set[T]) Remove(elems ...T) {
	for i := range elems {
		delete(set, elems[i])
	}
}

func (set Set[T]) Collect() []T {
	results := make([]T, 0, set.Len())
	for elem := range set {
		results = append(results, elem)
	}

	return results
}

func (set Set[T]) Len() int {
	return len(set)
}
