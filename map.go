package indexmap

// IndexMap is a map supports seeking data with more indexes.
// Serializing a IndexMap as JSON results in the same as serializing a map,
// the result doesn't contain the index information, only data.
type IndexMap[K comparable, V any] struct {
	primaryIndex *PrimaryIndex[K, V]
	indexes      map[string]*SecondaryIndex[V]
}

// Create a IndexMap with a primary index,
// the primary index must be one-to-one.
func NewIndexMap[K comparable, V any](primaryIndex *PrimaryIndex[K, V]) *IndexMap[K, V] {
	return &IndexMap[K, V]{
		primaryIndex: primaryIndex,
		indexes:      make(map[string]*SecondaryIndex[V]),
	}
}

// Add a secondary index,
// build index for the data inserted,
// the return value indicates whether succeed to add index,
// false if the indexName existed.
func (imap *IndexMap[K, V]) AddIndex(indexName string, index *SecondaryIndex[V]) bool {
	if _, ok := imap.indexes[indexName]; ok {
		return false
	}

	imap.indexes[indexName] = index

	imap.primaryIndex.iterate(func(_ K, value *V) {
		index.insert(value)
	})

	return true
}

// Get value by the primary key,
// nil if key not exists.
func (imap *IndexMap[K, V]) Get(key K) *V {
	return imap.primaryIndex.get(key)
}

// Return one of the values for the given secondary key,
// No guarantee for which one is returned if more than one elements indexed by the key.
func (imap *IndexMap[K, V]) GetBy(indexName string, key any) *V {
	index, ok := imap.indexes[indexName]
	if !ok {
		return nil
	}

	elems := index.get(key)
	if len(elems) == 0 {
		return nil
	}

	return elems[0]
}

// Return all values the seeked by the key,
// nil if index or key not exists.
func (imap *IndexMap[K, V]) GetAllBy(indexName string, key any) []*V {
	index, ok := imap.indexes[indexName]
	if !ok {
		return nil
	}

	return index.get(key)
}

// Insert values into the map,
// also updates the indexes added,
// overwrite if a value with the same primary key existed
func (imap *IndexMap[K, V]) Insert(values ...*V) {
	for i := range values {
		imap.primaryIndex.insert(values[i])
		for _, index := range imap.indexes {
			index.insert(values[i])
		}
	}

}

// Remove values into the map,
// also updates the indexes added.
func (imap *IndexMap[K, V]) Remove(keys ...K) {
	for i := range keys {
		elem := imap.primaryIndex.get(keys[i])
		if elem == nil {
			continue
		}

		imap.primaryIndex.remove(keys[i])

		for _, index := range imap.indexes {
			index.remove(elem)
		}
	}
}

// Iterate all the elements,
// stop iteration if fn returns false,
// no any guarantee to the order.
func (imap *IndexMap[K, V]) Range(fn func(key K, value *V) bool) {
	for k, v := range imap.primaryIndex.inner {
		if !fn(k, v) {
			return
		}
	}
}

// Return all the keys and values.
func (imap *IndexMap[K, V]) Collect() ([]K, []*V) {
	var (
		keys   = make([]K, 0, imap.Len())
		values = make([]*V, 0, imap.Len())
	)
	for k, v := range imap.primaryIndex.inner {
		keys = append(keys, k)
		values = append(values, v)
	}

	return keys, values
}

// The number of elements.
func (imap *IndexMap[K, V]) Len() int {
	return len(imap.primaryIndex.inner)
}
