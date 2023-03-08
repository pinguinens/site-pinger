package messenger

type Messenger interface {
	Close() error
	Send(msg string) error
	Alarm(code int, method string, url string, addr string) error
}
