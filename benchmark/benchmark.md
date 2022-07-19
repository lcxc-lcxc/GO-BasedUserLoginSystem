

# 压测服务器
- 机器个数：1
- CPU：8核
- 内存：16G
- 系统：MaxOS
- 数据库MySQL：8.0.29
- 数据库数据量：10,000,000

# 观察指标
主要关注下面指标
- QPS
- 失败率
- 平均耗时
- 最长耗时
- TP90 TP99 TP999


# 压测工具简介
使用开源压测工具go-stress-testing，工具具体介绍如下
[go-stress-testing](https://github.com/link1st/go-stress-testing)。


- 第一步下载 [go-stress-testing](https://github.com/link1st/go-stress-testing/releases)。
记得要直接clone。不要使用release里面的  

- 第二步创建压测数据。
```
curl --location --request POST 'http://localhost:8080/api/user/login' \
-H 'Connection: keep-alive' \
-H 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'username=mytest5&password=1234567' \

```
- 第三步将上面请求数据写入curl/login.txt文件中。

- 第四步选择压测参数。
```
  -H value
    	自定义头信息传递给服务器 示例:-H 'Content-Type: application/json'
  -c uint
    	并发数 (default 1)
  -d string
    	调试模式 (default "false")
  -data string
    	HTTP POST方式传送数据
  -n uint
    	请求数(单个并发/协程) (default 1)
  -p string
    	curl文件路径
  -m int 
      连接数
  -k 
     开启长连接         
  -u string
    	压测地址
  -v string
    	验证方法 http 支持:statusCode、json webSocket支持:json
```
go run main.go -c 200 -n 500 -m 300 -k -p curl/get.txt

# 压测结果
## 压测结论
1. 所有压测请求均返回正确的结果
2. 200并发固定用户(get接口)QPS达到2501，未达到QPS大于3000的目标。
3. 200并发随机用户(login接口)QPS达到2000，达到QPS大于1000的目标.
4. 2000并发固定用户(get接口)QPS达到2105，达到QPS大于1500的目标。
5. 2000并发随机用户(login接口)QPS达到1559，达到QPS大于800的目标。

注：
压测过程中的tips:
1. 客户端发起的请求使用HTTP 1.1协议，使用长连接，不要使用短连接。
2. 客户端请求连接数可适当调大，避免出现`read: connection reset by peer`
3. 压测过程中关掉不必要的后台程序，避免产生 `Socket/File : too many open files(打开的文件过多) `、`EOF`报错。
4. 在Mac进行性能测试可能存在系统资源限制，可尝试下面指令：
```
sudo sysctl -w kern.ipc.somaxconn=2048
sudo sysctl -w kern.maxfiles=12288
ulimit -n 10000
```