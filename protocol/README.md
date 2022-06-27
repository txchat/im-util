# Protocol工具
Protocol工具提供了丰富的功能：
1. wallet账户的批量生成和导入导出
2. 鉴权token的生成和校验
3. 根据指定算法进行签名和验签
4. seed的加密和解密
5. IM数据帧的校验
6. dh算法加解密

注意：以下操作如非特殊说明，默认在工程根目录下（工程根目录为仓库根目录的protocol文件夹下）

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
protocol -h
```
查看版本号：
```shell
protocol --version
#或
protocol version
```

其他命令参考：[命令详细描述文档](doc/markdown/protocol.md)