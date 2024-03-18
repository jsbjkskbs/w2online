## Requirements

### A.Docker
程序推荐在Docker环境下运行(当然可以不用)，具体安装方式见[Docker官网](https://www.docker.com/)

### B.ElasticSearch插件(7.17.16)
#### 需要插件:
##### 1. elasticsearch-analysis-ik-7.17.16(用于中文分词)
请下载并提取到./pkg/configs/es/plugins内或在docker-compose.yaml中为elastic节点添加指令

    ./bin/plugin-install <plugin name or url>

## Config

### A.MySQL
1. MySQL初始表(entry-point)-[mysql.ini](../../pkg/configs/sql/init.sql)
2. MySQL配置-[docker-compose-env.env](../../docker-compose-env.env)&[docker-compose.yaml](../../docker-compose.yaml)

### B.Redis
1. Redis用户配置-[docker-compose-env.env](../../docker-compose-env.env)&[docker-compose.yaml](../../docker-compose.yaml)

### C.ElasticSearch&Kibana
1. ElasticSearch&Kibana配置-[docker-compose-env.env](/docker-compose-env.env)&[docker-compose.yaml](../../docker-compose.yaml)

### D.OSS(七牛云)
1. OSS秘钥、Bucket、域名-[config.json](../../config.json)

### E.Sentinel
1. QPS配置-[config.json](../../config.json)

## Launch

### Linux

于根目录输入指令

    bash docker_build.sh

完成镜像构建，可对[docker_build.sh](../../docker_build.sh)或者[Dockerfile](../../docker-build/Dockerfile)进行调整

随后输入

    bash docker_compose_release.sh

在Docker内启动服务端以及所需的中间件，不过由于中间件启动不能保证比服务端程序快，故服务端程序需要二次启动(其实可以在[docker-compose.yaml](../../docker-compose.yaml)添加个(container name).restart.always属性即可)

## 备注
由于[docker-compose.yaml](../../docker-compose.yaml)对所有容器应用桥接模式,共用一个虚拟网络,故go的连接地址应当是带有容器名并为容器内端口的一段地址,如`redis:6379`(`container-name = redis`   `port = 16379:6379`)

