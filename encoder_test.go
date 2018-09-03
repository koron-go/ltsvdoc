package ltsvdoc_test

import (
	"testing"

	"github.com/koron-go/ltsvdoc"
)

func TestEncode_SimpleSingle(t *testing.T) {
	b, err := ltsvdoc.EncodeAll(&Foo{"abc"})
	if err != nil {
		t.Fatal("EncodeAll failed:", err)
	}
	s := string(b)
	if s != "t:foo\tvalue:abc\n" {
		t.Fatalf("unexpected EncodeAll: %q", s)
	}
}

func TestEncode_SimpleMultiple(t *testing.T) {
	b, err := ltsvdoc.EncodeAll(
		&Foo{"aaa"},
		&Foo{"bbb"},
		&Foo{"ccc"},
		&Foo{"zzz"},
	)
	if err != nil {
		t.Fatal("failed EncodeAll:", err)
	}
	s := string(b)
	exp := `t:foo	value:aaa
t:foo	value:bbb
t:foo	value:ccc
t:foo	value:zzz
`
	if s != exp {
		t.Fatalf("unexpected EncodeAll: %q", s)
	}
}

func TestEncode_ComplexMultiple(t *testing.T) {
	b, err := ltsvdoc.EncodeAll(
		&Foo{"aaa"},
		&Bar{},
		&Baz{123},
		&Foo{"zzz"},
	)
	if err != nil {
		t.Fatal("failed EncodeAll:", err)
	}
	s := string(b)
	exp := `t:foo	value:aaa
t:bar
t:baz	integer:123
t:foo	value:zzz
`
	if s != exp {
		t.Fatalf("unexpected EncodeAll: %q", s)
	}
}
