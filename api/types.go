package api

type Operation struct {
	Name   string `json:"name"`
	Method string `json:"method"`
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

type RouterGroup struct {
	Name       string `json:"name"`
	Middleware string `json:"middleware"`
}

type Config struct {
	ProjectName string `json:"ProjectName"`
	OutPath     struct {
		Handle     string `json:"handle"`
		Router     string `json:"router"`
		Middleware string `json:"middleware"`
	} `json:"OutPath"`
	API         []Module      `json:"Api"`
	RouterGroup []RouterGroup `json:"RouterGroup"`
}
