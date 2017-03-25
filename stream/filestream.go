package stream

import (
	"os"
)

// NewFileStreamWriter is a factory method to get a StreamWriter
func NewFileStreamWriter(topic string) (sw Writer, err error) {
	topic += ".topic"
	f, err := os.OpenFile(topic, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	sw = NewWriter(f)
	return
}

// NewFileStreamReader is a factory method to get a StreamReader
func NewFileStreamReader(topic string) (sr Reader, err error) {
	topic += ".topic"
	f, err := os.OpenFile(topic, os.O_RDONLY, 0666)
	if err != nil {
		return
	}
	sr = NewReader(f, 0)
	return
}

// NewFileStreamReaderWriter is a factory method to get a StreamReadWriter
func NewFileStreamReaderWriter(topic string) (srw ReadWriter, err error) {
	topic += ".topic"
	f, err := os.OpenFile(topic, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	srw = ReadWriter{NewReader(f, 0), NewWriter(f)}
	return
}

// NewFileStreamReaderAt is a factory method to get a StreamReader with a
// determined start point (at)
func NewFileStreamReaderAt(topic string, offset int64) (sr Reader, err error) {
	topic += ".topic"
	f, err := os.OpenFile(topic, os.O_RDONLY, 0666)
	if err != nil {
		return
	}
	sr = NewReader(f, offset)
	return
}
