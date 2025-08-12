package copy

import (
	"bytes"
	"encoding/gob"
)

// DeepCopyGob performs a deep copy of src into dst using gob encoding
// It can't copy functions or channels.
func DeepCopyGob(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
