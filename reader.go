package ltsvdoc

import (
	"bufio"
	"bytes"
	"io"
)

// Reader reads LTSV values.
type Reader struct {
	rd *bufio.Reader
}

// NewReader creates a new LTSV reader.
func NewReader(r io.Reader) *Reader {
	var rd *bufio.Reader
	switch r2 := r.(type) {
	case *bufio.Reader:
		rd = r2
	default:
		rd = bufio.NewReader(r)
	}
	return &Reader{rd: rd}
}

func (r *Reader) readLine() ([]byte, error) {
	d, err := r.rd.ReadSlice('\n')
	if err == nil || (err == io.EOF && len(d) > 0) {
		return d, nil
	} else if err != bufio.ErrBufferFull {
		return nil, err
	}
	bb := bytes.NewBuffer(d)
	for {
		d2, err := r.rd.ReadSlice('\n')
		if len(d2) > 0 {
			if _, err := bb.Write(d2); err != nil {
				return nil, err
			}
		}
		if err == nil || err == io.EOF {
			return bb.Bytes(), nil
		}
		if err != bufio.ErrBufferFull {
			return nil, err
		}
	}
}

// Read read a LTSV value.
func (r *Reader) Read() (Values, error) {
	d, err := r.readLine()
	if err != nil {
		return nil, err
	}
	d = bytes.TrimLeft(d, " \n\r\t")
	d = bytes.TrimRight(d, "\n\r\t")
	vals := make(Values, 0, 10)
	for _, raw := range bytes.Split(d, []byte("\t")) {
		kv := bytes.SplitN(raw, []byte(":"), 2)
		if len(kv) != 2 {
			continue
		}
		vals = append(vals, &LabelValue{
			Label: string(kv[0]),
			Value: string(kv[1]),
		})
	}
	return vals, nil
}
