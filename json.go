package indexmap

import "encoding/json"

func (imap *IndexMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(imap.primaryIndex.inner)
}

func (imap *IndexMap[K, V]) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &imap.primaryIndex.inner); err != nil {
		return err
	}

	imap.primaryIndex.iterate(func(_ K, value *V) {
		imap.Insert(value)
	})

	return nil
}
