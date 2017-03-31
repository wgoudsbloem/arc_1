package store

import "testing"

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
