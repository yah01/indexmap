package indexmap


type PrimaryIndex[K comparable, V any] struct {
	extractField func(value *V) K

	inner map[K]*V
}

// Create an primary index,
// the extractField func must guarantee it makes the index one-to-one.
func NewPrimaryIndex[K comparable, V any](extractField func(value *V) K) *PrimaryIndex[K, V] {
	return &PrimaryIndex[K, V]{
		extractField: extractField,
		inner:        make(map[K]*V),
	}
}

func (index *PrimaryIndex[K, V]) get(key K) *V {
	return index.inner[key]
}

func (index *PrimaryIndex[K, V]) insert(elem *V) {
	index.inner[index.extractField(elem)] = elem
}

func (index *PrimaryIndex[K, V]) remove(key K) {
	delete(index.inner, key)
}

func (index *PrimaryIndex[K, V]) iterate(handler func(key K, value *V)) {
	for key, value := range index.inner {
		handler(key, value)
	}
}

type SecondaryIndex[V any] struct {
	extractField func(value *V) []any

	inner map[any]Set[*V]
}

// Create a secondary index,
// the extractField func returns the keys for seeking the value,
// It's OK that the same key seeks more than one values.
func NewSecondaryIndex[V any](extractField func(value *V) []any) *SecondaryIndex[V] {
	return &SecondaryIndex[V]{
		extractField: extractField,
		inner:        make(map[any]Set[*V]),
	}
}

func (index *SecondaryIndex[V]) get(key any) Set[*V] {
	set, ok := index.inner[key]
	if !ok {
		return nil
	}

	return set
}

func (index *SecondaryIndex[V]) insert(elem *V) {
	keys := index.extractField(elem)
	for i := range keys {
		elems, ok := index.inner[keys[i]]
		if !ok {
			elems = make(Set[*V])
			index.inner[keys[i]] = elems
		}

		elems.Insert(elem)
	}
}

func (index *SecondaryIndex[V]) remove(elem *V) {
	keys := index.extractField(elem)
	for i := range keys {
		elems, ok := index.inner[keys[i]]
		if ok {
			elems.Remove(elem)
		}
	}
}
