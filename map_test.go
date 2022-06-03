package indexmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexMap(t *testing.T) {
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

	persons := []Person{
		{1, "Ashe", "ashe"},
		{2, "Bob", "bob"},
		{3, "Cassidy", "McCree"},
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
	imap.AddIndex(FormerNameIndex, NewSecondaryIndex(func(value *Person) []any {
		return []any{value.FormerName}
	}))

	for i := range persons {
		assert.Equal(t,
			&persons[i], imap.GetBy(FormerNameIndex, persons[i].FormerName))

		result := imap.GetAllBy(FormerNameIndex, persons[i].FormerName)
		assert.Equal(t, 1, len(result))
		assert.Contains(t, result, &persons[i])
	}

	// Remove
	imap.Remove(persons[0].ID)
	assert.Nil(t, imap.Get(persons[0].ID))
	assert.Nil(t, imap.GetBy(NameIndex, persons[0].Name))
	assert.Nil(t, imap.GetBy(FormerNameIndex, persons[0].Name))
	assert.Empty(t, imap.GetAllBy(NameIndex, persons[0].Name))
	assert.Empty(t, imap.GetAllBy(FormerNameIndex, persons[0].Name))
}
