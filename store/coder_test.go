package store

import (
	"bytes"
	"testing"
)

func TestNonCoder(t *testing.T) {
	var bbin bytes.Buffer
	enc := NewEncoder(NonEncoding)

}

// func TestBase64Coder(t *testing.T) {
// 	in := []byte("testval")
// 	var o []byte
// 	base64.StdEncoding.Encode(o, in)
// 	t.Log(o)
// 	e := NewEncoder(Base64Encoding)
// 	out1 := e(in)
// 	if string(out1) == string(in) {
// 		t.Errorf("want %v got %v", string(in), string(out1))
// 	}
// 	d := NewDecoder(NonDecoding)
// 	out2 := d(out1)
// 	if string(out2) != string(in) {
// 		t.Errorf("want %v got %v", string(in), string(out2))
// 	}
// }
