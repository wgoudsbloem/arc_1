package store

import "testing"
import "bytes"

func TestCoder(t *testing.T) {
	enc := Newencoder(NonEncoding)
	if _, ok := enc.(Encoder); !ok {
		t.Fatal("Expected NewEncoder to return type Encoder")
	}
	in1 := []byte("testval")
	enc.Encode(in1)
	if _, ok := enc.((*NonEncoder)); !ok {
		t.Error("Expected NewEncoder to return concrete type NonEncoder")
	}
	enc2 := Newencoder(Base64Encoding)
	if _, ok := enc.(Encoder); !ok {
		t.Fatal("Expected NewEncoder to return type Encoder")
	}
	in2 := []byte("testval")
	enc2.Encode(in2)
	if _, ok := enc2.((*Base64Encoder)); !ok {
		t.Error("Expected NewEncoder to return concrete type Base64Encoder")
	}
}

func TestDecoder(t *testing.T) {
	dec := NewDecoder(NonDecoding)
	if _, ok := dec.(Decoder); !ok {
		t.Fatal("Expected NewDecoder to return type Decoder")
	}
	var bb bytes.Buffer
	enc := NonEncoder{&bb}
	exp := []byte("testval")
	enc.Encode(exp)
	dec2 := NonDecoder{&bb}
	var b []byte
	dec2.Decode(b)
	if string(exp) != string(b) {
		t.Errorf("want %v got %v", string(exp), string(b))
	}
}
