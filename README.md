## HomeProxy
HTTP/HTTPS proxy over SSH.

### QQ群
782391570

## Future
- [x] 支持树莓派3b+, 4b+
- [ ] 支持域名定向转发
- [ ] 支持替换轮训器,现在随机
- [ ] 支持http Basic auth
- [ ] 支持https 证书
- [ ] 支持控制台需要账号访问

## 更新 2020年10月9日
1. 去除sqlite3依赖 
2. 去除其他文件或目录依赖，改成命令行模式

## 手动构建
详情见build.sh

# 使用帮助
## 下载
根据自己的版本下载相应文件[下载](https://gitee.com/Dukeshi/home-proxy/releases)

## 启动进程
```bash
./proxy -h 0.0.0.0:80 -path db -name default
```

## 启动代理
![](/pic/start.jpg)

## 新增规则
![](/pic/add_role.jpg)

## 新增转发代理
![](/pic/add_instance.jpg)

## 配置http代理
[帮助](https://jingyan.baidu.com/article/72ee561a053a87e16138dfed.html)