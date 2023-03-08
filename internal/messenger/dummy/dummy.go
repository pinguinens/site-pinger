package msgsvc

type Dummy struct{}

func New() (*Dummy, error) {
	return &Dummy{}, nil
}

func (m *Dummy) Close() error {
	return nil
}

func (m *Dummy) Send(msg string) error {
	return nil
}

func (m *Dummy) Alarm(code int, method, url, addr string) error {
	return nil
}
