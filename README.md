init-project 项目初始化工具
===

通过配置初始化一个echo框架基础项目

使用方法
---
+ 请求地址 127.0.0.1:22102/init_project
+ 请求体
```json
{
    "project_name":"id_generator",
    "save_path":"/Users/caoshuyu/WorkSpace/GoWork/Test/src/github.com/caoshuyu/",
    "project_path":"github.com/caoshuyu/",
    "open_update_conf":true,
    "mysql":[
        "master",
        "slave"
    ],
    "redis":[
        "master"
    ]
}
```

 生成后使用goland报错处理
 ---
 + 选择GoLand -> Preferences -> Go Modules -> 打勾 Enable Go Modules integration
 + Environment 代理可设置为 GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
 
 运行生成项目
 ---
 + 项目生成后修改对应的.toml配置文件，将数据库，http端口，log存储地址配置成相应值
 + 在项目目录下执行 go mod download 下载依赖
 + 执行 go run main.go 应该可以看到程序正常启动




