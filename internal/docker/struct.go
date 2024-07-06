package docker

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
	Name       string   `json:"name"`
	Driver     string   `json:"driver"`
	Labels     []string `json:"labels"`
	Options    []string `json:"options"`
	Internal   bool     `json:"internal"`
	Attachable bool     `json:"attachable"`
	IPV6       bool     `json:"ipv6"`
	SubNet     string   `json:"subnet"`
	IPRange    string   `json:"ip-range"`
	Gateway    string   `json:"gateway"`
}
