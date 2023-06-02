package api

type Config struct {
	ProjectName string   `json:"ProjectName"`
	API         []Module `json:"Api"`
}

type Module struct {
	Module   string     `json:"module"`
	Business []Business `json:"business"`
}

type Business struct {
	Name    string      `json:"name"`
	Operate []Operation `json:"operate"`
}

type Operation struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	Remark string `json:"remark"`
}
