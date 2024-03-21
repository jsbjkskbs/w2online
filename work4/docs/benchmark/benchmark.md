## Benchmark

### 注：测试样本(传入的Request以及数据库数据量)不够大

|已测试接口|结果|备注|
|-|-|-|
|/user/register|[user.register.txt](./result/user.register.txt)|
|/user/login|[user.login.nomfa.txt](./result/user.login.nomfa.txt)|用户无MFA|
|/user/login|[user.login.mfa.txt](./result/user.login.mfa.txt)|用户MFA验证错误|
|/user/info|[user.info.txt](./result/user.info.txt)||
|/video/feed|[video.feed.txt](./result/video.feed.txt)||
|/video/list|[video.list.txt](./result/video.list.txt)||
|/video/popular|[video.popular.txt](./result/video.popular.txt)||
|/video/search|[video.search.txt](./result/video.search.txt)||
|/action/like|[action.like.txt](./result/action.lile.txt)||
|/like/list|[like.list.txt](./result/like.list.txt)||
|/comment/list|[comment.list.txt](./result/comment.list.txt)||
|/following/list|[following.list.txt](./result/following.list.txt)||
|/friend/list|[friend.list.txt](./result/friend.list.txt)||