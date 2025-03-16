package main

import (
	"net"
	"sync"
)

type listener struct {
	net.Listener
}

func (a listener) Accept() (net.Conn, error) {
	conn, err := a.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return &connection{
		Conn: conn,
	}, nil
}

type connection struct {
	net.Conn
	lock sync.Mutex
	buf  *[]byte
	n    int
}

func (a *connection) Read(b []byte) (n int, err error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if len(b) == 0 {
		bp := make([]byte, 1)
		n, err = a.Conn.Read(bp)
		if err == nil && n > 0 {
			if a.buf == nil {
				a.buf = &bp
			} else {
				*a.buf = append(*a.buf, bp[0])
			}
			a.n += n
		}
		return 0, err
	}
	if a.n > 0 {
		n = copy(b, *a.buf)
		*a.buf = (*a.buf)[n:]
		a.n -= n
		return n, nil
	}
	return a.Conn.Read(b)
}
