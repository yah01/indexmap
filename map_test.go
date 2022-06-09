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
		assert.Equal(t,
			&persons[i], imap.Get(persons[i].ID))

		assert.Equal(t,
			&persons[i], imap.GetBy(NameIndex, persons[i].Name))

		result := imap.GetAllBy(NameIndex, persons[i].Name)
		assert.Equal(t, 1, len(result))
		assert.Contains(t, result, &persons[i])
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
