package ltsvdoc_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/koron-go/ltsvdoc"
)

func TestDecode_SimpleSingle(t *testing.T) {
	in := "t:foo\tvalue:abc\n"
	vals, err := ltsvdoc.DecodeAll([]byte(in))
	if err != nil {
		t.Fatal("failed DecodeAll:", err)
	}
	exp := []ltsvdoc.Unmarshaler{
		&Foo{Value: "abc"},
	}
	if diff := cmp.Diff(exp, vals); diff != "" {
		t.Fatal("unexpected DecodeAll:\n", diff)
	}
}

func TestDecode_SimpleMultiple(t *testing.T) {
	in := `t:foo	value:aaa
t:foo	value:bbb
t:foo	value:ccc
t:foo	value:zzz
`
	vals, err := ltsvdoc.DecodeAll([]byte(in))
	if err != nil {
		t.Fatal("failed DecodeAll:", err)
	}
	exp := []ltsvdoc.Unmarshaler{
		&Foo{"aaa"},
		&Foo{"bbb"},
		&Foo{"ccc"},
		&Foo{"zzz"},
	}
	if diff := cmp.Diff(exp, vals); diff != "" {
		t.Fatal("unexpected DecodeAll:\n", diff)
	}
}

func TestDecode_ComplexMultiple(t *testing.T) {
	in := `t:foo	value:aaa
t:bar
t:baz	integer:123
t:foo	value:zzz
`
	vals, err := ltsvdoc.DecodeAll([]byte(in))
	if err != nil {
		t.Fatal("failed DecodeAll:", err)
	}
	exp := []ltsvdoc.Unmarshaler{
		&Foo{"aaa"},
		&Bar{},
		&Baz{123},
		&Foo{"zzz"},
	}
	if diff := cmp.Diff(exp, vals); diff != "" {
		t.Fatal("unexpected DecodeAll:\n", diff)
	}
}
