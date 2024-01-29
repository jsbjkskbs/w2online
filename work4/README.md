## Uncomplete

#### ElasticSearch(7.17.16)采用的插件有：
1. elasticsearch-analysis-ik-7.17.16
#### 请下载并提取到./pkg/configs/es/plugins内
#### 或在docker-compose.yaml中为elastic节点添加指令
    ./bin/plugin-install <plugin name or url>

### 2024.1.29 (1st upload)
- 采用hertz框架
- 完成用户模块
- 完成视频模块(其中视频上传为分片上传，被拆分为四个接口)
- 引入Mysql+Redis+ElasticSearch+OSS(qiniuyun)
