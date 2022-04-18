package device

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	xgrpc "github.com/txchat/dtalk/pkg/net/grpc"
	pb "github.com/txchat/im-util/client-pressurev2/pkg/device/api"
	"github.com/txchat/im-util/client-pressurev2/pkg/user"
	xproto "github.com/txchat/imparse/proto"
	"google.golang.org/grpc"
)

type ScreenDriver interface {
}

type UserDeviceOpt struct {
	AppId, Server string

	Address        string
	PriKey, PubKey []byte

	Uuid, DeviceName string
	DeviceType       xproto.Device
}

type NetScreenDriver struct {
	nss    *NetScreenDriverServer
	device *Device

	appId, server string

	address        string
	priKey, pubKey []byte

	uuid, deviceName string
	deviceType       xproto.Device
	log              zerolog.Logger
}

func NewNetScreenDriver(opt *UserDeviceOpt, c *xgrpc.ServerConfig) *NetScreenDriver {
	//
	ns := &NetScreenDriver{
		nss:        nil,
		appId:      opt.AppId,
		server:     opt.Server,
		address:    opt.Address,
		priKey:     opt.PriKey,
		pubKey:     opt.PubKey,
		uuid:       opt.Uuid,
		deviceName: opt.DeviceName,
		deviceType: opt.DeviceType,
		log:        zerolog.New(os.Stderr).With().Timestamp().Logger(),
	}
	nss := NewNetScreenDriverServer(c, ns)
	ns.nss = nss
	return ns
}

func (ns *NetScreenDriver) StartUp() {
	log.Info().Msg("Start Up device")
	u := user.NewUser(ns.address, ns.priKey, ns.pubKey)
	ns.device = NewDevice(ns.uuid, ns.deviceName, ns.deviceType, ns.log, u)
}

func (ns *NetScreenDriver) Shutdown(ctx context.Context) {
	ns.nss.Shutdown(ctx)
	ns.device.Destroy()
}

func (ns *NetScreenDriver) Connect() error {
	return ns.device.ConnectIMServer(ns.appId, ns.server)
}

func (ns *NetScreenDriver) ReConnect() error {
	return ns.device.ReConnectIMServer(ns.appId, ns.server)
}

func (ns *NetScreenDriver) DisConnect() {
	ns.device.DisConnectIMServer()
}

func (ns *NetScreenDriver) ChangeCurrentPage(chType int32, target string) {
	ns.device.GetSession().ChangePage(xproto.Channel(chType), target)
}

func (ns *NetScreenDriver) Input(channelType, target, text string) {
	log.Info().Str("channelType", channelType).
		Str("target", target).
		Str("text", text).Msg("receive input")
	ns.device.SendMsg(channelType, target, text)
}

func (ns *NetScreenDriver) OutPut() (string, error) {
	if ns.device.GetIsDestroyed() {
		return "", io.EOF
	}
	select {
	case common := <-ns.device.GetOutputQueue():
		if common == nil {
			//
			return "", fmt.Errorf("common is nil")
		}
		chType := ""
		msg := ""
		switch common.GetChannelType() {
		case xproto.Channel_ToUser:
			chType = "私聊消息"
		case xproto.Channel_ToGroup:
			chType = "群聊消息"
		}

		if ns.device.GetUserId() != common.GetTarget() {
			return "", nil
		}

		if ns.device.GetUserId() == common.GetFrom() {
			msg = fmt.Sprintf("[%s]---[%s:%s]", common.GetMsg(), chType, common.GetFrom())
		} else {
			msg = fmt.Sprintf("[%s:%s]---[%s]", chType, common.GetFrom(), common.GetMsg())
		}
		return msg, nil
	case <-ns.device.GetDestroy():
		return "", io.EOF
	}
}

type NetScreenDriverServer struct {
	*xgrpc.Server
	pb.UnimplementedDriverServer
	ns *NetScreenDriver
}

func (nss *NetScreenDriverServer) Connect(ctx context.Context, req *pb.ConnectReq) (*pb.ConnectReply, error) {
	err := nss.ns.Connect()
	return &pb.ConnectReply{}, err
}

func (nss *NetScreenDriverServer) ReConnect(ctx context.Context, req *pb.ReConnectReq) (*pb.ReConnectReply, error) {
	err := nss.ns.ReConnect()
	return &pb.ReConnectReply{}, err
}

func (nss *NetScreenDriverServer) DisConnect(context.Context, *pb.DisConnectReq) (*pb.DisConnectReply, error) {
	nss.ns.DisConnect()
	return &pb.DisConnectReply{}, nil
}

func (nss *NetScreenDriverServer) ChangeCurrentPage(ctx context.Context, req *pb.ChangeCurrentPageReq) (*pb.ChangeCurrentPageReply, error) {
	nss.ns.ChangeCurrentPage(req.GetChannelType(), req.GetTarget())
	return &pb.ChangeCurrentPageReply{}, nil
}

func (nss *NetScreenDriverServer) Input(req pb.Driver_InputServer) error {
	for {
		data, err := req.Recv()
		if err == io.EOF {
			return req.SendAndClose(&pb.InputReply{})
		}
		if err != nil {
			return err
		}
		s := nss.ns.device.GetSession()
		chType, target := s.GetPage()
		nss.ns.Input(chType.String(), target, data.GetText())
	}
}

func (nss *NetScreenDriverServer) Output(req *pb.OutputReq, reply pb.Driver_OutputServer) error {
	for {
		msg, err := nss.ns.OutPut()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		err = reply.Send(&pb.OutputReply{
			Msg: msg,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func NewNetScreenDriverServer(c *xgrpc.ServerConfig, ns *NetScreenDriver) *NetScreenDriverServer {
	//serve gRPC
	nss := &NetScreenDriverServer{
		UnimplementedDriverServer: pb.UnimplementedDriverServer{},
		ns:                        ns,
	}
	connectionTimeout := grpc.ConnectionTimeout(time.Second * 7)
	ws := xgrpc.NewServer(
		c,
		connectionTimeout,
	)
	pb.RegisterDriverServer(ws.Server(), nss)
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	nss.Server = ws
	return nss
}
