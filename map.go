package indexmap

// V must be a sturct pointer
type IndexMap[K comparable, V any] struct {
	primaryIndex *PrimaryIndex[K, V]
	indexes      map[string]*SecondaryIndex[V]
}

func NewIndexMap[K comparable, V any](primaryIndex *PrimaryIndex[K, V]) *IndexMap[K, V] {
	return &IndexMap[K, V]{
		primaryIndex: primaryIndex,
		indexes:      make(map[string]*SecondaryIndex[V]),
	}
}

func (armap *IndexMap[K, V]) AddIndex(indexName string, index *SecondaryIndex[V]) {
	armap.indexes[indexName] = index

	armap.primaryIndex.iterate(func(_ K, value *V) {
		index.insert(value)
	})
}

func (armap *IndexMap[K, V]) Get(key K) *V {
	return armap.primaryIndex.get(key)
}

// Return one of the elements for the given secondary key,
// No guarantee for which one is returned if more than one elements indexed by the key
func (armap *IndexMap[K, V]) GetBy(indexName string, key any) *V {
	index, ok := armap.indexes[indexName]
	if !ok {
		return nil
	}

	elems := index.get(key)
	if len(elems) == 0 {
		return nil
	}

	return elems[0]
}

func (armap *IndexMap[K, V]) GetAllBy(indexName string, key any) []*V {
	index, ok := armap.indexes[indexName]
	if !ok {
		return nil
	}

	return index.get(key)
}

func (armap *IndexMap[K, V]) Insert(values ...*V) {
	for i := range values {
		armap.primaryIndex.insert(values[i])
		for _, index := range armap.indexes {
			index.insert(values[i])
		}
	}

}

func (armap *IndexMap[K, V]) Remove(keys ...K) {
	for i := range keys {
		elem := armap.primaryIndex.get(keys[i])
		if elem == nil {
			continue
		}

		armap.primaryIndex.remove(keys[i])

		for _, index := range armap.indexes {
			index.remove(elem)
		}
	}
}
