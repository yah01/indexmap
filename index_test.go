package indexmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	ID   int64
	Name string
	Age  int
	City string
	Like []string
}

const (
	InvalidIndex = "invalid"
	NameIndex    = "name"
	CityIndex    = "city"
	LikeIndex    = "like"
)

func TestPrimaryIndex(t *testing.T) {
	index := NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	})

	persons := []Person{
		{1, "Ashe", 38, "San Francisco", []string{"Bob", "Cassidy"}},
		{2, "Bob", 18, "San Francisco", nil},
		{3, "Cassidy", 40, "Shanghai", []string{"Bob", "Ashe"}},
	}

	for i := range persons {
		index.insert(&persons[i])
	}

	for i := range persons {
		assert.Equal(t,
			&persons[i], index.get(persons[i].ID))
	}

	// Insert overwrite
	overwritePerson := &Person{1, "Tracer", 23, "London", []string{"Bob"}}
	index.insert(overwritePerson)
	assert.Equal(t, overwritePerson, index.get(overwritePerson.ID))

	// Remove get nil
	index.remove(overwritePerson.ID)
	assert.Nil(t, index.get(overwritePerson.ID))
}

func TestSecondaryIndex(t *testing.T) {
	// Many-to-One index
	index := NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name, value.City}
	})

	persons := []Person{
		{1, "Ashe", 38, "San Francisco", []string{"Bob", "Cassidy"}},
		{2, "Bob", 18, "San Francisco", nil},
		{3, "Cassidy", 40, "Shanghai", []string{"Bob", "Ashe"}},
	}

	for i := range persons {
		index.insert(&persons[i])
	}

	for i := range persons {
		result := index.get(persons[i].Name)

		assert.Equal(t, 1, len(result))
		assert.Contains(t,
			result, &persons[i])

		result = index.get(persons[i].City)
		assert.Contains(t,
			result, &persons[i])
	}

	// Insert makes One-to-Many, Many-to-Many
	ashe2 := &Person{4, "Ashe", 83, "Chengdu", nil}
	index.insert(ashe2)

	result := index.get(ashe2.Name)
	assert.Equal(t, 2, len(result))
	assert.Contains(t, result, ashe2)
	assert.Contains(t, result, &persons[0])

	// Remove
	index.remove(ashe2)
	result = index.get(persons[0].Name)
	assert.Equal(t, 1, len(result))
	assert.Contains(t, result, &persons[0])
}
