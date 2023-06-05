# flow

FlowFrame框架的辅助工具

## 参数说明

```json
{
    "ProjectName": "TestProject", // 项目名称
    "OutPath": { // 生成文件输出路径(必须以/结尾并且是一个单词)
        "handle": "handle/"
    },
    "Api": [ // 接口
        {
            "module": "admin", // 模块
            "business": [ // 业务
                {
                    "name": "user", // 业务名称
                    "operate": [ // 操作
                        {
                            "name": "list", // 操作名称
                            "method": "GET", // 请求方式
                            "remark": "获取用户列表" // 备注
                        },
                        {
                            "name": "add",
                            "method": "POST",
                            "remark": "添加用户"
                        },
                        {
                            "name": "edit",
                            "method": "PUT",
                            "remark": "编辑用户"
                        },
                        {
                            "name": "del",
                            "method": "DELETE",
                            "remark": "删除用户"
                        }
                    ]
                }
            ]
        },
        {
            "module": "app",
            "business": [
                {
                    "name": "order",
                    "operate": [
                        {
                            "name": "list",
                            "method": "GET",
                            "remark": "获取订单列表"
                        },
                        {
                            "name": "add",
                            "method": "POST",
                            "remark": "添加订单"
                        },
                        {
                            "name": "edit",
                            "method": "PUT",
                            "remark": "编辑订单"
                        },
                        {
                            "name": "del",
                            "method": "DELETE",
                            "remark": "删除订单"
                        }
                    ]
                }
            ]
        }
    ]
}

```
