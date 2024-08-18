# 终端日志查看器
目前支持 [xray-core](https://github.com/XTLS/Xray-core) 的 `access_log` 和 [nginx](https://github.com/nginx/nginx) 默认格式的 `access_log`

### 使用方法：
1. 下载qqwry.dat
```bash
# 中国大陆可能需要设置终端代理
# export HTTPS_PROXY=http://127.0.0.1:8080
# export HTTP_PROXY=http://127.0.0.1:8080
go run main.go download
```

2. 统计日志数据
> 注意：`--search` 非必要参数，若指定匹配内容则禁止使用小括号来提取字符串
```bash
# 2.1 统计xray-core数据
go run main.go statistics --kind xray-access --logs access.log --search "[^\.]*\.google\.com"

# 2.2 统计nginx数据
go run main.go statistics --kind nginx --logs access.log --search "/favicon.png"
```

### 感谢以下项目（排名不分先后）:
* [qqwry.dat](https://github.com/FW27623/qqwry)
* [GeoCN](https://github.com/ljxi/GeoCN)
* [maxminddb-golang](https://github.com/oschwald/maxminddb-golang)
* [GeoLite](https://github.com/P3TERX/GeoLite.mmdb)
