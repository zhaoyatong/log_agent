# log_agent本地日志收集

功能：
+ 自动收集本地日志文件内容并将日志内容发送至kafka。
+ 依赖etcd做配置管理，支持配置热更新而无需重启。


agent配置文件为conf.ini，具体配置项可见注释。

etcd配置为json格式，示例如下：

``[{
	"path": "./my.log",
	"topic": "web_log"
}, {
	"path": "./aaa.log",
	"topic": "aab"
}]``


+ path：本地日志文件路径。
+ topic：对应kafka的topic。
