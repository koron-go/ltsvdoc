package ltsvdoc_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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

type testData interface {
	ltsvdoc.Marshaler
	ltsvdoc.Unmarshaler
}

func TestVarious(t *testing.T) {
	for ti, tc := range []struct {
		obj testData
		str string
	}{
		{&Foo{"aaa"}, "t:foo\tvalue:aaa\n"},
		{&Foo{"a\\b"}, "t:foo\tvalue:a\\\\b\n"},
	} {
		b, err := ltsvdoc.EncodeAll(tc.obj)
		if err != nil {
			t.Fatalf("encode failed #%d: %s", ti, err)
		}
		act, exp := string(b), tc.str
		if act != exp {
			t.Fatalf("ltsv not match #%d:\nwant=%q\n got=%q", ti, exp, act)
		}
		vals, err := ltsvdoc.DecodeAll([]byte(exp))
		if err != nil {
			t.Fatalf("decode failed #%d: %s", ti, err)
		}
		if diff := cmp.Diff([]ltsvdoc.Unmarshaler{tc.obj}, vals); diff != "" {
			t.Fatalf("decoded vals mismatch #%d: -want +got\n%s", ti, diff)
		}
	}
}
