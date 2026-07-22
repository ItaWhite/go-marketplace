package domain

import "encoding/json"

type Nullable[T any] struct {
	Set   bool
	Value *T
}

func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true

	if string(b) == "null" {
		n.Value = nil
		return nil
	}

	var value T

	err := json.Unmarshal(b, &value)
	if err != nil {
		return err
	}

	n.Value = &value

	return nil
}
