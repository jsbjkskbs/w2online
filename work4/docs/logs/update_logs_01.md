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

#### 2024.2.1 (4th upload)
- 完成WebSocket聊天

#### 2024.2.1 (5th upload)
- 完成所提供的全部接口
- 基本改良或统一各个package返回的error格式

#### 2024.2.1 (6th upload)
- 添加项目树文件(./tree.txt)