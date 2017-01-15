package nwbyteorder

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func ExampleWriteHelper() {
	buf := new(bytes.Buffer)
	wh := &WriteHelper{
		Writer: buf,
	}

	wh.Write([]byte{0x01, 02}, "w1")
	wh.Write(uint16(0x0304), "w2")

	fmt.Println(wh.Error)
	fmt.Println(hex.EncodeToString(buf.Bytes()))

	// Output:
	// <nil>
	// 01020304
}

func TestWriteHelper1(t *testing.T) {
	src := []byte("test")
	wh := &WriteHelper{
		Writer: new(bytes.Buffer),
	}

	wh.Write(src, "message")

	if wh.Error != nil {
		t.Fatal(wh.Error)
	}
	if buf := wh.Writer.(*bytes.Buffer).Bytes(); !reflect.DeepEqual(src, buf) {
		t.Fatal(hex.EncodeToString(buf))
	}
}

func TestWriteHelper2(t *testing.T) {
	src := uint32(0x01020304)
	wh := &WriteHelper{
		Writer: new(bytes.Buffer),
	}

	wh.Write(src, "message")

	if wh.Error != nil {
		t.Fatal(wh.Error)
	}
	if buf := wh.Writer.(*bytes.Buffer).Bytes(); !reflect.DeepEqual([]byte{0x01, 0x02, 0x03, 0x04}, buf) {
		t.Fatal(hex.EncodeToString(buf))
	}
}

func TestWriteHelper3(t *testing.T) {
	src := uint32(0x01020304)
	wh := &WriteHelper{
		Writer: new(bytes.Buffer),
		Error:  errors.New("dummy"),
	}

	wh.Write(src, "message")

	if l := wh.Writer.(*bytes.Buffer).Len(); l != 0 {
		t.Fatal(l)
	}
}

func TestWriteHelper4(t *testing.T) {
	val1 := []byte{0x01, 0x02, 0x03, 0x04}
	val2 := uint16(0x0506)
	wh := WriteHelper{
		Writer: new(bytes.Buffer),
	}

	wh.Write(val1, "message")

	wh.Do(func() error {
		val2 = uint16(0x0708)
		return errors.New("invalid")
	})

	wh.Write(val2, "message")

	if errors.Cause(wh.Error).Error() != "invalid" {
		t.Fatal(wh.Error)
	}
	if buf := wh.Writer.(*bytes.Buffer).Bytes(); !reflect.DeepEqual(val1, buf) {
		t.Fatal(hex.EncodeToString(buf))
	}
	if val2 != uint16(0x0708) {
		t.Fatalf("%x", val2)
	}
}
