package wire

import (
	"bytes"
	"io"
)

// registration multiple implementations for an interface, similar to go-data, but global
// multiple calls to Interface with the same base will return the same Registrar
func Interface(base interface{}) Registrar {

}

type Registrar struct {
}

func (r Registrar) RegisterImplementation(data interface{}, b byte) Registrar {

}

// if we want to read/write to bytes - eg. ReadBinaryBytes and BinaryBytes
func Marshal(v interface{}) ([]byte, error) {
	out := bytes.NewBuffer()
	enc := NewEncoder(out)
	err := enc.Encode(v)
	return out.Bytes(), err
}

func Unmarshal(data []byte, v interface{}) error {
	in := bytes.NewBuffer(data)
	dec := NewDecoder(in)
	err := enc.Decode(v)
	return err
}

// if we want to read/write to streams of data - eg. ReadBinary and WriteBinary
type Decoder struct{}

func NewDecoder(r io.Reader) *Decoder
func (dec *Decoder) Decode(v interface{}) error

type Encoder struct{}

func NewEncoder(w io.Writer) *Encoder
func (enc *Encoder) Encode(v interface{}) error
