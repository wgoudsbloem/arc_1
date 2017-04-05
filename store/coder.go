package store

import (
	"bytes"
	"encoding/base64"
	"io"
)

type Encoding struct {
	e Encoder
}

type Encoder interface {
	Encode(b []byte)
	io.Writer
}

func Newencoder(e Encoding) Encoder {
	return e.e
}

var NonEncoding = Encoding{e: NewNonEncoder()}
var Base64Encoding = Encoding{e: NewBase64Encoder()}

func NewNonEncoder() *NonEncoder {
	var bb bytes.Buffer
	return &NonEncoder{&bb}
}

type NonEncoder struct {
	io.Writer
}

func (ne *NonEncoder) Encode(b []byte) {
	ne.Write(b)
}

func NewBase64Encoder() *Base64Encoder {
	var bb bytes.Buffer
	base64.NewEncoder(base64.StdEncoding, &bb)
	return &Base64Encoder{&bb}
}

type Base64Encoder struct {
	io.Writer
}

func (b64 *Base64Encoder) Encode(b []byte) {
	b64.Write(b)
}

type Decoder interface {
	Decode(b []byte)
	io.Reader
}

func NewDecoder(d Decoding) Decoder {
	return d.d
}

type Decoding struct {
	d Decoder
}

var NonDecoding = Decoding{d: NewNonDecoder()}

func NewNonDecoder() *NonDecoder {
	return &NonDecoder{}
}

type NonDecoder struct {
	io.Reader
}

func (n *NonDecoder) Decode(b []byte) {
	for {
		if _, err := n.Read(b); err != nil {
			if err == io.EOF {
				err = nil
				return
			}
			return
		}
	}
}
