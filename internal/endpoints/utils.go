package endpoints

import (
	"encoding/json"
	"io"
)

func jsonUnmarshalReader(r io.Reader, a any) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&a)

	if err != nil {
		return err
	}
	return nil
}
