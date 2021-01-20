## HomeProxy
HTTP/HTTPS proxy over SSH. 在家庭使用的功能中心

### QQ群
782391570

## Future
- [x] 支持控制台需要账号访问
- [x] 服务与主机监控信息展示
- [x] 自定义dns
- [ ] IP访问控制，访问黑名单
- [x] 新增：DDNS
- [x] 代理认证
- [ ] SS实例的添加（被封的概率太大，暂时不添加）

## 更新预告
1. DDNS (aliDNS)
> [更新记录](/doc/update.md)

## 手动构建
详情见build.sh

# 使用帮助
## 下载
根据自己的版本下载相应文件[下载](https://github.com/lilacsheep/HomeCenter/releases)

## 启动进程
```bash
./proxy -h 0.0.0.0:80 -path db -name default
```
## Linux（树莓派） 开机自启 （Centos7以上版本）
修改systemd/system/proxy.service 中 /path/to/proxy 为你部署的具体路径

```bash
cp systemd/system/proxy.service /lib/systemd/system/
systemctl enable proxy
systemct start proxy
```
## 访问
访问http://your_address:port

## 默认账号
> username: admin
> password: !QAZ2wsx

## 代理
![](/pic/start.jpg)

## 文件管理
![](/pic/file.jpg)

## 应用监控
![](/pic/monitor.jpg)

## ddns
![](/pic/ddns.jpg)

## 配置http代理
[Win帮助](https://jingyan.baidu.com/article/72ee561a053a87e16138dfed.html)