package nwbyteorder

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

// WriteHelper is help to write on network order.
type WriteHelper struct {
	Error  error
	Writer io.Writer
}

// Write skips when Error != nil.
//
// Write is same as binary.Write with binary.BigEndian.
// when Write failed, Error set by errors.Wrapf with msg and msgargs.
func (s *WriteHelper) Write(v interface{}, msg string, msgargs ...interface{}) {
	if s.Error != nil {
		return
	}

	s.Error = binary.Write(s.Writer, binary.BigEndian, v)
	s.Error = errors.Wrapf(s.Error, msg, msgargs...)
}

// Do skips when Error != nil.
//
// Do calls f(), and updates Error
func (s *WriteHelper) Do(f func() error) {
	if s.Error != nil {
		return
	}

	s.Error = f()
}
