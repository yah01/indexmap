package indexmap

import (
	"fmt"
)

// V must be a sturct pointer
type IndexMap[K comparable, V any] struct {
	primaryIndex *PrimaryIndex[K, V]
	indexes      map[string]*Index[V]
}

func NewArindex[K comparable, V any](primaryIndex *PrimaryIndex[K, V]) *IndexMap[K, V] {
	return &IndexMap[K, V]{
		primaryIndex: primaryIndex,
		indexes:      make(map[string]*Index[V]),
	}
}

func (armap *IndexMap[K, V]) AddIndex(indexName string, index *Index[V]) {
	armap.indexes[indexName] = index
}

func (armap *IndexMap[K, V]) Get(key K) *V {
	return armap.primaryIndex.get(key)
}

// Return one of the elements for the given secondary key,
// No guarantee for which one is returned if more than one elements indexed by the key
func (armap *IndexMap[K, V]) GetBy(indexName string, key any) *V {
	index, ok := armap.indexes[indexName]
	if !ok {
		panic(fmt.Sprintf("no index for field %s", indexName))
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
		panic(fmt.Sprintf("no index for field %s", indexName))
	}

	return index.get(key)
}

func (armap *IndexMap[K, V]) Insert(value *V) {
	armap.primaryIndex.insert(value)
	for _, index := range armap.indexes {
		index.insert(value)
	}
}

func (armap *IndexMap[K, V]) Remove(key K) {
	elem := armap.primaryIndex.get(key)
	if elem == nil {
		return
	}

	for _, index := range armap.indexes {
		index.remove(elem)
	}
}
