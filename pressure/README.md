# Pressure工具
Pressure是提供批量启动客户端连接相关的工具：
1. 压力测试工具
2. 压力测试结果分析工具
3. 批量客户端连接工具

注意：以下操作如非特殊说明，默认在工程根目录下（工程根目录为仓库根目录的pressure文件夹下）

---
## 编译
### 依赖
1. git
2. golang1.17
3. gcc

在工程根目录下执行：
```shell
make build
```
交叉编译等其他详细编译选项参考:
```shell
make help
```

## 使用
查看帮助：
```shell
pressure -h
```
查看版本号：
```shell
pressure --version
#或
pressure version
```

其他命令参考：[命令详细描述文档](doc/markdown/pressure.md)