package main

type gospectorConf struct {
	Rules   []rule   `json:"rules`
	Subdirs []string `json:"subdirs"`
}

type rule struct {
	Extensions []string `json:"extensions"`
	Words      []string `json:"words"`
}
