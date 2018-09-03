package ltsvdoc_test

import (
	"fmt"
	"strconv"

	"github.com/koron-go/ltsvdoc"
)

type Foo struct {
	Value string
}

func (p *Foo) MarshalLTSV() (ltsvdoc.Values, error) {
	return ltsvdoc.NewValues("value", p.Value), nil
}

func (p *Foo) UnmarshalLTSV(vals ltsvdoc.Values) error {
	for _, v := range vals {
		switch v.Label {
		case "value":
			p.Value = v.ValueString()
		default:
			return fmt.Errorf("unknown label: %q", v.Label)
		}
	}
	return nil
}

type Bar struct {
}

func (p *Bar) MarshalLTSV() (ltsvdoc.Values, error) {
	return nil, nil
}

func (p *Bar) UnmarshalLTSV(vals ltsvdoc.Values) error {
	return nil
}

type Baz struct {
	Integer int
}

func (p *Baz) MarshalLTSV() (ltsvdoc.Values, error) {
	return ltsvdoc.NewValues("integer", p.Integer), nil
}

func (p *Baz) UnmarshalLTSV(vals ltsvdoc.Values) error {
	for _, v := range vals {
		switch v.Label {
		case "integer":
			n, err := strconv.Atoi(v.ValueString())
			if err != nil {
				return err
			}
			p.Integer = n
		default:
			return fmt.Errorf("unknown label: %q", v.Label)
		}
	}
	return nil
}

func init() {
	ltsvdoc.Register("foo", &Foo{})
	ltsvdoc.Register("bar", &Bar{})
	ltsvdoc.Register("baz", &Baz{})
}
