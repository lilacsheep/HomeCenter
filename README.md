* [Story Start] (#Story Start)
    * [Warehouse address](#Warehouse address)
    * [Update Record](#Update Record)
    * [Contact author](#Contact author)
    * [Detailed function](#Detailed function)
    * [Use help](#use help)
        * [Start process](#Start process)
        * [Start-up self-start] (# Start-up self-start)
        * [Proxy configuration](#Proxy configuration)

## 中文
[中文文档](/README2.md)

## The story begins
Often tinker with small things at home, such as:
1. SkyDrive Seafile, h5ai, nextcloud, Cloudreve
2. ssr of proxy class
3. Offline download: aria2
There are so many complicated things, and there is no unified tool to complete my needs. I just have a Raspberry Pi 3b+, and I intend to let this thing play its final value. 


## Warehouse Address
> [Github download](https://github.com/lilacsheep/HomeCenter/releases)

> [Gitee download](https://gitee.com/Dukeshi/HomeCenter)

## update record
> [Update Record](/doc/update.md)

## Contact the author
QQ: 521287094
Email: lilacsheep@hotmail.com

## Detailed
1. Http proxy function, based on the ssh protocol, provides support for multi-rule distribution, and other protocols will be added later.
2. Offline download, currently only supports file download of http protocol, supports resumable upload, multi-threaded download.
3. File management function, support file preview (for the time being supports pictures and videos)
4. Process monitoring, real-time view of the program's CPU, memory, traffic and other information.
5. DDNS function

## Manually build
See build.sh for details

# Using help
## Start the process
```bash
./proxy -h 0.0.0.0:80 -path db
```
## Linux (Raspberry Pi) Boot from boot (Centos7 and above)
Modify the specific path that /path/to/proxy in systemd/system/proxy.service is deployed for you

```bash
cp systemd/system/proxy.service /lib/systemd/system/
systemctl enable proxy
systemct start proxy
```

## Visit
Visit http://your_address:port

## Default account
> username: admin
> password: !QAZ2wsx

## Proxy configuration help
[Win Help](https://jingyan.baidu.com/article/72ee561a053a87e16138dfed.html)