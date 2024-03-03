# 启动命令
```bash
docker run -it -p 8080:8080 lipangeng/github-proxy
```

# 使用方式
以本项目部署在本地，即 `127.0.0.1:8080` 为例
## git clone
支持2种模式：
以 https://github.com/lipangeng/github-proxy 为例

可以直接做域名替换:
- http://127.0.0.1:8080/lipangeng/github-proxy

可以在地址前面增加本项目地址即可:
- http://127.0.0.1:8080/github.com/lipangeng/github-proxy
- http://127.0.0.1:8080/https://github.com/lipangeng/github-proxy

## other 
在地址前面增加本项目地址即可，即增加`http://127.0.0.1:8080/`

以分支源码为例: 
- http://127.0.0.1:8080/github.com/auther/project/archive/master.zip
- http://127.0.0.1:8080/https://github.com/auther/project/archive/master.zip