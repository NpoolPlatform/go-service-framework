# Npool go service framework

[![Test](https://github.com/NpoolPlatform/go-service-framework/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/NpoolPlatform/go-service-framework/actions/workflows/main.yml)

## 目录
* [功能](#功能)
* [命令](#命令)
* [步骤](#步骤)
* [最佳实践](#最佳实践)
* [关于mysql](#关于mysql)

-----------
### 功能
- [x] 创建sample service
- [x] 封装日志库
- [x] 统一service cli框架
- [x] 集成cli框架(https://github.com/urfave/cli)
- [x] 集成http server框架(https://github.com/go-chi/chi.git 不需要封装)
- [x] 集成http client框架(https://github.com/go-resty/resty 不需要封装)
- [x] 集成consul注册与发现
- [x] 全局主机环境参数解析
- [x] 集成apollo配置中心(https://github.com/philchia/agollo.git)
- [x] 集成redis访问
- [x] 集成mysql访问框架(https://github.com/ent/ent)
* [x] 集成版本信息
* [x] 集成rabbitmq访问
* [ ] 完善rabbitmq API
* [x] 集成GRPC
* [ ] 添加GRPC proto编译支持

### 命令
* make init ```初始化仓库，创建go.mod```
* make verify ```验证开发环境与构建环境，检查code conduct```
* make verify-build ```编译目标```
* make test ```单元测试```

### 步骤
* 在github上将模板仓库https://github.com/NpoolPlatform/go-service-app-template.git import为https://github.com/NpoolPlatform/my-service-name.git
* git clone https://github.com/NpoolPlatform/my-service-name.git
* cd my-service-name
* mv cmd/service-sample cmd/my-service
* 修改cmd/my-service/main.go中的serviceName为My Service
* 重命名cmd/my-service/ServiceSample.viper.yaml为cmd/my-service/MyService.viper.yaml
* 将cmd/my-service/MyService.viper.yaml中的内容修改为本服务对应内容

### 最佳实践
* 每个服务只提供单一可执行文件，有利于docker镜像打包与k8s部署管理
* 每个服务提供http调试接口，通过curl获取调试信息
* 集群内服务间direct call调用通过服务发现获取目标地址进行调用
* 集群内服务间event call调用通过rabbitmq解耦

### 关于mysql
* 参见https://entgo.io/docs/sql-integration
* 创建app后，从app.Mysql()获取本地mysql client
