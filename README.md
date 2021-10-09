# Npool go service framework

[![Test](https://github.com/NpoolPlatform/go-service-framework/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/NpoolPlatform/go-service-framework/actions/workflows/main.yml)

## 目录
* [功能](#功能)
* [命令](#命令)
* [步骤](#步骤)

-----------
### 功能
- [x] 创建sample service
- [x] 封装日志库
- [x] 统一service cli框架
- [x] 集成cli框架(https://github.com/urfave/cli)
- [ ] 集成http server框架(https://github.com/go-chi/chi.git)
- [ ] 集成http client框架(https://github.com/go-resty/resty)
- [x] 集成consul注册与发现
- [x] 全局主机环境参数解析
- [x] 集成apollo配置中心(https://github.com/philchia/agollo.git)
- [ ] 集成redis访问
- [ ] 集成mysql访问框架(https://github.com/ent/ent)
* [x] 集成版本信息

### 命令
* make init ```初始化仓库，创建go.mod```
* make verify ```验证开发环境与构建环境，检查code conduct```
* make verify-build ```编译目标```

### 步骤
* 在github上将模板仓库https://github.com/NpoolPlatform/go-template.git import为https://github.com/NpoolPlatform/my-service-name.git
* git clone https://github.com/NpoolPlatform/my-service-name.git
* cd my-service-name
* mkdir cmd/my-service
* curl https://raw.githubusercontent.com/NpoolPlatform/go-service-framework/master/cmd/service-sample/main.go -o cmd/my-service/main.go
