## /like/action

点赞操作属于高频操作，故先将记录保存至Redis。服务端内实现周期性同步，每隔一段时间将Redis内的信息同步到MySQL和ElasticSearch


![Alt text](likeaction.png)