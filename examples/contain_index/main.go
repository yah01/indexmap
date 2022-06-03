package main

import (
	"fmt"

	"github.com/yah01/indexmap"
)

type Replica struct {
	ID    int64
	Nodes []int64
}

const (
	NodeContainIndex = "node_contain"
)

func main() {
	replicas := indexmap.NewIndexMap(indexmap.NewPrimaryIndex(func(value *Replica) int64 {
		return value.ID
	}))
	replicas.AddIndex(NodeContainIndex, indexmap.NewSecondaryIndex(func(value *Replica) []any {
		keys := make([]any, 0, len(value.Nodes))
		for _, node := range value.Nodes {
			keys = append(keys, node)
		}
		return keys
	}))

	ashe := &Replica{
		ID:    1,
		Nodes: []int64{1, 2, 3},
	}
	bob := &Replica{
		ID:    2,
		Nodes: []int64{2, 3, 4},
	}
	cindy := &Replica{
		ID:    3,
		Nodes: []int64{3, 4, 5},
	}

	replicas.Insert(ashe)
	replicas.Insert(bob)
	replicas.Insert(cindy)

	fmt.Printf("%+v\n", replicas.Get(ashe.ID))
	fmt.Printf("%+v\n", replicas.Get(bob.ID))

	for nodeID := int64(1); nodeID <= 5; nodeID++ {
		for _, replica := range replicas.GetAllBy(NodeContainIndex, nodeID) {
			fmt.Printf("%+v", replica)
		}
		fmt.Println()
	}

	fmt.Println("remove bob")
	replicas.Remove(bob.ID)
	for nodeID := int64(1); nodeID <= 5; nodeID++ {
		for _, replica := range replicas.GetAllBy(NodeContainIndex, nodeID) {
			fmt.Printf("%+v", replica)
		}
		fmt.Println()
	}
}

/*
Outputs:
&{ID:1 Nodes:[1 2 3]}
&{ID:2 Nodes:[2 3 4]}
&{ID:1 Nodes:[1 2 3]}
&{ID:1 Nodes:[1 2 3]}&{ID:2 Nodes:[2 3 4]}
&{ID:1 Nodes:[1 2 3]}&{ID:2 Nodes:[2 3 4]}&{ID:3 Nodes:[3 4 5]}
&{ID:2 Nodes:[2 3 4]}&{ID:3 Nodes:[3 4 5]}
&{ID:3 Nodes:[3 4 5]}
remove bob
&{ID:1 Nodes:[1 2 3]}
&{ID:1 Nodes:[1 2 3]}
&{ID:3 Nodes:[3 4 5]}&{ID:1 Nodes:[1 2 3]}
&{ID:3 Nodes:[3 4 5]}
&{ID:3 Nodes:[3 4 5]}
*/
