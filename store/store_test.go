package store

import (
	"arcessio/pubsub"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

type TestStream struct {
	bytes.Buffer
}

func (t *TestStream) Seek(offset int64, whence int) (int64, error) {
	return int64(t.Len()), nil
}

func (t *TestStream) ReadAt(p []byte, off int64) (n int, err error) {
	return t.Read(p)
}

// store a string into a buffer
func TestPut(t *testing.T) {
	in1 := "teststring"
	in2 := "teststring2"
	exp := in1 + "\n" + in2 + "\n"
	var bb TestStream
	s := store{stream: &bb}
	var bbin1 bytes.Buffer
	bbin1.WriteString(in1)
	_, err := s.Put(&bbin1)
	if err != nil {
		t.Error(err)
	}
	var bbin2 bytes.Buffer
	bbin2.WriteString(in2)
	offset, err := s.Put(&bbin2)
	if err != nil {
		t.Error(err)
	}
	expLen := int64(len(in1) + len("\n"))
	if offset != expLen {
		t.Errorf("want %v, got %v", expLen, offset)
	}
	res, _ := bb.ReadString('\t')
	if exp != res {
		t.Errorf("want %v, got %v", exp, res)
	}
}

func TestGet(t *testing.T) {
	in1 := "teststring"
	in2 := "teststring2"
	var bb TestStream
	s := store{stream: &bb}
	var bbin1 bytes.Buffer
	bbin1.WriteString(in1)
	_, err := s.Put(&bbin1)
	if err != nil {
		t.Error(err)
	}
	var bbin2 bytes.Buffer
	bbin2.WriteString(in2)
	offset, err := s.Put(&bbin2)
	if err != nil {
		t.Error(err)
	}
	var bbout1 bytes.Buffer
	err = s.Get(&bbout1)
	if err != nil {
		t.Error(err)
	}
	if bbout1.String() != in1 {
		t.Errorf("want %v got %v", in1, bbout1.String())
	}
	if bbout1.String() == in2 {
		t.Errorf("want %v got %v", in1, bbout1.String())
	}
	s2 := store{stream: &bb, end: offset}
	var bbout2 bytes.Buffer
	err = s2.Get(&bbout2)
	if err != nil {
		t.Error(err)
	}
	if bbout2.String() != in2 {
		t.Errorf("want [%v] got [%v]", in2, bbout2.String())
	}
}

func TestStorer(t *testing.T) {
	var w TestStream
	s := NewStorer(&w)
	_, ok := s.(Storer)
	if !ok {
		t.Error("Type is not a Storer")
	}
}

func TestFileStorer(t *testing.T) {
	fn := "test"
	s := NewFileStorer(fn)
	_, ok := s.(Storer)
	if !ok {
		t.Error("Type is not a Storer")
	}
	in1 := `{"test":"value"}`
	var bbin1 bytes.Buffer
	bbin1.WriteString(in1)
	_, err := s.Put(&bbin1)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(3 * time.Second)
	_, err = ioutil.ReadFile(fn + ".topic")
	if err != nil {
		t.Error(err)
	}
	s2 := NewFileStorer(fn)
	var bbout1 bytes.Buffer
	err = s2.Get(&bbout1)
	if err != nil {
		t.Error(err)
	}
	if string(in1) != bbout1.String() {
		t.Errorf("want %v got %v", in1, bbout1.String())
	}
}

type MockPubSub struct {
	t *testing.T
}

func (m *MockPubSub) Subscribe(fn pubsub.Subscriber) {

}

func (m *MockPubSub) Notify(in interface{}) {
	if _, ok := in.(Storer); !ok {
		m.t.Error("want Storer got something else...")
	}
}

func TestSubscribe(t *testing.T) {
	in1 := "teststring"
	in2 := "teststring2"
	var bb TestStream
	s := store{stream: &bb}
	s.PubSub = &MockPubSub{t}
	s.Subscribe(func(in interface{}) (err error) { return nil })
	var bbin1 bytes.Buffer
	bbin1.WriteString(in1)
	_, err := s.Put(&bbin1)
	if err != nil {
		t.Error(err)
	}
	var bbin2 bytes.Buffer
	bbin2.WriteString(in2)
	_, err = s.Put(&bbin2)
	if err != nil {
		t.Error(err)
	}

}

func TestInternalLastEntry(t *testing.T) {
	expVal1 := `{"test1":"val1"}`
	expVal2 := `{"test2":"val2"}`
	testVal1 := expVal1 + "\n"
	testVal2 := expVal2 + "\n"
	testVal3 := testVal1 + testVal2
	res1 := lastEntry([]byte(testVal1))
	if string(res1) != expVal1 {
		t.Errorf("want: '%v' got: '%v'", expVal1, string(res1))
	}
	res2 := lastEntry([]byte(testVal3))
	if string(res2) != expVal2 {
		t.Errorf("want: '%v' got: '%v'", expVal2, string(res2))
	}
}

func TestCleanup(t *testing.T) {
	fi, err := ioutil.ReadDir("./")
	if err != nil {
		t.Error(err)
	}
	for _, f := range fi {
		if strings.HasSuffix(f.Name(), ".topic") {
			os.Remove(f.Name())
		}
	}
}
