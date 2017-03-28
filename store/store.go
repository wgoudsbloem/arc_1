package store

import (
	"io"
)

type Storer interface {
	Put(p []byte) (offset int64, err error)
}

type store struct {
	ws  io.WriteSeeker
	end int64
}

func NewStorer(ws io.WriteSeeker) Storer {
	e, err := ws.Seek(0, io.SeekEnd)
	if err != nil {
		panic(err)
	}
	return store{ws, e}
}

func (s store) Put(p []byte) (offset int64, err error) {
	n, err := s.ws.Write(p)
	if err != nil {
		return
	}
	offset = s.end + int64(n)
	return
}
