# common service

### 目录介绍
```
├── common
├── component
├── essential
│   ├── api
│   └── config
│   └── logger
│   └── scheduler
│   └── storage
│   └── ticker
│   └── notification 事件通知
├── main
│   └── _your_app_
├── migration
├── module
│   ├── agent
│   └── agentgroup
```

- common 通用的一些函数或者常量，无状态的
- essential 通用的一些模块，是业务运行必须的模块，有状态的
- main 主程序目录
- migration 数据库migration代码
- module 业务模块，属于可有可无的，互相之间需要不能包含依赖关系，如果有依赖关系，需要经过essential调用

### 如何使用


