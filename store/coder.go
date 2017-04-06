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
type Encode func(w io.Reader) Encoder

func NewEncoder(e Encoding, r io.Reader) Encoder {
	return e.Encode(r)
}

var NonEncoding = Encoding{func(r io.Reader) Encoder {
	return func(w io.Writer) {
		io.Copy(w, r)
	}
}}

var Base64Encoding = Encoding{func(r io.Reader) Encoder {
	return func(w io.Writer) {
		wout := base64.NewEncoder(base64.StdEncoding, w)
		io.Copy(wout, r)
	}
}}

// Decoding
type Decoding struct {
	Decode
}

type Decoder func(w io.Writer)
type Decode func(w io.Reader) Decoder

func NewDecoder(d Decoding, r io.Reader) Decoder {
	return d.Decode(r)
}

var NonDecoding = Decoding{func(r io.Reader) Decoder {
	return func(w io.Writer) {
		io.Copy(w, r)
	}
}}

var Base64Decoding = Decoding{func(r io.Reader) Decoder {
	return func(w io.Writer) {
		rout := base64.NewDecoder(base64.StdEncoding, r)
		io.Copy(w, rout)
	}
}}
