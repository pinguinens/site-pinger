package site

import "net/url"

type Collection struct {
	List  []Site
	hosts []HostTable
}

func (c *Collection) GetHostsList() ([]HostTable, error) {
	if c.hosts != nil {
		return c.hosts, nil
	}

	var err error
	c.hosts, err = c.makeHosts()
	if err != nil {
		return nil, err
	}

	return c.hosts, err
}

func (c *Collection) makeHosts() ([]HostTable, error) {
	var maxHosts int
	for _, s := range c.List {
		if len(s.Target.Hosts) > maxHosts {
			maxHosts = len(s.Target.Hosts)
		}
	}

	hostTables := make([]HostTable, maxHosts, maxHosts)
	for i, _ := range hostTables {
		hostTables[i] = HostTable{}
	}
	for _, s := range c.List {
		uri, err := url.Parse(s.Target.URI)
		if err != nil {
			return nil, err
		}
		for _, h := range s.Target.Hosts {
			for _, ht := range hostTables {
				if _, ok := ht[uri.Host]; !ok {
					ht.Add(uri.Host, h)
					break
				}
			}
		}
	}

	return hostTables, nil
}
