package types

type Inventory struct {
	Hosts  []*Host
	Groups []*Group
	Vars   map[string]string
}

type Host struct {
	Ip          string            `json:"ip,omitempty"`
	Name        string            `json:"name,omitempty"`
	Port        int32             `json:"port,omitempty"`
	User        string            `json:"user,omitempty"`
	Password    string            `json:"password,omitempty"`
	PrivateKey  string            `json:"privateKey,omitempty"`
	ProxyConfig *ProxyConfig      `json:"proxyConfig,omitempty"`
	Vars        map[string]string `json:"vars,omitempty"`
}

type ProxyConfig struct {
	Enable   bool   `json:"enable,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Ip       string `json:"ip,omitempty"`
	Port     int32  `json:"port,omitempty"`
}

type Group struct {
	Name     string            `json:"name,omitempty"`
	Hosts    []string          `json:"hosts,omitempty"`
	Children []string          `json:"children,omitempty"`
	Vars     map[string]string `json:"vars,omitempty"`
}
