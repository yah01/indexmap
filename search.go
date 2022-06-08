package indexmap

import (
	"github.com/yah01/container"
	"github.com/yah01/indexmap/utils"
)

const (
	lookupTypeNone lookupType = iota
	lookupTypeAnd
	lookupTypeOr
	lookupTypeExclude
)

type lookupType int

type indexLookup struct {
	lookupType lookupType
	indexName  string
	keys       []any
}

type SearchStream[K comparable, V any] struct {
	chain []*indexLookup

	imap *IndexMap[K, V]
}

func (stream *SearchStream[K, V]) And(indexName string, keys ...any) *SearchStream[K, V] {
	stream.chain = append(stream.chain, &indexLookup{
		lookupType: lookupTypeAnd,
		indexName:  indexName,
		keys:       keys,
	})

	return stream
}

func (stream *SearchStream[K, V]) Or(indexName string, keys ...any) *SearchStream[K, V] {
	stream.chain = append(stream.chain, &indexLookup{
		lookupType: lookupTypeOr,
		indexName:  indexName,
		keys:       keys,
	})

	return stream
}

func (stream *SearchStream[K, V]) Exclude(indexName string, keys ...any) *SearchStream[K, V] {
	stream.chain = append(stream.chain, &indexLookup{
		lookupType: lookupTypeExclude,
		indexName:  indexName,
		keys:       keys,
	})

	return stream
}

func (stream *SearchStream[K, V]) Excute() []*V {
	values := make(container.Set[*V])

	for _, lookup := range stream.chain {
		result := stream.executeLookup(lookup)

		switch lookup.lookupType {
		case lookupTypeAnd:
			utils.IntersectSet(values, result)

		case lookupTypeOr:
			utils.UnionSet(values, result)

		case lookupTypeExclude:
			utils.ExceptSet(values, result)
		}
	}

	return values.Collect()
}

func (stream *SearchStream[K, V]) executeLookup(lookup *indexLookup) container.Set[*V] {
	result := make(container.Set[*V])
	for i := range lookup.keys {
		utils.UnionSet(result,
			stream.imap.getAllBy(lookup.indexName, lookup.keys[i]))
	}

	return result
}

func (imap *IndexMap[K, V]) Search() *SearchStream[K, V] {
	return &SearchStream[K, V]{
		chain: make([]*indexLookup, 0),
		imap:  imap,
	}
}
