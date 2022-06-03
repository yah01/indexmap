# IndexMap
We often created a map with $ID \to Object$ to seek data, but this limits us to seek the data with only ID. to seek data with any field without SQL in database, IndexMap is the data structure you can reach this.

## Get Started
First, to create a IndexMap with primary index:
```golang
type Person struct {
	ID   int64
	Name string
	Age  int
}

persons := indexmap.NewIndexMap(indexmap.NewPrimaryIndex(func(value *Person) int64 {
    return value.ID
}))
```

Now it's just like the common map type, but then you can add more indexes to seek person with name:
```golang
persons.AddIndex("name", indexmap.NewSecondaryIndex(func(value *Person) []any {
    return []any{value.Name}
}))
```
You have to provide the way to extract keys for the inserted object, all keys must be comparable.

The insertion updates indexes automatically:
```golang
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
```

And seek data with primary index or the added index:
```golang
fmt.Printf("%+v\n", persons.Get(ashe.ID))
fmt.Printf("%+v\n", persons.GetBy("name", ashe.Name))
fmt.Printf("%+v\n", persons.Get(bob.ID))
fmt.Printf("%+v\n", persons.GetBy("name", bob.Name))
```
which outputs:
```
&{ID:1 Name:Ashe Age:39}
&{ID:1 Name:Ashe Age:39}
&{ID:2 Name:Bob Age:18}
&{ID:2 Name:Bob Age:18}
```

## One-To-Many/Many-To-Many Index
It's OK to create an index that's not one-to-one, The `GetBy()` method returns one of the object if many ones exist, `GetAllBy()` return a slice with all matched objects. For the example of many-to-many index, refer [contain_index_example](./examples/contain_index/main.go)

## Performance
Let $n$ be the number of elements inserted, $m$ be the number of indexes:
| Operation | Complexity |
| --------- | ---------- |
| Get       | $O(1)$     |
| GetBy     | $O(1)$     |
| Insert    | $O(m)$     |
| Remove    | $O(m)$     |
| AddIndex  | $O(n)$     |

The more indexes, the slower the write operations.