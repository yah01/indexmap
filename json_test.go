package indexmap

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonMarshal(t *testing.T) {
	imap := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))

	imap.AddIndex(NameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name}
	}))

	persons := map[int64]*Person{
		1: {1, "Ashe", 38, "San Francisco", []string{"Bob", "Cassidy"}},
		2: {2, "Bob", 18, "San Francisco", nil},
		3: {3, "Cassidy", 40, "Shanghai", []string{"Bob", "Ashe"}},
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
	imap := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))

	imap.AddIndex(NameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name}
	}))

	persons := map[int64]*Person{
		1: {1, "Ashe", 38, "San Francisco", []string{"Bob", "Cassidy"}},
		2: {2, "Bob", 18, "San Francisco", nil},
		3: {3, "Cassidy", 40, "Shanghai", []string{"Bob", "Ashe"}},
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
	imap.AddIndex(CityIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.City}
	}))

	for id := range persons {
		assert.Contains(t,
			imap.GetAllBy(CityIndex, persons[id].City), persons[id])
	}

	// Mock JSON with syntax error
	mapData = append(mapData, '}')
	err = json.Unmarshal(mapData, imap)
	assert.Error(t, err)
}
