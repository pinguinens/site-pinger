package messenger

import (
	"errors"
	"fmt"
	"net"
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
	message, err := m.conn.Read(rBytes)
	if err != nil {
		return err
	}
	if string(message) != "OK" {
		return errors.New(fmt.Sprintf("message server return: %v", message))
	}

	return nil
}

func (m *Messenger) Alarm(code int, method, url, addr string) error {
	return m.Send(fmt.Sprintf("%v\n%v\n%v\n%v", code, method, url, addr))
}
