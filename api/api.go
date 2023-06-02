package api

import (
	"fmt"
	"os"

	stringx "github.com/liuxiaobopro/gobox/string"
)

var (
	conf       *Config
	handlePath = "handle/"
)

type fileInfo struct {
	folderPath   string // 文件夹路径
	packname     string // 包名
	baseFilename string // 基础文件名
	daoFilename  string // dao文件名

	baseStructName  string // 基础结构体名
	reqStructName   string // 请求结构体名
	replyStructName string // 响应结构体名
}

func Run(config *Config) {
	conf = config
	for _, v := range conf.API {
		for _, v1 := range v.Business {
			for _, v2 := range v1.Operate {
				folderPath := fmt.Sprintf("%s%s/%s/%s", handlePath, v.Module, v1.Name, v2.Name)
				// 判断文件夹是否存在
				if _, err := os.Stat(folderPath); os.IsNotExist(err) {
					// 创建文件夹
					err := os.MkdirAll(folderPath, os.ModePerm)
					if err != nil {
						fmt.Println("Create folder failed.")
						os.Exit(1)
					}
				}

				firstUpOperateName := stringx.FirstUp(v2.Name)
				firstUpBusinessName := stringx.FirstUp(v1.Name)

				genFileInfo := &fileInfo{
					folderPath:   folderPath,
					packname:     v2.Name,
					baseFilename: v2.Name + ".go",
					daoFilename:  v2.Name + "_dao.go",

					baseStructName:  firstUpBusinessName + firstUpOperateName,
					reqStructName:   v1.Name + firstUpOperateName + "Req",
					replyStructName: v1.Name + firstUpOperateName + "Reply",
				}
				genFile(genFileInfo)
			}
		}
	}
}

func genFile(f *fileInfo) {
	// 判断文件是否存在
	filePath := fmt.Sprintf("%s/%s", f.folderPath, f.baseFilename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 创建文件
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Create file failed.")
			os.Exit(1)
		}
		defer file.Close()

		// 写入文件
		_, err = file.WriteString(baseFileContent(f))
		if err != nil {
			fmt.Println("Write file failed.")
			os.Exit(1)
		}
	} else {
		fmt.Println("File already exists.")
		os.Exit(1)
	}

	// 判断文件是否存在
	filePath = fmt.Sprintf("%s/%s", f.folderPath, f.daoFilename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 创建文件
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Create file failed.")
			os.Exit(1)
		}
		defer file.Close()

		// 写入文件
		_, err = file.WriteString(fmt.Sprintf("package %s\n\n", f.packname))
		if err != nil {
			fmt.Println("Write file failed.")
			os.Exit(1)
		}
	} else {
		fmt.Println("File already exists.")
		os.Exit(1)
	}
}

func baseFileContent(f *fileInfo) string {
	return `package ` + f.packname + `

import (
	"github.com/liuxiaobopro/gobox/gin/ctx"
	replyx "github.com/liuxiaobopro/gobox/reply"
)

type ` + f.baseStructName + `Flow struct {
	ctx.Flow
}

type ` + f.reqStructName + `Req struct {
}

type ` + f.replyStructName + `Reply struct {
}

func (d *` + f.baseStructName + `Flow) FlowHandle() {
	var req ` + f.reqStructName + `Req
	if err := d.ShouldBind(&req); err != nil {
		d.ReturnJson(replyx.ParamErrT)
		return
	}

	d.SetReq(&req)
}

func (d *` + f.baseStructName + `Flow) FlowValidate() {}

func (d *` + f.baseStructName + `Flow) FlowLogic() {
	req, ok := d.GetReq().(*` + f.reqStructName + `Req)
	if !ok {
		d.ReturnJson(replyx.InternalErrT)
		return
	}
	var (
		out = &` + f.replyStructName + `Reply{}
	)

	_ = req

	d.ReturnSucc(out)
}`
}
