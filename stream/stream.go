package stream

import (
	"encoding/json"
	"io"

	"arcessio/pubsub"
)

// Reader is the abstract
type Reader struct {
	ras ReadAtSeeker
	*json.Decoder
}

// Writer is the abstract
type Writer struct {
	ws io.WriteSeeker
	*json.Encoder
	pubsub.PubSub
}

// ReadWriter is the abstract
type ReadWriter struct {
	Reader
	Writer
}

// ReadAtSeeker is an composite interface
type ReadAtSeeker interface {
	io.ReadSeeker
	io.ReaderAt
}

// NewWriter is a factory to get a Writer
func NewWriter(ws io.WriteSeeker) Writer {
	enc := json.NewEncoder(ws)
	n := &pubsub.PubSuber{}
	return Writer{ws, enc, n}
}

// WriteByte will write a byte to the stream
func (sw *Writer) WriteByteArray(jsn []byte) (offset int64, err error) {
	if err = sw.Encode(jsn); err != nil {
		return
	}
	offset, err = sw.ws.Seek(0, io.SeekCurrent)
	sw.Notify(nil)
	return
}

// NewReader is a factory to get a Reader
func NewReader(rs ReadAtSeeker, offset int64) Reader {
	if _, err := rs.Seek(offset, io.SeekStart); err != nil {
		panic(err)
	}
	return Reader{rs, json.NewDecoder(rs)}
}

// WriteByte will write the stream to the byte
func (sr *Reader) Read(jsn *[]byte) (offset int64, err error) {
	if err = sr.Decode(jsn); err != nil {
		return
	}
	offset, err = sr.ras.Seek(0, io.SeekCurrent)
	return
}

const (
	buffSize int64 = 512
)

// LastJSON returns the last json from the stream
func (sr *Reader) LastJSON() (jsn []byte, end int64, err error) {
	end, err = sr.ras.Seek(0, io.SeekEnd)
	if err != nil {
		return
	}
	//get position to read from
	b := make([]byte, buffSize)
	buff := buffSize
	if buffSize > end {
		buff = end
	}
	_, err = sr.ras.ReadAt(b, end-buff)
	if err == io.EOF {
		//err = nil
	}
	err = json.Unmarshal(lastEntry(b), &jsn)
	return
}

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
