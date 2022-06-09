package indexmap

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

func GenPersons() map[int64]*Person {
	return map[int64]*Person{
		0: {0, "Ashe", 38, "San Francisco", []string{"Bob", "Cassidy"}},
		1: {1, "Bob", 18, "San Francisco", nil},
		2: {2, "Cassidy", 40, "Shanghai", []string{"Bob", "Ashe"}},
	}
}

func InsertData[K comparable, V any](imap *IndexMap[K, V], data map[K]*V) {
	for _, v := range data {
		imap.Insert(v)
	}
}
