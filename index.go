package indexmap

import (
	"github.com/yah01/container"
)

type PrimaryIndex[K comparable, V any] struct {
	extractField func(value *V) K

	inner map[K]*V
}

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

type SecondaryIndex[V any] struct {
	extractField func(value *V) []any

	inner map[any]container.Set[*V]
}

func NewSecondaryIndex[V any](extractField func(value *V) []any) *SecondaryIndex[V] {
	return &SecondaryIndex[V]{
		extractField: extractField,
		inner:        make(map[any]container.Set[*V]),
	}
}

func (index *SecondaryIndex[V]) get(key any) []*V {
	elems, ok := index.inner[key]
	if !ok {
		return nil
	}

	return elems.Collect()
}

func (index *SecondaryIndex[V]) insert(elem *V) {
	keys := index.extractField(elem)
	for i := range keys {
		elems, ok := index.inner[keys[i]]
		if !ok {
			elems = make(container.Set[*V])
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
