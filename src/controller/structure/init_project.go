package structure

type InitProjectInput struct {
	ProjectName    string   `json:"project_name"`
	SavePath       string   `json:"save_path"`
	Mysql          []string `json:"mysql"`
	Redis          []string `json:"redis"`
	ProjectPath    string   `json:"project_path"`
	OpenUpdateConf bool     `json:"open_update_conf"`
}

type InitProjectOut struct {
}

type BuildModelFileInput struct {
}

type BuildModelFileOut struct {
}
