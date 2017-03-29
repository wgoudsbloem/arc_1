package store

import (
	"bytes"
	"io"
	"os"
)

type Storer interface {
	Put(p []byte) (offset int64, err error)
	Get() (p []byte, err error)
}

type ReadAtWriteSeeker interface {
	io.ReadWriteSeeker
	io.ReaderAt
}

type store struct {
	stream ReadAtWriteSeeker
	end    int64
}

func NewStorer(rw ReadAtWriteSeeker) Storer {
	e, err := rw.Seek(0, io.SeekEnd)
	if err != nil {
		panic(err)
	}
	return &store{rw, e}
}

func NewFileStorer(topic string) Storer {
	f, err := os.OpenFile(topic+".topic", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	return &store{stream: f}
}

func (s *store) Put(p []byte) (offset int64, err error) {
	p = append(p, '\n')
	n, err := s.stream.Write(p)
	if err != nil {
		return
	}
	s.end += int64(n)
	offset = s.end
	return
}

func (s *store) Get() (p []byte, err error) {
	b := make([]byte, 1)
	var bb bytes.Buffer
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
			_, err = bb.Write(b)
			if err != nil {
				return
			}
		} else {
			break
		}
	}
	p = bb.Bytes()
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