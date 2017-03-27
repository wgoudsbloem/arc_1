package stream

import (
	"os"
	"testing"
)

var lastOffset int64

var testTopic = "filestream.test"

var (
	expectedValue1 = []byte(`{"key_test":"value_test"}`)
	expectedValue2 = []byte(`{"key_test2":"value_test2"}`)
	expectedValue3 = []byte(`{"key_test3":"value_test3"}`)
)

var expectedOffset int64

func TestNewFileStreamWriter(t *testing.T) {
	sw, err := NewFileStreamWriter(testTopic)
	if err != nil {
		t.Fatal(err)
	}
	offset, err := sw.WriteByteArray(expectedValue1)
	if err != nil {
		t.Error(err)
	}
	if offset < 1 {
		t.Errorf("expected offset to be %v, but got %v", 1, offset)
	}
	//t.Logf("write file offset: %v", offset)
	offset, err = sw.WriteByteArray(expectedValue2)
	if err != nil {
		t.Error(err)
	}
	if offset < 1 {
		t.Errorf("expected offset to be %v, but got %v", 1, offset)
	}
	//t.Logf("write file offset: %v", offset)
	lastOffset = offset
	offset, err = sw.WriteByteArray(expectedValue3)
	if err != nil {
		t.Error(err)
	}
	if offset < 1 {
		t.Errorf("expected offset to be %v, but got %v", 1, offset)
	}
	//t.Logf("write file offset: %v", offset)
	expectedOffset = offset
}

func TestNewFileStreamReader(t *testing.T) {
	sr, err := NewFileStreamReader(testTopic)
	if err != nil {
		t.Fatal(err)
	}
	var b []byte
	offset, err := sr.Read(&b)
	if err != nil {
		t.Error(err)
	}
	if offset < 1 {
		t.Errorf("expected offset to be %v, but got %v", 1, offset)
	}
	//t.Logf("read file offset: %v", offset)
	//t.Log(string(b))
	var b1 []byte
	offset, err = sr.Read(&b1)
	if err != nil {
		t.Error(err)
	}
	if offset < 1 {
		t.Errorf("expected offset to be %v, but got %v", 1, offset)
	}
	//t.Logf("read2 file offset: %v", offset)
	//t.Log(string(b1))
}

func TestNewFileStreamReaderAt(t *testing.T) {
	sr, err := NewFileStreamReaderAt(testTopic, lastOffset)
	if err != nil {
		t.Fatal(err)
	}
	var b []byte
	offset, err := sr.Read(&b)
	if err != nil {
		t.Error(err)
	}
	if offset < 1 {
		t.Errorf("expected offset to be %v, but got %v", 1, offset)
	}
	//t.Logf("read file offset: %v", offset)
	//t.Log(string(b))
}

var testtopic2 = "testtopic2"

func TestFileLastJson(t *testing.T) {
	srw, err := NewFileStreamReaderWriter(testtopic2)
	if err != nil {
		t.Fatal(err)
	}
	expVal := []byte(`{"test","val"}`)
	srw.WriteByteArray(expVal)
	b, _, err := srw.LastJSON()
	if err != nil {
		t.Error(err)
	}
	if string(b) != string(expVal) {
		t.Errorf("expected b to be %v, but got %v", string(expVal), string(b))
	}
	//	if offset != expectedOffset {
	//		t.Errorf("expected offset to be %v, but got %v", expectedOffset, offset)
	//	}
	//t.Log(offset)
}

func TestCleanup2(t *testing.T) {
	if err := os.Remove(testTopic + ".topic"); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(testtopic2 + ".topic"); err != nil {
		t.Fatal(err)
	}
}
