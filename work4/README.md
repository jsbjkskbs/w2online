## West2-Online(Golang 下半年综合考核)

### 备注

#### ElasticSearch(7.17.16)采用的插件有：
1. elasticsearch-analysis-ik-7.17.16
#### 请下载并提取到./pkg/configs/es/plugins内
#### 或在docker-compose.yaml中为elastic节点添加指令
    ./bin/plugin-install <plugin name or url>

### 接口实现

|已实现接口|备注|
|-|-|
|/user/register|无|
|/user/login|无|
|/user/info|双token生成的新access-token会放置在response的Header内|
|/user/avatar/upload|用redis缓存上传状态，防止上传过程中仍执行上传操作|
|/video/publish|被拆分为以下四个接口|
|/video/publish/start|请求新建上传事务,服务端用redis记录上传请求|
|/video/publish/uploading|上传分片,客户端需将视频分片为m3u8和ts文件,服务端用redis记录上传过程并检验md5|
|/video/publish/cancle|取消上传,服务端执行取消动作|
|/video/publish/complete|上传结束,服务端比对上传操作及记录,无误后上传至OSS|
|/video/search|用elasticsearch实现|
|/video/feed|用elasticsearch实现|
|/video/visit/:id|测试用接口|
|/video/popular|用redis实现|
|/comment/publish|可对评论进行评论,但所有子评论都只有一个母评论,即不会出现子评论的子评论|
|/like/action|无|
|/like/list|无|
|/comment/list|无|
|/comment/delete|无|
|relation/action|无|
|follower/list|无|
|following/list|无|
|friend/list|无|

### 数据同步逻辑

#### 初始化
- MySQL会将所有数据同步到Redis,elasticsearch上

#### 运行时
- 评论、视频、关注一类非高频数据会直接保存到MySQL,同时在Redis(评论、视频)或elasticsearch(视频)保存(该类数据一定有ID)
- 点赞、浏览量一类高频数据不会直接保存到MySQL,刚开始只在Redis保存(该类数据不可能有ID)
- Redis会定期将点赞、浏览等数据同步到MySQL和elastisearch内

#### 其他(如运行时出错)
- 暂无

### 数据删除逻辑
- 以一天为有效日期,会对未完成的上传请求进行删除操作(删除暂存文件夹,删除Redis上传记录)

### 更新日志

#### 2024.1.29 (1st upload)
- 采用hertz框架
- 完成用户模块
- 完成视频模块(其中视频上传为分片上传，被拆分为四个接口)
- 引入Mysql+Redis+ElasticSearch+OSS(qiniuyun)

#### 2024.1.30 (2nd upload)
- 完成互动模块

#### 2024.1.31 (3rd upload)
- 完成社交模块(至此完成全部寒假要求的接口,仅剩余部分额外要求的接口)
- 修改了docker-compose中的连接逻辑(work由host模式变为bridge模式,全部容器由一个共同net连接)
- 由于社交中关注与粉丝关系在Redis内的数据结构为Set,无法正常完成分页功能,故访问数据库以获取列表。而好友关系由Redis中Set的交集获取，可能会产生无序的问题。