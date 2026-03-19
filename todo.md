第二步：引入连接分桶（Bucket）
目标：优化并发性能，减少锁竞争

需要实现：

将连接分散到多个 Bucket
每个 Bucket 独立加锁
根据用户 ID hash 到不同 Bucket
新增文件：


├── bucket.go         # 连接分桶
验收标准：

支持 1000+ 并发连接
锁竞争明显降低
第三步：实现房间（Room）功能
目标：支持群聊和广播

需要实现：

Room 结构体
用户加入/退出房间
房间消息广播
新增文件：


├── room.go           # 房间管理
验收标准：

支持创建房间
房间内消息能广播给所有成员
第四步：拆分为微服务架构
目标：将单体应用拆分为 Comet + Logic

需要实现：
0
Comet 服务：只负责连接管理
Logic 服务：负责业务逻辑和推送接口
两者通过 gRPC 通信
文件结构：


goim-v2/
├── cmd/
│   ├── comet/
│   │   └── main.go
│   └── logic/
│       └── main.go
├── internal/
│   ├── comet/
│   │   ├── server.go
│   │   ├── bucket.go
│   │   └── grpc.go
│   └── logic/
│       ├── http.go
│       └── grpc.go
└── api/
    └── comet.proto    # gRPC 接口定义
验收标准：

Comet 和 Logic 独立部署
Logic 通过 gRPC 调用 Comet 推送消息
第五步：引入消息队列（Kafka）
目标：异步推送，提高吞吐量

需要实现：

Job 服务：消费 Kafka 消息
Logic 将推送请求写入 Kafka
Job 从 Kafka 读取并调用 Comet
新增：


├── cmd/
│   └── job/
│       └── main.go
└── internal/
    └── job/
        ├── consumer.go
        └── push.go
验收标准：

推送请求异步处理
支持高并发推送
第六步：服务发现和负载均衡
目标：支持多 Comet 节点

需要实现：

Comet 注册到服务发现（etcd/consul）
Logic 通过服务发现找到 Comet
负载均衡策略
验收标准：

支持多个 Comet 实例
Logic 能正确路由到用户所在的 Comet
第七步：性能优化
目标：支持百万级连接

优化点：

内存池（sync.Pool）
零拷贝（减少内存分配）
协议优化（二进制协议）
读写分离