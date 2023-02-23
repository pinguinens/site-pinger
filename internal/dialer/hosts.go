package dialer

type HostTable map[string]Host

type Host struct {
	Domain string
	Addr   string
}

func (t *HostTable) Add(domain, addr string) {}
