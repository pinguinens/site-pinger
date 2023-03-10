package resource

type HostTable map[string]Host

type Host struct {
	Domain string
	Addr   string
}

func (t HostTable) Add(domain, addr string) {
	t[domain] = Host{
		Domain: domain,
		Addr:   addr,
	}
}
