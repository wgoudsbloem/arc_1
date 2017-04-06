package store

import (
	"arcessio/pubsub"
	"bytes"
	"io"
	"os"
)

type Storer interface {
	Put(r io.Reader) (offset int64, err error)
	Get(w io.Writer) (err error)
	pubsub.Subscriberer
}

type ReadAtWriteSeeker interface {
	io.ReadWriteSeeker
	io.ReaderAt
}

type store struct {
	stream ReadAtWriteSeeker
	end    int64
	pubsub.PubSub
}

func NewStorer(rw ReadAtWriteSeeker) Storer {
	e, err := rw.Seek(0, io.SeekEnd)
	if err != nil {
		panic(err)
	}
	s := &store{stream: rw, end: e}
	s.PubSub = &pubsub.PubSuber{}
	return s

}

func NewFileStorer(topic string) Storer {
	f, err := os.OpenFile(topic+".topic", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	return &store{stream: f}
}

func (s *store) Put(r io.Reader) (index int64, err error) {
	var bb bytes.Buffer
	_, err = bb.ReadFrom(r)
	if err != nil {
		return
	}
	err = bb.WriteByte('\n')
	if err != nil {
		return
	}
	n, err := s.stream.Write(bb.Bytes())
	index = s.end
	s.end += int64(n)
	if s.PubSub != nil {
		s.Notify(s)
	}
	return
}

func (s *store) Get(w io.Writer) (err error) {
	b := make([]byte, 1)
	for {
		_, err = s.stream.Read(b)
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}
		if b[0] != '\n' {
			_, err = w.Write(b)
			if err != nil {
				return
			}
		} else {
			break
		}
	}
	return
}

var buf int64 = 500

func lastEntry(sample []byte) (result []byte) {
	var start, end int
	for n := len(sample) - 1; n > -1; n-- {
		if sample[n] == '\n' {
			if end == 0 {
				end = n
			} else {
				start = n + 1
				break
			}
		}
	}
	result = sample[start:end]
	return
}
