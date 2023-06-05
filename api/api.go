package api

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	stringx "github.com/liuxiaobopro/gobox/string"
)

var (
	conf       *Config
	handlePath string
	routerPath string
)

type (
	handleFileInfo struct {
		FolderPath   string // 文件夹路径
		Packname     string // 包名
		BaseFilename string // 基础文件名
		DaoFilename  string // dao文件名

		BaseStructName  string // 基础结构体名
		ReqStructName   string // 请求结构体名
		ReplyStructName string // 响应结构体名
	}
)

func initConfig() {
	handlePath = conf.OutPath.Handle
}

func Run(config *Config) {
	conf = config
	initConfig()
	format()
}

// format 格式化
func format() {
	for _, v := range conf.API {
		for _, v1 := range v.Business {
			for _, v2 := range v1.Operate {
				FolderPath := fmt.Sprintf("%s%s/%s/%s", handlePath, v.Module, v1.Name, v2.Name)
				// 判断文件夹是否存在
				if _, err := os.Stat(FolderPath); os.IsNotExist(err) {
					// 创建文件夹
					err := os.MkdirAll(FolderPath, os.ModePerm)
					if err != nil {
						fmt.Println("Create folder failed.")
						os.Exit(1)
					}
				}

				firstUpOperateName := stringx.FirstUp(v2.Name)
				firstUpBusinessName := stringx.FirstUp(v1.Name)

				handle := &handleFileInfo{
					FolderPath:   FolderPath,
					Packname:     v2.Name,
					BaseFilename: v2.Name + ".go",
					DaoFilename:  "dao.go",

					BaseStructName:  firstUpBusinessName + firstUpOperateName,
					ReqStructName:   v1.Name + firstUpOperateName + "Req",
					ReplyStructName: v1.Name + firstUpOperateName + "Reply",
				}

				genHandleFile(handle)
			}
		}
	}
}

// genHandleFile 生成handle文件
func genHandleFile(f *handleFileInfo) {
	// 判断文件是否存在
	filePath := fmt.Sprintf("%s/%s", f.FolderPath, f.BaseFilename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 创建文件
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Create file failed.")
			os.Exit(1)
		}
		defer file.Close()

		// 写入文件
		// 创建模板对象
		tmpl, err := template.New("baseFileHandle").Parse(baseFileHandleContent())
		if err != nil {
			fmt.Println("Error creating template:", err)
			os.Exit(1)
		}

		// 渲染模板并保存结果到字符串变量
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, f)
		if err != nil {
			fmt.Println("Error rendering template:", err)
			os.Exit(1)
		}

		_, err = file.WriteString(buf.String())
		if err != nil {
			fmt.Println("Write file failed.")
			os.Exit(1)
		}
	}

	// 判断文件是否存在
	filePath = fmt.Sprintf("%s/%s", f.FolderPath, f.DaoFilename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 创建文件
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Create file failed.")
			os.Exit(1)
		}
		defer file.Close()

		// 写入文件
		_, err = file.WriteString(fmt.Sprintf("package %s\n\n", f.Packname))
		if err != nil {
			fmt.Println("Write file failed.")
			os.Exit(1)
		}
	}
}

// baseFileHandleContent 基础文件内容
func baseFileHandleContent() string {
	return `package {{.Packname}}

import (
	"github.com/liuxiaobopro/gobox/gin/ctx"
	replyx "github.com/liuxiaobopro/gobox/reply"
)

type {{.BaseStructName }}Flow struct {
	ctx.Flow
}

type {{.ReqStructName }}Req struct {
}

type {{.ReplyStructName }}Reply struct {
}

func (d *{{.BaseStructName }}Flow) FlowHandle() *replyx.T {
	var req {{.ReqStructName }}Req
	if err := d.ShouldBind(&req); err != nil {
		return replyx.ParamErrT
	}

	d.SetReq(&req)
	return nil
}

func (d *{{.BaseStructName }}Flow) FlowValidate() *replyx.T {
	return nil
}

func (d *{{.BaseStructName }}Flow) FlowLogic() *replyx.T {
	req, ok := d.GetReq().(*{{.ReqStructName }}Req)
	if !ok {
		return replyx.InternalErrT
	}
	var (
		out = &{{.ReplyStructName }}Reply{}
	)

	_ = req

	d.ReturnSucc(out)
	return nil
}`
}
