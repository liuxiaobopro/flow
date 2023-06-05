package api

type Operation struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	Rg     string `json:"rg"`
	Remark string `json:"remark"`
}

type Business struct {
	Name    string      `json:"name"`
	Operate []Operation `json:"operate"`
}

type Module struct {
	Module   string     `json:"module"`
	Business []Business `json:"business"`
}

type Config struct {
	ProjectName string `json:"ProjectName"`
	OutPath     struct {
		Handle string `json:"handle"`
	} `json:"OutPath"`
	API []Module `json:"Api"`
}
