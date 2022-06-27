# im-util
仓库提供了如协议校验工具，压测工具、模拟客户端等工具。
1. [Protocol](#Protocol工具)

## Protocol工具
> Protocol工具，是为了便于在IM客户端开发过程中校验各种格式是否正确的工具。

在如下场景你可以会需要使用[Protocol工具](protocol/README.md)：
1. 作为客户端与IM服务端进行通讯协议的握手时，使用工具校验鉴权帧。
2. 调用接口时鉴权token的校验。
3. 消息协议数据帧的校验。
4. 需要批量生成用户时
5. 采用DH加解密数据，对算法正确性验证。
6. 使用密码加密和解密seed时。

### License

im-util is under the MIT license. See the [LICENSE](LICENSE) file for details.
