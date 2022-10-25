package bitcoin

type NodeInfo struct {
	Blockbook *Blockbook `json:"blockbook"`
	Backend   *Backend   `json:"backend"`
}

type Blockbook struct {
	Version    string `json:"version"`
	BestHeight int64  `json:"bestHeight"`
	InSync     bool   `json:"inSync"`
}

type Backend struct {
	Subversion string `json:"subversion"`
}
