package indexmap

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonMarshal(t *testing.T) {
	const (
		NameIndex       = "name_index"
		FormerNameIndex = "former_name_index"
	)

	imap := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))

	imap.AddIndex(NameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name}
	}))

	persons := map[int64]*Person{
		1: {1, "Ashe", "ashe"},
		2: {2, "Bob", "bob"},
		3: {3, "Cassidy", "McCree"},
	}

	for _, v := range persons {
		imap.Insert(v)
	}

	imapData, err := json.Marshal(imap)
	assert.NoError(t, err)

	decodedMap := make(map[int64]*Person)
	err = json.Unmarshal(imapData, &decodedMap)
	assert.NoError(t, err)

	assert.EqualValues(t, persons, decodedMap)
}

func TestJsonUnmarshal(t *testing.T) {
	const (
		NameIndex       = "name_index"
		FormerNameIndex = "former_name_index"
	)

	imap := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))

	imap.AddIndex(NameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name}
	}))

	persons := map[int64]*Person{
		1: {1, "Ashe", "ashe"},
		2: {2, "Bob", "bob"},
		3: {3, "Cassidy", "McCree"},
	}

	for _, v := range persons {
		imap.Insert(v)
	}

	mapData, err := json.Marshal(persons)
	assert.NoError(t, err)

	err = json.Unmarshal(mapData, imap)
	assert.NoError(t, err)

	for i := range persons {
		assert.Equal(t,
			persons[i], imap.Get(persons[i].ID))

		assert.Equal(t,
			persons[i], imap.GetBy(NameIndex, persons[i].Name))

		result := imap.GetAllBy(NameIndex, persons[i].Name)
		assert.Equal(t, 1, len(result))
		assert.Contains(t, result, persons[i])
	}

	// Add index after inserting data
	imap.AddIndex(FormerNameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.FormerName}
	}))

	for i := range persons {
		assert.Equal(t,
			persons[i], imap.GetBy(FormerNameIndex, persons[i].FormerName))

		result := imap.GetAllBy(FormerNameIndex, persons[i].FormerName)
		assert.Equal(t, 1, len(result))
		assert.Contains(t, result, persons[i])
	}
}
