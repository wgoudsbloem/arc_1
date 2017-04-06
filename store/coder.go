package store

import (
	"encoding/base64"
	"io"
)

type Encoding struct {
	Encode
}

// Encoding
type Encoder func(w io.Writer)
type Encode func(w io.Writer) Encoder

func NewEncoder(e Encoding, w io.Writer) Encoder {
	return e.Encode(w)
}

var NonEncoding = Encoding{func(w1 io.Writer) Encoder {
	return func(w2 io.Writer) {
		return
	}
}}

var Base64Encoding = Encoding{func(w1 io.Writer) Encoder {
	return func(w2 io.Writer) {
		w2 = base64.NewEncoder(base64.StdEncoding, w1)
		return
	}
}}

// Decoding
type Decoding struct {
	d Decoder
}

type Decoder func(p []byte) (b []byte)

func NewDecoder(d Decoding) Decoder {
	return d.d
}

var NonDecoding = Decoding{d: func(p []byte) (b []byte) {
	return p
}}

var Base64Decoding = Decoding{d: func(p []byte) (b []byte) {
	base64.StdEncoding.Decode(b, p)
	return
}}
