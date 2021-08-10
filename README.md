# mengjiang

## 安装

```shell
make install
```

## 使用

```shell
# 直接启动代理默认监听http的8080和socks的1086端口
mengjiang

# 指定端口
mengjiang --http :8080 --socks :1086

# 透明代理 只支持linux
mengjiang -mod nat
```

## 透明代理

透明代理需要使用iptables的nat表将流量导入到代理的8080端口
