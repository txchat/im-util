## pressure conn keep

启动指定数量的客户端连接

### Synopsis

启动指定数量的客户端连接，并保持心跳，不发送消息

```
pressure conn keep [flags]
```

### Options

```
  -a, --appId string     (default "dtalk")
  -h, --help            help for keep
  -i, --in string       users store file path (default "./users.txt")
      --rs string       存储用户信息的字段分隔符[默认：,] (default ",")
  -s, --server string   server address (default "172.16.101.107:3102")
  -t, --time string      (default "720h")
  -u, --users int       users number (default 2)
```

### SEE ALSO

* [pressure conn](pressure_conn.md)	 - 批量连接相关命令

###### Auto generated by spf13/cobra on 11-Jul-2022