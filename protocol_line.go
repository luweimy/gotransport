package gotransport

import (
	"bufio"
	"io"
)

type lineProtocol struct {
	data []byte
}

func LineProtocol() Protocol {
	return &lineProtocol{}
}

func (l *lineProtocol) Payload() []byte {
	return l.data
}

func (l *lineProtocol) SetPayload(payload []byte) {
	l.data = payload
}

func (l *lineProtocol) SetFlags(value interface{}) error {
	return ErrFlagsNotSupport
}

func (l *lineProtocol) Flags() Flags {
	return FlagsValue{}
}

func (l *lineProtocol) WriteTo(w io.Writer) (int, error) {
	if len(l.data) <= 0 {
		return 0, nil
	}
	if l.data[len(l.data)-1] != '\n' {
		l.data = append(l.data, '\n')
	}
	return w.Write(l.data)
}

//func (p *lineProtocol) ReadFrom(r io.Reader) (int, error) {
//	var (
//		reader = bufio.NewReader(r)
//		buffer = bytes.Buffer{}
//
//		part   []byte
//		prefix bool
//		err    error
//	)
//	for {
//		if part, prefix, err = reader.ReadLine(); err != nil {
//			return buffer.Len(), err
//		}
//		buffer.Write(part)
//		if !prefix {
//			break
//		}
//	}
//
//	p.data = buffer.Bytes()
//	return len(p.data), nil
//}

func (p *lineProtocol) ReadFrom(r io.Reader) (int, error) {
	line, err := bufio.NewReader(r).ReadBytes('\n')
	n := len(line)
	if err != nil {
		return n, err
	}

	// Handle the case "\r\n".
	if len(line) > 0 && line[len(line)-1] == '\n' {
		drop := 1
		if len(line) > 1 && line[len(line)-2] == '\r' {
			drop = 2
		}
		line = line[:len(line)-drop]
	}

	p.data = line
	return n, nil
}
