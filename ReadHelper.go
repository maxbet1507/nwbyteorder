package nwbyteorder

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

// ReadHelper is help to read on network order.
type ReadHelper struct {
	Error  error
	Reader io.Reader
}

// Read skips when Error != nil.
//
// Read is same as binary.Read with binary.BigEndian.
// when Read failed, Error set by errors.Wrapf with msg and msgargs.
func (s *ReadHelper) Read(v interface{}, msg string, msgargs ...interface{}) {
	if s.Error != nil {
		return
	}

	s.Error = binary.Read(s.Reader, binary.BigEndian, v)
	s.Error = errors.Wrapf(s.Error, msg, msgargs...)
}

// Do skips when Error != nil.
//
// Do calls f(), and updates Error
func (s *ReadHelper) Do(f func() error) {
	if s.Error != nil {
		return
	}

	s.Error = f()
}
