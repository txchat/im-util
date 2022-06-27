## protocol frame auth check

校验鉴权帧

### Synopsis

通过标签-d将鉴权帧的完整数据内容传递给命令工具来校验鉴权是否通过，可以选择是否开启过期时间。

```
protocol frame auth check [flags]
```

### Options

```
  -d, --data string   the frame data encoded by base64
  -h, --help          help for check
      --timeout       check timeout enable -t=[false]
  -T, --type string   check app type -T=[dtalk] (default "dtalk")
```

### SEE ALSO

* [protocol frame auth](protocol_frame_auth.md)	 - 鉴权数据帧相关命令

###### Auto generated by spf13/cobra on 27-Jun-2022