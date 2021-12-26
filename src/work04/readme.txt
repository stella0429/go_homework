【week04作业】

1. 按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包含对数据层、业务层、API 注册，
以及 main 函数对于服务的注册和启动，信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。

解答：基于goin框架和go-kratos框架构建的基本目录结构和工程说明如下（代码见同目录下项目）：
.
└── gogin
    ├── api                        --api目录
    │   ├── myapp01                --应用服务1
    │   │   └── v1                 --包含自测请求api，路由，以及用protobuf生成文件实践
    │   │       ├── order          --订单模块，路由及api test
    │   │       └── user           --用户模块，路由及api test
    │   └── myapp02                --应用服务2，基于go-kratos框架（有些bug待处理）
    │       └── v1                 --api定义
    ├── cmd                        --服务主干目录，main.go
    │   ├── myapp01                --应用服务1，main.go以及用wire构建的依赖实践
    │   └── myapp02                --应用服务2，main.go以及用wire构建的依赖实践
    ├── configs                    --json配置文件
    ├── docs                       --说明文件目录
    ├── internal                   --私用应用程序和库代码
    │   ├── app
    │   │   ├── myapp01            
    │   │   │   ├── biz
    │   │   │   ├── conf
    │   │   │   ├── data
    │   │   │   ├── server
    │   │   │   └── service
    │   │   └── myapp02             --基于go-kratos框架（有些bug待处理），注册数据层、业务层、API 注册
    │   │       ├── biz
    │   │       ├── conf
    │   │       ├── data
    │   │       ├── server
    │   │       └── service
    │   └── pkg                     --内部公用包
    ├── pkg                         --内外部公用包
    │   ├── cache
    │   │   ├── memcache
    │   │   └── redis
    │   └── conf
    ├── routers                     --路由
    ├── third_party                 --第三方库包
    └── utils                       --搬运的go-kratos框架库，有点问题
        ├── api
        │   └── metadata
        ├── encoding
        │   ├── form
        │   ├── json
        │   ├── proto
        │   ├── xml
        │   └── yaml
        ├── errors
        ├── internal
        │   ├── context
        │   ├── endpoint
        │   ├── host
        │   ├── httputil
        │   └── testdata
        │       ├── binding
        │       ├── complex
        │       ├── encoding
        │       └── helloworld
        ├── log
        ├── middleware
        │   ├── auth
        │   │   └── jwt
        │   ├── circuitbreaker
        │   ├── logging
        │   ├── metadata
        │   ├── metrics
        │   ├── ratelimit
        │   ├── recovery
        │   ├── selector
        │   ├── tracing
        │   └── validate
        ├── registry
        ├── selector
        │   ├── filter
        │   ├── node
        │   │   ├── direct
        │   │   └── ewma
        │   ├── p2c
        │   ├── random
        │   └── wrr
        └── transport
            ├── grpc
            │   └── resolver
            │       ├── direct
            │       └── discovery
            └── http
                ├── binding
                ├── pprof
                └── status

86 directories