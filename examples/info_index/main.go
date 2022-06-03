package main

import (
	"fmt"

	"github.com/yah01/indexmap"
)

type Person struct {
	ID   int64
	Name string
	Age  int
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
	}
	bob := &Person{
		ID:   2,
		Name: "Bob",
		Age:  18,
	}

	persons.Insert(ashe)
	persons.Insert(bob)

	fmt.Printf("%+v\n", persons.Get(ashe.ID))
	fmt.Printf("%+v\n", persons.GetBy("name", ashe.Name))
	fmt.Printf("%+v\n", persons.Get(bob.ID))
	fmt.Printf("%+v\n", persons.GetBy("name", bob.Name))
}

/*
Outputs:
&{ID:1 Name:Ashe Age:39}
&{ID:1 Name:Ashe Age:39}
&{ID:2 Name:Bob Age:18}
&{ID:2 Name:Bob Age:18}
*/
