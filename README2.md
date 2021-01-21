* [故事开始](#故事开始)
    * [仓库地址](#仓库地址)
    * [更新记录](#更新记录)
    * [联系作者](#联系作者)
    * [详细功能](#详细功能)
    * [使用帮助](#使用帮助)
        * [启动进程](#启动进程)
        * [开机自启](#开机自启)
        * [代理配置](#代理配置)

## 故事开始
在家里经常鼓捣一些小东西，例如:
1. 网盘类seafile,h5ai,nextcloud,Cloudreve
2. 代理类的ssr等
3. 离线下载的: aria2
东西又多有繁杂，没有一个统一的工具来完成自己的需求，刚刚好手上有个树莓派3b+， 打算让此物发挥最后的价值。所以自己做一个简单易用的小工具。

## 仓库地址
> [Github下载](https://github.com/lilacsheep/HomeCenter/releases)

> [Gitee下载](https://gitee.com/Dukeshi/HomeCenter)

## 更新记录
> [更新记录](/doc/update.md)

## 联系作者
QQ: 521287094
Email: lilacsheep@hotmail.com

## 详细功能
1. Http代理功能，基于ssh协议,提供支持多规则分发，后期也会新增其他的协议。
2. 离线下载，目前只支持http协议的文件下载，支持断点续传，多线程下载。
3. 文件管理功能，支持文件预览（暂时支持图片和视频）
4. 进程监控，实时查看程序的CPU，内存，流量等信息。
5. DDNS功能

## 手动构建
详情见build.sh

# 使用帮助
## 启动进程
```bash
./proxy -h 0.0.0.0:80 -path db
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

## 代理配置使用帮助
[Win帮助](https://jingyan.baidu.com/article/72ee561a053a87e16138dfed.html)