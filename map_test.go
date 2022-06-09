package indexmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexMap(t *testing.T) {
	imap := NewIndexMap(NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))

	imap.AddIndex(NameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name}
	}))

	persons := []Person{
		{1, "Ashe", 38, "San Francisco", []string{"Bob", "Cassidy"}},
		{2, "Bob", 18, "San Francisco", nil},
		{3, "Cassidy", 40, "Shanghai", []string{"Bob", "Ashe"}},
	}

	for i := range persons {
		imap.Insert(&persons[i])
	}

	for i := range persons {
		assert.True(t, imap.Contain(persons[i].ID))

		assert.Equal(t,
			&persons[i], imap.Get(persons[i].ID))

		assert.Equal(t,
			&persons[i], imap.GetBy(NameIndex, persons[i].Name))

		assert.Nil(t, imap.GetBy(InvalidIndex, persons[i].Name))

		result := imap.GetAllBy(NameIndex, persons[i].Name)
		assert.Equal(t, 1, len(result))
		assert.Contains(t, result, &persons[i])

		assert.Nil(t, imap.getAllBy(InvalidIndex, persons[i].Name))
	}

	// Add index after inserting data
	imap.AddIndex(CityIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.City}
	}))

	for i := range persons {
		assert.Equal(t,
			&persons[i], imap.GetBy(NameIndex, persons[i].Name))

		result := imap.GetAllBy(CityIndex, persons[i].City)
		assert.Contains(t, result, &persons[i])
	}

	// Remove
	imap.Remove(persons[0].ID)
	assert.Nil(t, imap.Get(persons[0].ID))
	assert.Nil(t, imap.GetBy(NameIndex, persons[0].Name))
	assert.Empty(t, imap.GetAllBy(NameIndex, persons[0].Name))

	// Update
	imap.Insert(&persons[0])
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

	imap.AddIndex(NameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name}
	}))

	ok := imap.AddIndex(NameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.City}
	}))

	assert.False(t, ok)
}
