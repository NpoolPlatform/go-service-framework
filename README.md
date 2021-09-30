# Npool go template

[![Test](https://github.com/NpoolPlatform/go-template/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/NpoolPlatform/go-template/actions/workflows/main.yml)

## 目录
* [功能](#功能)
* [命令](#命令)

-----------
### 功能
- [x] 创建sample service
- [x] 封装日志库
- [ ] 集成cli框架(https://github.com/urfave/cli)
- [ ] 集成http server框架(https://github.com/urfave/cli)
- [ ] 集成http client框架(https://github.com/go-resty/resty)
- [ ] 集成consul注册与发现
- [ ] 全局主机环境参数解析
- [ ] 集成apollo配置中心
- [ ] 集成viper配置库(https://github.com/spf13/viper.git)
- [ ] 集成redis访问
- [ ] 集成mysql访问框架(https://github.com/ent/ent)

### 命令
* make init ```初始化仓库，创建go.mod```
* make verify ```验证开发环境与构建环境，检查code conduct```
* make verify-build ```编译目标```
