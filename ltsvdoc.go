package ltsvdoc

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// LabelValue is a pair of label and value.
type LabelValue struct {
	Label string
	Value interface{}
}

// NewLabelValue creates a labeled value.
func NewLabelValue(label string, value interface{}) *LabelValue {
	return &LabelValue{
		Label: label,
		Value: value,
	}
}

// ValueString get a string representation of value.
func (v LabelValue) ValueString() string {
	s := fmt.Sprint(v.Value)
	if strings.IndexAny(s, "\t\n\\") < 0 {
		return s
	}
	q := strconv.Quote(s)
	return q[1 : len(q)-1]
}

// RawValueString return value in string form without any quotes.
func (v LabelValue) RawValueString() string {
	return fmt.Sprint(v.Value)
}

// Values is an array of LabelValue
type Values []*LabelValue

// NewValues creates a Values with capacity.
func NewValues(args ...interface{}) Values {
	if len(args)%2 == 1 {
		panic("NewValues require even args")
	}
	vals := make(Values, 0, len(args)/2)
	for i := 0; i+1 < len(args); i += 2 {
		vals.Add(fmt.Sprint(args[i]), args[i+1])
	}
	return vals
}

// Add addes a value to Values.
func (vals *Values) Add(label string, value interface{}) *Values {
	*vals = append(*vals, &LabelValue{Label: label, Value: value})
	return vals
}

// Append adds multiple values to Values.
func (vals *Values) Append(args ...interface{}) *Values {
	if len(args)%2 == 1 {
		panic("Append require even args")
	}
	for i := 0; i+1 < len(args); i += 2 {
		vals.Add(fmt.Sprint(args[i]), args[i+1])
	}
	return vals
}

// WriteTo writes LTSV string to io.Writer
func (vals Values) WriteTo(w io.Writer) (int64, error) {
	b := make([]byte, 0, 1024)
	for i, v := range vals {
		if i > 0 {
			b = append(b, '\t')
		}
		b = append(b, v.Label...)
		b = append(b, ':')
		b = append(b, v.ValueString()...)
	}
	b = append(b, '\n')
	n, err := w.Write(b)
	if err != nil {
		return 0, err
	}
	return int64(n), nil
}

// Marshaler defines LTSV row marshaler
type Marshaler interface {
	MarshalLTSV() (Values, error)
}

// Unmarshaler defines LTSV marshaler
type Unmarshaler interface {
	UnmarshalLTSV(Values) error
}

// Doc defines required functions for LTSV marshal/unmarshal-able type.
type Doc interface {
	Marshaler
	Unmarshaler
}

// DefaultRegister is a Registry which used as default in this package.
var DefaultRegister = NewRegistry()

// Register registers a type as LTSV marshal/unmarshal-able.
func Register(name string, doc Doc) {
	err := DefaultRegister.Add(name, doc)
	if err != nil {
		panic(err.Error())
	}
}

// EncodeAll encodes all marshalers as LTSV to []byte.
func EncodeAll(vals ...Marshaler) ([]byte, error) {
	bb := &bytes.Buffer{}
	err := DefaultRegister.EncodeAll(bb, vals...)
	if err != nil {
		return nil, err
	}
	return bb.Bytes(), nil
}

// DecodeAll decodes all unmarshalers from []byte.
func DecodeAll(b []byte) ([]Unmarshaler, error) {
	r := bytes.NewReader(b)
	return DefaultRegister.DecodeAll(r)
}
