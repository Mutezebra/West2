#### 用户数量 100
#### 总消息条数 100000

#### mysql-start 660705
#### mysql记录数 760705        
**test结束之后又等待了一段时间，mysql的数据才全部保存下来**

#### 总耗时 59.83s
#### 总循环时间 20.6033s
#### cpu使用时间 4.43s
#### 时间比 4.43:1

```
当前只是简单的通过一个go func异步处理消息，
可以看到时间比再次降低，但还是有阻塞，所以打算
下一步打算多开Manager。
```