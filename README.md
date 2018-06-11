go-tlog
=======

TLOG工具库

# 功能组件:

### 生成可执行文件tlog
1. 如果修改了[cmd/main.go](./cmd/main.go)，要`make`生成新的可执行文件
2. 生成的可执行文件在`dist`文件夹

### 使用可执行文件tlog（xml转换为go代码）
1. 在[Releases](https://github.com/goofansu/go-tlog/releases)中下载对应平台的可执行文件
2. `tlog`默认根据腾讯方提供的`tlog.xml`生成`tlogevt.go`，更多设置参见`tlog -h`
3. `tlogevt.go`中，指针表示必填项，其他表示可选项

### tlogrus（生成TLOG格式的日志）
1. 提供`tlogrus.Log(event interface{})`接口，任何struct都可记录成TLOG格式数据。
2. 安装: `go get github.com/goofansu/go-tlog/tlogrus`
3. 使用方法：参见[examples/basic.go](examples/basic/basic.go)

生成的日志格式如下：
```bash
# StructName|v1|v2|...
PlayerRegister|1|2017-02-20 19:48:29|100695782|0|1|0FEB999268EF9FEA26D4CB219C37910D|NULL|NULL|NULL|Unicom|NULL|0|0|0|1|NULL|0|NULL|NULL|NULL
PlayerExpFlow|2|2017-02-20 19:48:29|100695783|1|2|0FEB999268EF9FEA26D4CB219C37910C|0|0|0|0|0|0
```

### UDP hook
1. 提供`hooks/udp`，可通过UDP形式发送日志
2. 安装：`go get github.com/goofansu/go-kit/hooks/udp`
3. 使用方法：参见[hooks/udp/udp.go](hooks/udp/udp.go)

### FPM(File Per Message) hook
1. 提供`hooks/fpm`，可以把日志按事件记录到各自的log中，比如`PlayerLogin.log`, `MoneyFlow.log`
2. 安装：`go get github.com/goofansu/go-kit/hooks/fpm`
3. 使用方法：参见[hooks/fpm/fpm.go](hooks/fpm/fpm.go)

# 新增一条日志的步骤
1. 从腾讯经分处获取需求，制定日志格式
2. 生成tlogevt.go
3. 编写业务逻辑
4. 检查日志是否生成
