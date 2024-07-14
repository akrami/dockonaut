package docker

import (
	"time"
)

type Config struct {
	Projects   []Project  `json:"Projects"`
	Dependency Dependency `json:"Dependency"`
}

type Project struct {
	Name       string   `json:"Name"`
	Local      string   `json:"Local"`
	Repository string   `json:"Repository"`
	Branch     string   `json:"Branch"`
	Path       string   `json:"Path"`
	PreAction  []string `json:"PreAction"`
	PostAction []string `json:"PostAction"`
	Depends    []string `json:"Depends"`
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
	Networks []Network `json:"Networks"`
	Volumes  []Volume  `json:"Volumes"`
	Scripts  []string  `json:"Scripts"`
}

type Volume struct {
	Name       string            `json:"Name"`
	Driver     string            `json:"Driver"`
	Labels     map[string]string `json:"Labels"`
	Mountpoint string            `json:"Mountpoint"`
	Options    map[string]string `json:"Options"`
	Scope      string            `json:"Scope"`
	CreatedAt  time.Time         `json:"CreatedAt"`
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
