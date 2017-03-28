package store

import (
	"bytes"
	"testing"
)

type TestSeeker struct {
	bytes.Buffer
}

func (t *TestSeeker) Seek(offset int64, whence int) (int64, error) {
	return int64(t.Len()), nil
}

// store a string into a buffer
func TestPut(t *testing.T) {
	exp := []byte("teststring")
	var bb TestSeeker
	s := store{&bb, 0}
	offset, err := s.Put(exp)
	if err != nil {
		t.Error(err)
	}
	expLen := int64(len(exp))
	if offset != expLen {
		t.Errorf("want %v, got %v", expLen, offset)
	}
	res, _ := bb.ReadBytes('\t')
	if exp != res {
		t.Errorf("want %v, got %v", exp, res)
	}
}
