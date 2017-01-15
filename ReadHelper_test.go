package nworder

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func ExampleReadHelper() {
	rh := ReadHelper{
		Reader: bytes.NewReader([]byte{0x01, 0x02, 0x03, 0x04}),
	}

	var buf1 uint16
	var buf2 = make([]byte, 2)
	var buf3 uint8

	rh.Read(&buf1, "r1")
	rh.Read(buf2, "r2")

	fmt.Printf("%04x", buf1)
	fmt.Println(hex.EncodeToString(buf2))
	fmt.Println(rh.Error)

	rh.Read(&buf3, "r3")

	fmt.Println(rh.Error)

	// Output:
	// 01020304
	// <nil>
	// r3: EOF
}

func TestReadHelper1(t *testing.T) {
	src := []byte("test")
	rh := ReadHelper{
		Reader: bytes.NewBuffer(src),
	}

	buf := make([]byte, 4)
	rh.Read(buf, "message")

	if rh.Error != nil {
		t.Fatal(rh.Error)
	}
	if !reflect.DeepEqual(src, buf) {
		t.Fatal(hex.EncodeToString(buf))
	}
}

func TestReadHelper2(t *testing.T) {
	src := []byte("test")
	rh := ReadHelper{
		Reader: bytes.NewBuffer(src),
	}

	buf := make([]byte, 5)
	rh.Read(buf, "message")

	if errors.Cause(rh.Error) != io.ErrUnexpectedEOF {
		t.Fatal(rh.Error)
	}
}

func TestReadHelper3(t *testing.T) {
	src := []byte{0x01, 0x02}
	rh := ReadHelper{
		Reader: bytes.NewBuffer(src),
	}

	var buf uint16
	rh.Read(&buf, "message")

	if rh.Error != nil {
		t.Fatal(rh.Error)
	}
	if buf != 0x0102 {
		t.Fatalf("%x", buf)
	}
}

func TestReadHelper4(t *testing.T) {
	src := []byte{0x01, 0x02}
	rh := ReadHelper{
		Reader: bytes.NewBuffer(src),
	}

	var buf uint32
	rh.Read(&buf, "message")

	if errors.Cause(rh.Error) != io.ErrUnexpectedEOF {
		t.Fatal(rh.Error)
	}
}

func TestReadHelper5(t *testing.T) {
	src := []byte{0x01, 0x02}
	rh := ReadHelper{
		Reader: bytes.NewBuffer(src),
		Error:  errors.New("dummy"),
	}

	var buf uint16
	rh.Read(&buf, "message")

	if l := rh.Reader.(*bytes.Buffer).Len(); l != 2 {
		t.Fatal(l)
	}
}

func TestReadHelper6(t *testing.T) {
	src := []byte{0x01, 0x02, 0x03, 0x04}
	rh := ReadHelper{
		Reader: bytes.NewBuffer(src),
	}

	var buf1, buf2 uint16
	rh.Read(&buf1, "message")

	rh.Do(func() error {
		buf1 = 0x0103
		return errors.New("invalid")
	})

	rh.Read(&buf2, "message")

	if errors.Cause(rh.Error).Error() != "invalid" {
		t.Fatal(rh.Error)
	}
	if buf1 != 0x0103 {
		t.Fatalf("%x", buf1)
	}
	if buf2 != 0 {
		t.Fatal(buf2)
	}
}
