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

var names = []string{"James", "Mary", "Robert", "Patricia", "John", "Jennifer", "Michael", "Linda", "David", "Elizabeth", "William", "Barbara", "Richard", "Susan", "Joseph", "Jessica", "Thomas", "Sarah", "Charles", "Karen", "Christopher", "Lisa", "Daniel", "Nancy", "Matthew", "Betty", "Anthony", "Margaret", "Mark", "Sandra", "Donald", "Ashley", "Steven", "Kimberly", "Paul", "Emily", "Andrew", "Donna", "Joshua", "Michelle", "Kenneth", "Carol", "Kevin", "Amanda", "Brian", "Dorothy", "George", "Melissa", "Timothy", "Deborah", "Ronald", "Stephanie", "Edward", "Rebecca", "Jason", "Sharon", "Jeffrey", "Laura", "Ryan", "Cynthia", "Jacob", "Kathleen", "Gary", "Amy", "Nicholas", "Angela", "Eric", "Shirley", "Jonathan", "Anna", "Stephen", "Brenda", "Larry", "Pamela", "Justin", "Emma", "Scott", "Nicole", "Brandon", "Helen", "Benjamin", "Samantha", "Samuel", "Katherine", "Gregory", "Christine", "Alexander", "Debra", "Frank", "Rachel", "Patrick", "Carolyn", "Raymond", "Janet", "Jack", "Catherine", "Dennis", "Maria", "Jerry", "Heather", "Tyler", "Diane", "Aaron", "Ruth", "Jose", "Julie", "Adam", "Olivia", "Nathan", "Joyce", "Henry", "Virginia", "Douglas", "Victoria", "Zachary", "Kelly", "Peter", "Lauren", "Kyle", "Christina", "Ethan", "Joan", "Walter", "Evelyn", "Noah", "Judith", "Jeremy", "Megan", "Christian", "Andrea", "Keith", "Cheryl", "Roger", "Hannah", "Terry", "Jacqueline", "Gerald", "Martha", "Harold", "Gloria", "Sean", "Teresa", "Austin", "Ann", "Carl", "Sara", "Arthur", "Madison", "Lawrence", "Frances", "Dylan", "Kathryn", "Jesse", "Janice", "Jordan", "Jean", "Bryan", "Abigail", "Billy", "Alice", "Joe", "Julia", "Bruce", "Judy", "Gabriel", "Sophia", "Logan", "Grace", "Albert", "Denise", "Willie", "Amber", "Alan", "Doris", "Juan", "Marilyn", "Wayne", "Danielle", "Elijah", "Beverly", "Randy", "Isabella", "Roy", "Theresa", "Vincent", "Diana", "Ralph", "Natalie", "Eugene", "Brittany", "Russell", "Charlotte", "Bobby", "Marie", "Mason", "Kayla", "Philip", "Alexis", "Louis", "Lori"}

func InsertData[K comparable, V any](imap *IndexMap[K, V], data map[K]*V) {
	for _, v := range data {
		imap.Insert(v)
	}
}
