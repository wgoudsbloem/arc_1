package store

import "encoding/base64"

type Encoding struct {
	Encoder
}

// Encoding
type Encoder func(b []byte) (p []byte)

func NewEncoder(e Encoding) Encoder {
	return e.Encoder
}

var NonEncoding = Encoding{func(b []byte) (p []byte) {
	return b
}}

var Base64Encoding = Encoding{func(b []byte) (p []byte) {
	base64.StdEncoding.Encode(p, b)
	return
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
