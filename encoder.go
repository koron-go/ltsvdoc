package ltsvdoc

import (
	"fmt"
	"io"
)

// Encoder encodes Marshalers as LTSV lines.
type Encoder struct {
	reg *Registry
	w   io.Writer
}

// NewEncoder cerates a new Encoder using io.Writer.
func NewEncoder(w io.Writer) *Encoder {
	return NewEncoderWithRegistry(DefaultRegister, w)
}

// NewEncoderWithRegistry creates a new Encoder with Registry.
func NewEncoderWithRegistry(reg *Registry, w io.Writer) *Encoder {
	return &Encoder{
		reg: reg,
		w:   w,
	}
}

// Encode encodes a Marshaler
func (e *Encoder) Encode(v Marshaler) error {
	n, ok := e.reg.getName(v)
	if !ok {
		return fmt.Errorf("not found the type:%T in Registry", v)
	}
	vals, err := v.MarshalLTSV()
	if err != nil {
		return err
	}
	row := make(Values, 0, len(vals)+1)
	row = append(row, &LabelValue{Label: "t", Value: n})
	row = append(row, vals...)
	_, err = row.WriteTo(e.w)
	if err != nil {
		return err
	}
	return nil
}
