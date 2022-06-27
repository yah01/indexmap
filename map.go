package indexmap

import "sync"

// IndexMap is a map supports seeking data with more indexes.
// Serializing a IndexMap as JSON results in the same as serializing a map,
// the result doesn't contain the index information, only data.
// NOTE: DO NOT insert nil value into the IndexMap
type IndexMap[K comparable, V any] struct {
	primaryIndex *PrimaryIndex[K, V]
	indexes      map[string]*SecondaryIndex[V]
	lock         sync.RWMutex
	version      int64
}

// Create a IndexMap with a primary index,
// the primary index must be a one-to-one mapping.
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
	imap.lock.Lock()
	defer imap.lock.Unlock()

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
	imap.lock.RLock()
	defer imap.lock.RUnlock()

	return imap.primaryIndex.get(key)
}

// Return one of the values for the given secondary key,
// No guarantee for which one is returned if more than one elements indexed by the key.
func (imap *IndexMap[K, V]) GetBy(indexName string, key any) *V {
	imap.lock.RLock()
	defer imap.lock.RUnlock()

	index, ok := imap.indexes[indexName]
	if !ok {
		return nil
	}

	elems := index.get(key)
	if len(elems) == 0 {
		return nil
	}

	for value := range elems {
		return value
	}

	return nil
}

// Return all values the seeked by the key,
// nil if index or key not exists.
func (imap *IndexMap[K, V]) GetAllBy(indexName string, key any) []*V {
	imap.lock.RLock()
	defer imap.lock.RUnlock()

	values := imap.getAllBy(indexName, key)
	if values == nil {
		return nil
	}

	return values.Collect()
}

// Return true if the value with given key exists,
// false otherwise.
func (imap *IndexMap[K, V]) Contain(key K) bool {
	return imap.Get(key) != nil
}

// Insert values into the map,
// also updates the indexes added,
// overwrite if a value with the same primary key existed.
// NOTE: insert an modified existed value with the same address may confuse the index, use Update() to do this.
func (imap *IndexMap[K, V]) Insert(values ...*V) {
	imap.lock.Lock()
	defer imap.lock.Unlock()

	imap.insert(values...)
}

// insert is the lock free version of Inser
func (imap *IndexMap[K, V]) insert(values ...*V) {

	for i := range values {
		oldKey := imap.primaryIndex.extractField(values[i])
		// don't use Get(oldKey) that rlock on locked map (dead lock)
		old := imap.primaryIndex.get(oldKey)
		imap.primaryIndex.insert(values[i])
		for _, index := range imap.indexes {
			if old != nil {
				index.remove(old)
			}

			index.insert(values[i])
		}
	}
}

// An UpdateFn modifies the given value,
// and returns the modified value, they could be the same object,
// true if the object is modified,
// false otherwise
type UpdateFn[V any] func(value *V) (*V, bool)

// Update the value for the given key,
// it removes the old one if exists, and inserts updateFn(old) if modified and not nil.
func (imap *IndexMap[K, V]) Update(key K, updateFn UpdateFn[V]) {
	imap.lock.Lock()
	defer imap.lock.Unlock()

	// don't use Get(key) that rlock on locked map (dead lock)
	old := imap.primaryIndex.get(key)
	if old != nil {
		imap.remove(key)
	}

	new, modified := updateFn(old)
	if modified && new != nil {
		imap.insert(new)
	}
}

// Update the values for the given index and key,
// it removes the old ones if exist, and inserts updateFn(old) for every old ones if not nil.
// NOTE: the modified values have to be with unique primary key
func (imap *IndexMap[K, V]) UpdateBy(indexName string, key any, updateFn UpdateFn[V]) {
	imap.lock.Lock()
	defer imap.lock.Unlock()

	oldValueSet := imap.getAllBy(indexName, key)
	if len(oldValueSet) == 0 {
		return
	}

	oldValues := oldValueSet.Collect()

	imap.removeValues(oldValues...)

	for _, old := range oldValues {
		new, modified := updateFn(old)
		if modified && new != nil {
			imap.insert(new)
		}
	}
}

// Remove values into the map,
// also updates the indexes added.
func (imap *IndexMap[K, V]) Remove(keys ...K) {
	imap.lock.Lock()
	defer imap.lock.Unlock()

	imap.remove(keys...)
}

// remove is the lock free  version Remove
func (imap *IndexMap[K, V]) remove(keys ...K) {

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

// Remove values into the map,
// also updates the indexes added.
func (imap *IndexMap[K, V]) RemoveBy(indexName string, keys ...any) {
	imap.lock.Lock()
	defer imap.lock.Unlock()

	imap.removeBy(indexName, keys...)
}

// removBy is the lock free verison of RemoveBy
func (imap *IndexMap[K, V]) removeBy(indexName string, keys ...any) {
	for i := range keys {
		values := imap.getAllBy(indexName, keys[i])
		if values == nil {
			continue
		}

		imap.removeValueSet(values)
	}
}

// Remove all values.
func (imap *IndexMap[K, V]) Clear() {
	imap.lock.Lock()
	defer imap.lock.Unlock()

	for k := range imap.primaryIndex.inner {
		delete(imap.primaryIndex.inner, k)
	}

	for i := range imap.indexes {
		for k := range imap.indexes[i].inner {
			delete(imap.indexes[i].inner, k)
		}
	}
}

// Iterate all the elements,
// stop iteration if fn returns false,
// no any guarantee to the order.
// don't use modifying calls to this indexmap while the Range is nunning
// that may cause dead locks.
func (imap *IndexMap[K, V]) Range(fn func(key K, value *V) bool) {
	imap.lock.RLock()
	defer imap.lock.RUnlock()

	for k, v := range imap.primaryIndex.inner {
		if !fn(k, v) {
			return
		}
	}
}

// Return all the keys and values.
func (imap *IndexMap[K, V]) Collect() ([]K, []*V) {
	imap.lock.RLock()
	defer imap.lock.RUnlock()

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
	imap.lock.RLock()
	defer imap.lock.RUnlock()

	return len(imap.primaryIndex.inner)
}

// getAllBy ist the lock free version if GetAllBy(...)
func (imap *IndexMap[K, V]) getAllBy(indexName string, key any) Set[*V] {
	index, ok := imap.indexes[indexName]
	if !ok {
		return nil
	}
	return index.get(key)
}

// All values must exists
func (imap *IndexMap[K, V]) removeValues(values ...*V) {
	for i := range values {
		imap.remove(imap.primaryIndex.extractField(values[i]))
	}
}

// All values must exists
func (imap *IndexMap[K, V]) removeValueSet(values Set[*V]) {
	for value := range values {
		imap.remove(imap.primaryIndex.extractField(value))
	}
}
