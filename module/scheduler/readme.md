# 调度器设计

# 概念梳理
- 调度器是一个可以指定namespace，不同的namespace，是不同的调度器，调度器，定时轮询查询状态，默认有一个default的调度器
- JOB，任务
- Provider，添加、删除、查询Job、结果反馈，与生产端关联比较大
- Consumer，负责任务的消费，以及结果反馈，与消费端关联比较大

# job状态
- Init，初始化创建
- Pending，正在加入，未运行
- Running，正在运行
- Killed，停止的任务
- Retryable，重试
- Failed
- Succeeded

# job的调度对象的类型，就是谁来执行这个job
- group 表示由哪个组来执行
- anyone 表示由这个组的任意一个来执行
- everyone 表示由这个组的每一个来执行
- :id   表示由这个id来执行

# 数据结构
zset ,用于存储优先级以及job的key 
hash，用于存储job的内容，为数据库job的缓存，有超时时间
sqldb中job，三者由id关联

# job的probe
- type job类型
- PendingCallback 
- ProgressCallback
- DoneCallback
- RunningCallback

# api
- 新建scheduler
- 获取scheduler
- 新建provider
- 新建consumer
- 创建job request
- 添加job
- 获取job
- 注册job回调
