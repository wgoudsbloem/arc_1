package stream

import (
	"bytes"
	"testing"
)

type ByteBuffer struct {
	b    bytes.Buffer
	seek int64
}

func (bb *ByteBuffer) Write(p []byte) (n int, err error) {
	n, err = bb.b.Write(p)
	bb.seek = bb.seek + int64(n)
	return
}

func (bb *ByteBuffer) Read(p []byte) (n int, err error) {
	return bb.b.Read(p)
}

func (bb *ByteBuffer) ReadAt(p []byte, off int64) (n int, err error) {
	return bb.b.Read(p)
}

// returns the length, which is the same as offset with one write
func (bb *ByteBuffer) Seek(offset int64, whence int) (int64, error) {
	return bb.seek, nil
}

func TestWriteByte(t *testing.T) {
	bb := ByteBuffer{}
	sw := NewWriter(&bb)
	offset, err := sw.WriteByteArray([]byte(`{"key":"value"}`))
	if err != nil {
		t.Error(err)
	}
	if offset < 1 {
		t.Errorf("expected offset to be > %v, but got %v", 1, offset)
	}
	//t.Logf("write offset: %v", offset)
}

func TestReadJson(t *testing.T) {
	bb := ByteBuffer{}
	sw := NewWriter(&bb)
	_, err := sw.WriteByteArray([]byte(`{"key":"value"}`))
	if err != nil {
		t.Error(err)
	}
	sr := NewReader(&bb, 0)
	b := make([]byte, bb.b.Len())
	offset, err := sr.Read(&b)
	if err != nil {
		t.Error(err)
	}
	if offset < 1 {
		t.Errorf("expected offset to be %v, but got %v", 1, offset)
	}
	//t.Logf("read offset: %v", offset)
	//t.Log(string(b))
}

func TestSubscribe(t *testing.T) {
	var b2 bytes.Buffer
	bb2 := ByteBuffer{b2, 0}
	sw := NewWriter(&bb2)
	fn := func(msg interface{}) error {
		//t.Log("executed")
		return nil
	}
	sw.Subscribe(fn)
	_, err := sw.WriteByteArray([]byte(`{"key":"value"}`))
	if err != nil {
		t.Error(err)
	}
}

func TestLastJson(t *testing.T) {
	exp0 := `{"keyx":"valuey"}`
	var b0 bytes.Buffer
	bb2 := ByteBuffer{b: b0}
	sw := NewWriter(&bb2)
	_, err := sw.WriteByteArray([]byte(exp0))
	if err != nil {
		t.Error(err)
	}
	sr := NewReader(&bb2, 0)
	b, _, err := sr.LastJSON()
	if err != nil {
		t.Error(err)
	}
	if string(b) != exp0 {
		t.Errorf("want: '%v' got '%v'", exp0, string(b))
	}
	//write another entry
	exp2 := `{"key1":"value2"}`
	_, err = sw.WriteByteArray([]byte(exp2))
	if err != nil {
		t.Error(err)
	}
	b2, _, err := sr.LastJSON()
	if err != nil {
		t.Error(err)
	}
	if string(b2) != exp2 {
		t.Errorf("want: '%v' got '%v'", exp2, string(b2))
	}
}

// in:"eyJrZXkxIjoidmFsdWUyIn0="

// out:"eyJrZXkxIjoidmFsdWUyIn0="
