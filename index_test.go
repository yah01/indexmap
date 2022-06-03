package indexmap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Person struct {
	ID         int64
	Name       string
	FormerName string
}

func TestPrimaryIndex(t *testing.T) {
	index := NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	})

	persons := []Person{
		{1, "Ashe", "ashe"},
		{2, "Bob", "bob"},
		{3, "Cassidy", "McCree"},
	}

	for i := range persons {
		index.insert(&persons[i])
	}

	for i := range persons {
		assert.Equal(t,
			&persons[i], index.get(persons[i].ID))
	}

	// Insert overwrite
	overwritePerson := &Person{1, "Tracer", "tracer"}
	index.insert(overwritePerson)
	assert.Equal(t, overwritePerson, index.get(overwritePerson.ID))

	// Remove get nil
	index.remove(overwritePerson.ID)
	assert.Nil(t, index.get(overwritePerson.ID))
}

func TestSecondaryIndex(t *testing.T) {
	// Many-to-One index
	index := NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name, value.FormerName}
	})

	persons := []Person{
		{1, "Ashe", "ashe"},
		{2, "Bob", "bob"},
		{3, "Cassidy", "McCree"},
	}

	for i := range persons {
		index.insert(&persons[i])
	}

	for i := range persons {
		result := index.get(persons[i].Name)

		assert.Equal(t, 1, len(result))
		assert.Contains(t,
			result, &persons[i])

		result = index.get(persons[i].FormerName)

		assert.Equal(t, 1, len(result))
		assert.Contains(t,
			result, &persons[i])
	}

	// Insert makes One-to-Many, Many-to-Many
	overwritePerson := &Person{4, "Ashe", "alice"}
	index.insert(overwritePerson)

	result := index.get(overwritePerson.Name)
	assert.Equal(t, 2, len(result))
	assert.Contains(t, result, overwritePerson)
	assert.Contains(t, result, &persons[0])
	result = index.get(overwritePerson.FormerName)

	assert.Equal(t, 1, len(result))
	assert.Contains(t,
		result, overwritePerson)

	// Remove
	index.remove(overwritePerson)
	result = index.get(persons[0].Name)
	assert.Equal(t, 1, len(result))
	assert.Contains(t, result, &persons[0])

	result = index.get(persons[0].FormerName)
	assert.Equal(t, 1, len(result))
	assert.Contains(t, result, &persons[0])
}
