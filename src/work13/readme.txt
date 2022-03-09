【毕业项目】

对当下自己项目中的业务，进行一个微服务改造，需要考虑如下技术点：

微服务架构（BFF、Service、Admin、Job、Task 分模块）
API 设计（包括 API 定义、错误码规范、Error 的使用）
gRPC 的使用
Go 项目工程化（项目结构、DI、代码分层、ORM 框架）
并发的使用（errgroup 的并行链路请求）
微服务中间件的使用（ELK、Opentracing、Prometheus、Kafka）
缓存的使用优化（一致性处理、Pipeline 优化）

解答：

一.背景：
1.平时工作中主要以PHP开发为主，涉及不到Go项目开发。同时，平时工作也有点忙，工作中用的技术栈也有限。
2.因此，本次毕业项目主要以学习相关工具为主。


二.毕业项目：
1.环境搭建（基于腾讯云服务器）：
a.docker+logstash+kibana+es+mysql环境搭建，logstash拉去mysql表数据，放到es，并在kibana进行展示数据
b.docker+zookeeper+kafka环境搭建
c.docker+redis（单机集群）搭建

2.golang操作各类中间件使用案列：见同目录下/lastproject/cmd/middlewaretest/目录下对应名称的demo文件
1) es基本操作：添加，查询，批量查询，更新，条件查询更新删除文档等；
2）mysql基本操作：增删改查事务等操作；
3) redis集群基本操作：
a.基本键值，list，set，hash，sorted set等各类数据结构和pipeline事务处理的使用；
b.pipeline优化：
    (1)目的是将一批命令打包到一个内部维护的queue里，然后建立socket与server交互，这时只会发送一次命令，也就是只会交互一次;
    (2)然后queue内的命令都执行完后会一起返回结果，这样大大减少了通信的次数，降低了通信所耗费的时间;
    (3)queue是先进先出，所以可以保证执行顺序;
    (4)server对pipeline的命令结果进行缓存处理，会消耗很多内存，同时client执行命令后，结果会缓存在client-revice-buffer中;
    (5)如果缓存满了，通知server停止发送数据，因此要控制好每次pipeline的大小，保持效率最高。
    (6)适用场景：对实时性要求不高；批量将数据写入 redis，允许一定比例的写入失败。
    (7)不适用场景：要求可靠性高，每次都需要实时知道这次操作是否成功，数据是否写入redis了等对实时性的这种需求都不适合。
4）kafka基本操作：简单封装生产者和消费者的demo实现，以及一些常见问题的处理：
a.producer把消息发送给broker时产生的丢失：网络抖动；master接收到消息，在未将消息同步给follower之前，挂掉了；master接收到消息，master未成功将消息同步给每个follower。
    解决：producer设置acks参数，config.Producer.RequiredAcks = sarama.WaitForAll
b.某个broker消息尚未从内存缓冲区持久化到磁盘，就挂掉了。
    解决：设置参数，加快消息持久化的频率，能在一定程度上减少这种情况发生的概率。但提高频率自然也会影响性能。
c.consumer成功拉取到了消息，consumer挂了。
    解决：调整相关业务逻辑（自动提交），或者设置手动sync，消费成功才提交。复杂场景，可以考虑二阶段提交。
d.订阅Kafka的消费者按照消息顺序写入mysql，而不是随机写入？
    解决：初始化的时候，设置选择分区的策略为Hash：config.Producer.Partitioner = sarama.NewHashPartitioner；生成消息前，设置消息的key值，如下：
    msg := &sarama.ProducerMessage{
        Topic: "testAutoSyncOffset",
        Value: sarama.StringEncoder("hello"),
        Key: sarama.StringEncoder(strconv.Itoa(RecvID)),
    }
e.多线程情况下一个partition的乱序处理.
    解决：可以通过写 N 个内存 queue，具有相同 key 的数据都到同一个内存 queue；然后对于 N 个线程，每个线程分别消费一个内存 queue 即可，这样就能保证顺序性。PS：就像4 % 10 = 4，14 % 10 = 4，他们取余都是等于4，所以落到了一个partition，但是key值不一样啊，我们可以自己再取余，放到不同的queue里面。
f.重复消费和消息幂等问题。
    解决：
    （1）如果是存在redis中不需要持久化的数据，比如string类型，set具有天然的幂等性，无需处理。
    （2）插入mysql之前，进行一次query操作，针对每个客户端发的消息，生成一个唯一的ID（雪花算法），或者直接把消息的ID设置为唯一索引。

3.基于grpc和protobuf的客户端和服务端的简单实时通讯demo实现：
1）见同目录下/lastproject/cmd/helloworld/目录下文件；
2）功能：客户端和服务器端简单实时通讯，创建双向数据流，包括head头数据等传递。

