package main

import (
	"fmt"

	"github.com/yah01/indexmap"
)

type Person struct {
	ID   int64
	Name string
	Age  int
	City string
	Like []string
}

func main() {
	persons := indexmap.NewIndexMap(indexmap.NewPrimaryIndex(func(value *Person) int64 {
		return value.ID
	}))

	persons.AddIndex("name", indexmap.NewSecondaryIndex(func(value *Person) []any {
		return []any{value.Name}
	}))

	ashe := &Person{
		ID:   1,
		Name: "Ashe",
		Age:  39,
		City: "San Francisco",
		Like: []string{"Bob", "Cassidy"},
	}
	bob := &Person{
		ID:   2,
		Name: "Bob",
		Age:  18,
		City: "San Francisco",
	}
	cassidy := &Person{
		ID:   3,
		Name: "Cassidy",
		Age:  40,
		City: "Shanghai",
		Like: []string{"Ashe", "Bob"},
	}

	persons.Insert(ashe)
	persons.Insert(bob)
	persons.Insert(cassidy)

	persons.AddIndex("city", indexmap.NewSecondaryIndex(func(value *Person) []any {
		return []any{value.City}
	}))

	// Like is a "contain" index
	persons.AddIndex("like", indexmap.NewSecondaryIndex(func(value *Person) []any {
		like := make([]any, 0, len(value.Like))
		for i := range value.Like {
			like = append(like, value.Like[i])
		}
		return like
	}))

	fmt.Println("Search with ID or Name:")
	fmt.Printf("%+v\n", persons.Get(ashe.ID))
	fmt.Printf("%+v\n", persons.GetBy("name", ashe.Name))

	fmt.Println("\nSearch persons come from San Francisco:")
	for _, person := range persons.GetAllBy("city", "San Francisco") {
		fmt.Printf("%+v\n", person)
	}

	fmt.Println("\nSearch persons like Bob")
	for _, person := range persons.GetAllBy("like", "Bob") {
		fmt.Printf("%+v\n", person)
	}

	// DO NOT modify a internal value outside
	// person := persons.GetBy("name", "Ashe")
	// person.City = "Shanghai"
	// persons.Insert(person)

	// Modify the internal value with Update()/UpdateBy()
	persons.UpdateBy("name", "Ashe", func(value *Person) (*Person, bool) {
		if value.City == "Shanghai" {
			return value, false
		}
		value.City = "Shanghai"
		return value, true
	})
}

/*
Outputs:
Search with ID or Name:
&{ID:1 Name:Ashe Age:39 City:San Francisco Like:[Bob Cassidy]}
&{ID:1 Name:Ashe Age:39 City:San Francisco Like:[Bob Cassidy]}

Search persons come from San Francisco:
&{ID:1 Name:Ashe Age:39 City:San Francisco Like:[Bob Cassidy]}
&{ID:2 Name:Bob Age:18 City:San Francisco Like:[]}

Search persons like Bob
&{ID:3 Name:Cassidy Age:40 City:Shanghai Like:[Ashe Bob]}
&{ID:1 Name:Ashe Age:39 City:San Francisco Like:[Bob Cassidy]}
*/
