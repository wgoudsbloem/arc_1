package store

import (
	"bytes"
	"io"
	"testing"
)

func TestNonCoder(t *testing.T) {
	exp := "test"
	var in, out, in2, out2, out3 bytes.Buffer
	in.WriteString(exp)
	enc := NewEncoder(NonEncoding, &in)
	enc(&out)
	val, err := out.ReadString('\n')
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	if val != exp {
		t.Errorf("want %v got %v", exp, val)
	}
	in2.WriteString(exp)
	enc2 := NewEncoder(NonEncoding, &in2)
	enc2(&out2)
	dec := NewDecoder(NonDecoding, &out2)
	dec(&out3)
	val2, err := out3.ReadString('\n')
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	if val2 != exp {
		t.Errorf("want %v got %v", exp, val2)
	}
}

func TestBase64Coder(t *testing.T) {
	exp := "test"
	var in, out, in2, out2, out3 bytes.Buffer
	in.WriteString(exp)
	enc := NewEncoder(Base64Encoding, &in)
	enc(&out)
	val, err := out.ReadString('\n')
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	if val == "" {
		t.Error("val cannot be empty")
	}
	if val == exp {
		t.Errorf("do not want %v but got %v", exp, val)
	}
	in2.WriteString(exp)
	enc2 := NewEncoder(NonEncoding, &in2)
	enc2(&out2)
	dec := NewDecoder(NonDecoding, &out2)
	dec(&out3)
	val2, err := out3.ReadString('\n')
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	if val2 != exp {
		t.Errorf("want %v got %v", exp, val2)
	}
}
