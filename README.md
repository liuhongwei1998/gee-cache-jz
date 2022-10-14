# gee-cache-jz
gee-cache-jz是仿照groupcache实现的分布式缓存

groupcache为Go语言版本的memcached
为学习实现的玩具

支持特性有：

1、单机缓存和基于 HTTP 的分布式缓存

2、最近最少访问(Least Recently Used, LRU) 缓存策略

3、使用 Go 锁机制防止缓存击穿

4、使用一致性哈希选择节点，实现负载均衡

5、使用 protobuf 优化节点间二进制通信

