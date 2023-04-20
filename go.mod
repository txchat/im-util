module github.com/txchat/im-util

go 1.15

require (
	github.com/33cn/chain33 v1.67.3
	github.com/Terry-Mao/goim v0.0.0-20210523140626-e742c99ad76e
	github.com/ethereum/go-ethereum v1.10.16
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/gin-gonic/gin v1.8.2
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/oofpgDLD/kafka-go v1.1.0
	github.com/rs/zerolog v1.28.0
	github.com/spf13/cobra v1.5.0
	github.com/stretchr/testify v1.8.2
	github.com/txchat/dtalk v0.1.2
	github.com/txchat/im v0.0.1
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	github.com/zeromicro/go-zero v1.4.3
	golang.org/x/crypto v0.0.0-20221005025214-4161e89ecf1b
	google.golang.org/protobuf v1.28.1
)

replace (
	github.com/txchat/dtalk => ../dtalk
	github.com/txchat/im => ../im
)
