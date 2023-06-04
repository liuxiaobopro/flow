package api

import (
	"bytes"
	"fmt"
	"os"
	"strings"
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

	routerFileInfo struct {
		ProjectName string            // 项目名
		Packname    string            // 包名
		GroupName   string            // 组名
		FuncName    string            // 函数名
		Routers     []*routerTempInfo // 路由
	}

	routerTempInfo struct {
		ProjectName string // 项目名
		Method      string // 方法
		Path        string // 路径
		Router      string // 路由
		HandlerName string // 处理函数
		Remark      string // 备注
	}

	routerVars struct {
		GroupName string            // 名称
		routers   []*routerTempInfo // 路由
	}
)

func initConfig() {
	handlePath = conf.OutPath.Handle
	routerPath = conf.OutPath.Router
}

func Run(config *Config) {
	conf = config
	initConfig()
	format()
}

// format 格式化
func format() {
	var routerGroup = make(map[string]*routerVars)

	// 判断router文件夹是否存在
	if _, err := os.Stat(routerPath); os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(routerPath, os.ModePerm)
		if err != nil {
			fmt.Println("Create folder failed.")
			os.Exit(1)
		}
	}

	if len(conf.RouterGroup) > 0 {
		for _, v := range conf.RouterGroup {
			routerGroup[v.Name] = &routerVars{
				GroupName: v.Name,
				routers:   []*routerTempInfo{},
			}
		}
	}

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

				//#region handle
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
				//#endregion

				//#region router
				if v2.Rg != "" {
					routerDetail, ok := routerGroup[v2.Rg]
					if !ok {
						// 路由不存在, 已经跳过
						fmt.Println("router not exist, skip")
					} else {
						routerDetail.routers = append(routerDetail.routers, &routerTempInfo{
							ProjectName: conf.ProjectName,
							Method:      v2.Method,
							Path:        fmt.Sprintf("%s%s/%s/%s", conf.OutPath.Handle, v.Module, v1.Name, v2.Name),
							Router:      fmt.Sprintf("/%s/%s", v1.Name, v2.Name),
							HandlerName: handle.Packname + "." + handle.BaseStructName,
							Remark:      v2.Remark,
						})
					}
				}
				//#endregion
			}
		}
	}

	if len(conf.RouterGroup) > 0 {
		for _, v := range routerGroup {
			if len(v.routers) == 0 {
				continue
			}
			genRouterFile(&routerFileInfo{
				ProjectName: conf.ProjectName,
				Packname:    strings.TrimRight(conf.OutPath.Router, "/"),
				GroupName:   v.GroupName,
				FuncName:    stringx.FirstUp(v.GroupName),
				Routers:     v.routers,
			})
		}
	}
}

// genRouterFile 生成router文件
func genRouterFile(r *routerFileInfo) {
	// 判断文件是否存在
	filePath := fmt.Sprintf("%s/%s.go", routerPath, r.GroupName)

	var err error
	if _, err = os.Stat(filePath); err == nil {
		// 文件存在, 删除文件
		err1 := os.Remove(filePath)
		if err1 != nil {
			fmt.Println("Remove file failed.")
			os.Exit(1)
		}
	} else if !os.IsNotExist(err) {
		fmt.Println("Stat file failed.")
		os.Exit(1)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Create file failed.")
		os.Exit(1)
	}
	defer file.Close()

	// 写入文件
	// 创建模板对象
	tmpl, err := template.New("routerFile").Parse(routerFileContent())
	if err != nil {
		fmt.Println("Error creating template:", err)
		os.Exit(1)
	}

	// 渲染模板并保存结果到字符串变量
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, r)
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

// routerFileContent router文件内容
func routerFileContent() string {
	return `package {{.Packname}}

import (
	"github.com/gin-gonic/gin"
	{{range .Routers}}
	"{{.ProjectName}}/{{.Path}}" {{end}}

	"github.com/liuxiaobopro/gobox/gin/ctx"
)

func Add{{.FuncName}}(rg *gin.RouterGroup) {
	{{range .Routers}}
	rg.Handle("{{.Method}}", "{{.Router}}", ctx.Use(new({{.HandlerName}}Flow))) // {{.Remark}} {{end}}
}`
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

func (d *{{.BaseStructName }}Flow) FlowHandle() {
	var req {{.ReqStructName }}Req
	if err := d.ShouldBind(&req); err != nil {
		d.ReturnJson(replyx.ParamErrT)
		return
	}

	d.SetReq(&req)
}

func (d *{{.BaseStructName }}Flow) FlowValidate() {}

func (d *{{.BaseStructName }}Flow) FlowLogic() {
	req, ok := d.GetReq().(*{{.ReqStructName }}Req)
	if !ok {
		d.ReturnJson(replyx.InternalErrT)
		return
	}
	var (
		out = &{{.ReplyStructName }}Reply{}
	)

	_ = req

	d.ReturnSucc(out)
}`
}
