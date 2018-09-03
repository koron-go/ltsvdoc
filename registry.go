package ltsvdoc

import (
	"fmt"
	"io"
	"reflect"
)

type docType struct {
	name string
	typ  reflect.Type
}

// Registry is a regisry of LTSV marshal/unmarshal-able types.
type Registry struct {
	xName map[string]*docType
	xType map[reflect.Type]*docType
}

// NewRegistry creates a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		xName: make(map[string]*docType),
		xType: make(map[reflect.Type]*docType),
	}
}

// Add adds a type as LTSV marshal/unmarshal-able.
func (reg *Registry) Add(name string, doc Doc) error {
	if _, ok := reg.xName[name]; ok {
		return fmt.Errorf("duplicated name: %s", name)
	}
	dt := &docType{name: name, typ: toType(doc)}
	reg.xName[dt.name] = dt
	reg.xType[dt.typ] = dt
	return nil
}

// getName returns name of type which implements Marshaler.
func (reg *Registry) getName(v Marshaler) (string, bool) {
	dt, ok := reg.xType[toType(v)]
	if !ok {
		return "", false
	}
	return dt.name, true
}

func (reg *Registry) newUnmarshaler(name string) (Unmarshaler, bool) {
	dt, ok := reg.xName[name]
	if !ok {
		return nil, false
	}
	return reflect.New(dt.typ).Interface().(Unmarshaler), true
}

func toType(v interface{}) reflect.Type {
	return reflect.Indirect(reflect.ValueOf(v)).Type()
}

// EncodeAll encodes all marshalers as LTSV to io.Writer.
func (reg *Registry) EncodeAll(w io.Writer, vals ...Marshaler) error {
	enc := NewEncoderWithRegistry(reg, w)
	for _, v := range vals {
		err := enc.Encode(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// DecodeAll decodes all unmarshalers from io.Reader.
func (reg *Registry) DecodeAll(r io.Reader) ([]Unmarshaler, error) {
	var retval []Unmarshaler
	dec := NewDecoderWithRegistry(reg, r)
	for {
		v, err := dec.Decode()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		retval = append(retval, v)
	}
	return retval, nil
}
