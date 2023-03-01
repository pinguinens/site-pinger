package site

type HostTable map[string]string

func (t HostTable) Add(domain, ip string) {
	t[domain] = ip
}
