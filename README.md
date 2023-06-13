Orca
======

Proxy for [OpenAI API](https://platform.openai.com/docs/api-reference), written by golang

解决两个问题：

1. 缓存重复请求，节省费用
2. 内部使用不需要 token，走流量监控

访问记录会记在 prometheus 和 elastic（可选）

注意目前只能代理 `chat` 和 `completions` 两个接口，其他像 fine-tune 什么的还要走长连接，太麻烦了还没搞

## 为什么要写这个

OpenAI 目前的后台非常惨，一个 org 下好几个帐号，都没法看清每个帐号每天花了多少钱，有一天费用暴增，我都不知道是不是我还是同事消耗的，还是 token 泄漏了，我需要一个更详细、实时的监控。而且我那个程序本来是有缓存的，但是 dev 和 prod 没法共享缓存，所以在想如果走一个统一的代理还能节省每个程序单独写缓存。基于这两点，写了这个代理，目前在公司内部小量使用。

## 部署

在 `server/build` 目录，将 `sample-env.sh` 复制为 `env.sh` 并修改，

```
export OPENAI_API_KEY="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"
export ORCA_WEB=":22035"
export ORCA_ES_ADDR="https://10.0.84.49:9200/"
export ORCA_ES_USER=""
export ORCA_ES_PASS=""
```

`ORCA_WEB` 是代理对外监听地址，不写 IP 表示监听所有

如果不使用 elastic 可以把 `ORCA_ES_USER` 和 `ORCA_ES_PASS` 留空

改完后直接敲 `make` 可运行

注意所有请求会缓存在 `static/cache` 目录

## 使用

声明变量 `OPENAI_API_KEY` 和 `OPENAI_API_BASE`，大部分 OpenAI 脚本应该都会认这两个变量

由于 `OPENAI_API_KEY` 不再是真正的 key（而是上一节 `env.sh` 里的），建议写个自己好懂的，比如 `zhengkai.common`，这样能在统计时做区分

```
export OPENAI_API_KEY="zhengkai.common"
export OPENAI_API_BASE=http://10.0.84.49:22035/v1
```

prometheus endpoint 在服务启动后 `/_metrics` 路径

流量 log 打在 prometheus 和 elastic 后可以自己分析了，这是简单的效果截图

![Grafana](misc/img/grafana.png)

![Kibana](misc/img/kibana.png)

## Docker

`misc/docker` 下有个简单脚本用来构建镜像

先按 **部署** 一节修改 `env.sh`，如果使用 elastic tls 需要将 `http-ca.crt` 替换成你的

`make run` 是尝试本地运行 docker 镜像

`make install` 是拷到目标主机上将镜像作为服务后台运行

具体操作需要看一下 `Makefile` 文件，主要是目录、端口等习惯问题

我运行的方式比较笨拙，只是公司内网跑没考虑太多安全性（比方说 scp env.sh 这个操作），如果有更优雅的方式欢迎指点我，提前感谢

## TODO

* 完善 docker image
* 内置 web 端查询流量/查看每个请求
