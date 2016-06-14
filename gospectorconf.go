package main

type gospectorConf struct {
	Rules   []rule   `json:"rules`
	Subdirs []string `json:"subdirs"`
	Excluded []string `json:"excluded"`
}

type rule struct {
	Extensions []string `json:"extensions"`
	Words      []string `json:"words"`
}
