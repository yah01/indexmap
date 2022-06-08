package utils

import "github.com/yah01/container"

func UnionSet[T comparable](dst container.Set[T], src container.Set[T]) {
	for elem := range src {
		dst.Insert(elem)
	}
}

func IntersectSet[T comparable](dst container.Set[T], src container.Set[T]) {
	for elem := range dst {
		if !src.Contain(elem) {
			delete(dst, elem)
		}
	}
}

func ExceptSet[T comparable](dst container.Set[T], src container.Set[T]) {
	for elem := range src {
		delete(dst, elem)
	}
}
