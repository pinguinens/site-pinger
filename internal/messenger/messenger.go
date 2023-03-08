package messenger

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

type Messenger struct {
	conn  net.Conn
	codes []string
}

func New(address string, codes []string) (*Messenger, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Messenger{
		conn:  conn,
		codes: codes,
	}, nil
}

func (m *Messenger) Close() error {
	return m.conn.Close()
}

func (m *Messenger) Send(msg string) error {
	fmt.Fprintf(m.conn, "%v", msg)
	rBytes := make([]byte, 128)
	_, err := m.conn.Read(rBytes)
	if err != nil {
		return err
	}
	if string(rBytes) != "OK" {
		return errors.New(fmt.Sprintf("message server not OK: %v", string(rBytes)))
	}

	return nil
}

func (m *Messenger) Alarm(code int, method, url, addr string) error {
	if m.isAlarm(code) {
		return m.Send(fmt.Sprintf("%v\n%v\n%v\n%v", code, method, url, addr))
	}

	return nil
}

func (m *Messenger) isAlarm(code int) bool {
	for _, c := range m.codes {
		switch c {
		case "1**":
			if !(code >= 100 && code <= 199) {
				return false
			}
		case "2**":
			if !(code >= 200 && code <= 299) {
				return false
			}
		case "3**":
			if !(code >= 300 && code <= 399) {
				return false
			}
		case "4**":
			if !(code >= 400 && code <= 499) {
				return false
			}
		case "5**":
			if !(code >= 500 && code <= 599) {
				return false
			}
		}
		if strconv.Itoa(code) == c {
			return true
		}
	}

	return true
}
