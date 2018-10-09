package ltsvdoc

import (
	"fmt"
	"strconv"
)

// Proxy represents a proxy for Values.
type Proxy struct {
	vals   Values
	xLabel map[string][]int
	errors []error
}

// Proxy represents a proxy for LabelValue.
type ProxyValue struct {
	p   *Proxy
	l   string
	v   *LabelValue
	err error
}

// NewProxy creates a new Proxy for Values.
func NewProxy(vals Values) *Proxy {
	x := make(map[string][]int)
	for i, v := range vals {
		x[v.Label] = append(x[v.Label], i)
	}
	return &Proxy{
		vals:   vals,
		xLabel: x,
	}
}

// Error returns a first error occurrence.
func (p *Proxy) Error() error {
	if len(p.errors) == 0 {
		return nil
	}
	return p.errors[0]
}

func (p *Proxy) addErr(label string, err error) *ProxyValue {
	p.errors = append(p.errors, err)
	return &ProxyValue{
		p:   p,
		l:   label,
		err: err,
	}
}

// Get returns a proxy for a value at first.
func (p *Proxy) Get(label string) *ProxyValue {
	x, ok := p.xLabel[label]
	if !ok {
		return p.addErr(label, fmt.Errorf("label %q not found", label))
	}
	if len(x) > 1 {
		return p.addErr(label, fmt.Errorf("multiple labels for %q", label))
	}
	return &ProxyValue{
		p: p,
		l: label,
		v: p.vals[x[0]],
	}
}

// Opt return a proxy for a value without error propagation.
func (p *Proxy) Opt(label string) *ProxyValue {
	x, ok := p.xLabel[label]
	if !ok {
		return &ProxyValue{
			l:   label,
			err: fmt.Errorf("label %q not found", label),
		}
	}
	if len(x) > 1 {
		return &ProxyValue{
			l:   label,
			err: fmt.Errorf("multiple labels for %q", label),
		}
	}
	return &ProxyValue{
		l: label,
		v: p.vals[x[0]],
	}
}

// Has checks Values has a label or not.
func (p *Proxy) Has(label string) bool {
	x, ok := p.xLabel[label]
	if !ok {
		return false
	}
	if len(x) > 1 {
		return false
	}
	return true
}

func (pv *ProxyValue) setErr(err error) {
	if pv.err == nil {
		pv.err = err
		if pv.p != nil {
			pv.p.addErr(pv.l, err)
		}
	}
}

// Error gets a last occurred error.
func (pv *ProxyValue) Error() error {
	return pv.err
}

// Bool returns a bool value
func (pv *ProxyValue) Bool() bool {
	if pv.err != nil {
		return false
	}
	v, err := strconv.ParseBool(pv.v.RawValueString())
	if err != nil {
		pv.setErr(fmt.Errorf("label:%q failed to parse as bool: %s", pv.l, err))
		return false
	}
	return v
}

// String returns a string value.
func (pv *ProxyValue) String() string {
	if pv.err != nil {
		return ""
	}
	return pv.v.RawValueString()
}

// Int64 returns a int64 value.
func (pv *ProxyValue) Int64() int64 {
	if pv.err != nil {
		return 0
	}
	n, err := strconv.ParseInt(pv.v.RawValueString(), 10, 64)
	if err != nil {
		pv.setErr(fmt.Errorf("label:%q can't be parsed as int64: %s", pv.l, err))
		return 0
	}
	return n
}

// Int32 returns a int32 value.
func (pv *ProxyValue) Int32() int32 {
	if pv.err != nil {
		return 0
	}
	n, err := strconv.ParseInt(pv.v.RawValueString(), 10, 32)
	if err != nil {
		pv.setErr(fmt.Errorf("label:%q can't be parsed as int32: %s", pv.l, err))
		return 0
	}
	return int32(n)
}

// Float64 returns a float64 value.
func (pv *ProxyValue) Float64() float64 {
	if pv.err != nil {
		return 0
	}
	f, err := strconv.ParseFloat(pv.v.RawValueString(), 64)
	if err != nil {
		pv.setErr(fmt.Errorf("label:%q failed to parse as float64: %s", pv.l, err))
		return 0
	}
	return f
}
