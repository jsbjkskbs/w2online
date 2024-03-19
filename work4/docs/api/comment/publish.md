## /comment/publish
发布评论时，评论先通过生产者入队作为缓冲；消费者订阅队列，当队列不为空时，获得消费信息，将消费信息转化进入数据库。

![Alt text](publish.png)