package docker

import (
	"time"
)

type RepoList struct {
	Projects   []Project  `json:"project"`
	Dependency Dependency `json:"dependency"`
}

type Project struct {
	Name       string   `json:"name"`
	Local      string   `json:"local"`
	Repository string   `json:"repository"`
	Branch     string   `json:"branch"`
	Path       string   `json:"path"`
	PreAction  []string `json:"pre-action"`
	PostAction []string `json:"post-action"`
	Depends    []string `json:"depends"`
}

type Container struct {
	ID       string `json:"ID"`
	Name     string `json:"Name"`
	Image    string `json:"Image"`
	Labels   string `json:"Labels"`
	Networks string `json:"Networks"`
	Ports    string `json:"Ports"`
	Project  string `json:"Project"`
	State    string `json:"State"`
}

type Dependency struct {
	Network []Network `json:"network"`
	Volume  []Volume  `json:"volume"`
	Script  []string  `json:"script"`
}

type Volume struct {
	Name    string   `json:"name"`
	Driver  string   `json:"driver"`
	Labels  []string `json:"labels"`
	Options []string `json:"options"`
}

type Network struct {
	Name       string            `json:"Name"`
	Id         string            `json:"Id"`
	Created    time.Time         `json:"Created"`
	Scope      string            `json:"Scope"`
	Driver     string            `json:"Driver"`
	EnableIPv6 bool              `json:"EnableIPv6"`
	IPAM       IPAM              `json:"IPAM"`
	Internal   bool              `json:"Internal"`
	Attachable bool              `json:"Attachable"`
	Ingress    bool              `json:"Ingress"`
	ConfigOnly bool              `json:"ConfigOnly"`
	Options    map[string]string `json:"Options"`
	Labels     map[string]string `json:"Labels"`
}

type IPAM struct {
	Driver  string            `json:"Driver"`
	Options map[string]string `json:"Options"`
	Config  []IPAMConfig      `json:"Config"`
}

type IPAMConfig struct {
	Subnet  string `json:"Subnet"`
	Gateway string `json:"Gateway"`
	IPRange string `json:"IPRange"`
}

type Command struct {
	WorkingDirectory string
	Args             []string
}
