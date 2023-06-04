package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// 定义结构体
type Person struct {
	Name  string
	Age   int
	Email string
}

func main() {
	// 定义模板字符串
	tmplStr := "Name: {{.Name}}, Age: {{.Age}}, Email: {{.Email}}"

	// 创建模板对象
	tmpl, err := template.New("person").Parse(tmplStr)
	if err != nil {
		fmt.Println("Error creating template:", err)
		os.Exit(1)
	}

	// 创建结构体实例
	person := Person{
		Name:  "John",
		Age:   30,
		Email: "john@example.com",
	}

	// 渲染模板并保存结果到字符串变量
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, person)
	if err != nil {
		fmt.Println("Error rendering template:", err)
		os.Exit(1)
	}

	// 从Buffer中获取最终的字符串
	result := buf.String()

	// 打印最终的字符串
	fmt.Println(result)
}
