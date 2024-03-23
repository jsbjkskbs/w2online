## West2-Online(Golang 下半年综合考核)

### 备注
[如何启动？](./docs/quickstart/guide.md)

### 接口实现

|已实现接口|备注|
|-|-|
|/user/register|无|
|/user/login|无|
|/user/info|双token生成的新access-token会放置在response的Header内|
|/user/avatar/upload|无|
|/auth/mfa/qrcode|无|
|/auth/mfa/bind|无|
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
|/relation/action|无|
|/follower/list|无|
|/following/list|无|
|/friend/list|无|
|ws://127.0.0.1:10000|无|

### 数据同步逻辑

#### 初始化
- MySQL会将所有数据同步到Redis,elasticsearch上

#### 运行时
- 评论、视频、关注一类非高频数据会直接保存到MySQL,同时在Redis(评论、视频)或elasticsearch(视频)保存(该类数据一定有ID)
- 点赞、浏览量一类高频数据不会直接保存到MySQL,刚开始只在Redis保存(该类数据不可能有ID)
- Redis会定期将点赞、浏览等数据同步到MySQL和elastisearch内

#### 其他(如运行时出错)
- 暂无

#### 数据删除逻辑
- 以一天为有效日期,会对未完成的上传请求进行删除操作(删除暂存文件夹,删除Redis上传记录)

### 更新日志
[20240301之前](./docs/update_logs/update_logs_01.md)

### Benchmark
[benchmark.md](./docs/benchmark/benchmark.md)
