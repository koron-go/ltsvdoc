package ltsvdoc

import (
	"fmt"
	"io"
)

// Decoder provides LTSV doc decoding.
type Decoder struct {
	reg *Registry
	r   *Reader
}

// NewDecoder creates a new decoder.
func NewDecoder(r io.Reader) *Decoder {
	return NewDecoderWithRegistry(DefaultRegister, r)
}

// NewDecoderWithRegistry creates a new decoder with Registry.
func NewDecoderWithRegistry(reg *Registry, r io.Reader) *Decoder {
	return &Decoder{
		reg: reg,
		r:   NewReader(r),
	}
}

// Decode decodes a document.
func (d Decoder) Decode() (Unmarshaler, error) {
	row, err := d.r.Read()
	if err != nil {
		return nil, err
	}
	first, data := row[0], row[1:]
	if first.Label != "t" {
		return nil, fmt.Errorf(`first value should have label "t": %q`, first.Label)
	}
	n := first.RawValueString()
	um, ok := d.reg.newUnmarshaler(n)
	if !ok {
		return nil, fmt.Errorf("unsupported type: %s", n)
	}
	err = um.UnmarshalLTSV(data)
	if err != nil {
		return nil, err
	}
	return um, nil
}
