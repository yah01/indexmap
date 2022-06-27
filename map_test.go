package indexmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexMap(t *testing.T) {
	imap := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))

	ok := imap.AddIndex(NameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name}
	}))
	assert.True(t, ok)

	persons := GenPersons()
	InsertData(imap, persons)

	for i, person := range persons {
		assert.Equal(t,
			persons[i], imap.Get(persons[i].ID))

		assert.Equal(t,
			person, imap.Get(person.ID))

		assert.Equal(t,
			person, imap.GetBy(NameIndex, person.Name))

		assert.Nil(t, imap.GetBy(InvalidIndex, person.Name))

		result := imap.GetAllBy(NameIndex, person.Name)
		assert.Equal(t, 1, len(result))
		assert.Contains(t, result, person)

		assert.Nil(t, imap.getAllBy(InvalidIndex, person.Name))
	}

	// Add index after inserting data
	ok = imap.AddIndex(CityIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.City}
	}))
	assert.True(t, ok)

	for _, person := range persons {
		assert.Equal(t,
			person, imap.GetBy(NameIndex, person.Name))

		result := imap.GetAllBy(CityIndex, person.City)
		assert.Contains(t, result, person)
	}

	// Remove
	imap.Remove(persons[0].ID)
	assert.Nil(t, imap.Get(persons[0].ID))
	assert.Nil(t, imap.GetBy(NameIndex, persons[0].Name))
	assert.Empty(t, imap.GetAllBy(NameIndex, persons[0].Name))

	imap.RemoveBy(CityIndex, "San Francisco")
	assert.Empty(t, imap.GetAllBy(CityIndex, "San Francisco"))
	assert.Equal(t, 1, len(imap.GetAllBy(CityIndex, "Shanghai")))

	// Update
	imap.Clear()
	InsertData(imap, persons)
	imap.Update(persons[0].ID, func(value *Person) (*Person, bool) {
		value.Name = "Tracer"
		return value, true
	})
	assert.Equal(t, "Tracer", imap.Get(persons[0].ID).Name)

	count := len(imap.GetAllBy(CityIndex, "Shanghai"))
	imap.UpdateBy(CityIndex, "Shanghai", func(value *Person) (*Person, bool) {
		value.City = "Beijing"
		return value, true
	})
	assert.Empty(t, imap.GetAllBy(CityIndex, "Shanghai"))
	assert.Equal(t, count, len(imap.GetAllBy(CityIndex, "Beijing")))

	// Collect
	keys, values := imap.Collect()
	assert.Equal(t, imap.Len(), len(keys))
	assert.Equal(t, imap.Len(), len(values))
	for i := range keys {
		assert.Equal(t, values[i], imap.Get(keys[i]))
	}

	// Range
	count = 0
	imap.Range(func(key int64, value *Person) bool {
		count++
		assert.Equal(t, value, imap.Get(key))
		return true
	})
	assert.Equal(t, imap.Len(), count)
}

func TestAddExistedIndex(t *testing.T) {
	imap := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))

	ok := imap.AddIndex(NameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name}
	}))
	assert.True(t, ok)

	ok = imap.AddIndex(NameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.City}
	}))
	assert.False(t, ok)
}
