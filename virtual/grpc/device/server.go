package device

import (
	"context"
	xgrpc "github.com/txchat/dtalk/pkg/net/grpc"
	"github.com/txchat/im-util/internel/device"
	protoutil "github.com/txchat/im-util/internel/proto"
	"github.com/txchat/im-util/pkg/net"
	pb "github.com/txchat/im-util/virtual/grpc/device/api"
	comet "github.com/txchat/im/api/comet/grpc"
	xproto "github.com/txchat/imparse/proto"
	"google.golang.org/grpc"
	"io"
	"time"
)

type Server struct {
	pb.UnimplementedDeviceServer

	*xgrpc.Server
	device  *device.Device
	rb      chan *xproto.Common
	session Session
}

func (s *Server) Connect(ctx context.Context, req *pb.ConnectReq) (*pb.ConnectReply, error) {
	var err error
	return &pb.ConnectReply{}, err
}

func (s *Server) ReConnect(ctx context.Context, req *pb.ReConnectReq) (*pb.ReConnectReply, error) {
	var err error
	return &pb.ReConnectReply{}, err
}

func (s *Server) DisConnect(context.Context, *pb.DisConnectReq) (*pb.DisConnectReply, error) {
	return &pb.DisConnectReply{}, nil
}

func (s *Server) Focus(ctx context.Context, req *pb.FocusReq) (*pb.FocusReply, error) {
	s.session.SetTarget(req.GetTarget())
	s.session.SetChannel(xproto.Channel(req.GetChannelType()))
	return &pb.FocusReply{}, nil
}

func (s *Server) Input(req pb.Device_InputServer) error {
	for {
		data, err := req.Recv()
		if err == io.EOF {
			return req.SendAndClose(&pb.InputReply{})
		}
		if err != nil {
			return err
		}
		if s.session.GetTarget() == "" {
			return req.SendAndClose(&pb.InputReply{
				Err: "no focus target",
			})
		}
		s.device.SendMsg(xproto.Channel_name[int32(s.session.GetChannel())], s.session.GetTarget(), data.GetText())
	}
}

func (s *Server) Output(req *pb.OutputReq, reply pb.Device_OutputServer) error {
	for {
		msg := <-s.rb
		err := reply.Send(&pb.OutputReply{
			Msg:         string(msg.Msg),
			ChannelType: int32(msg.GetChannelType()),
			Target:      msg.GetTarget(),
		})
		if err != nil {
			return err
		}
	}
}

func (s *Server) OnReceive(c *net.IMConn, proto *comet.Proto) error {
	bizProto, err := protoutil.ConvertBizProto(proto.GetBody())
	if err != nil {
		return err
	}
	if bizProto.GetEventType() == xproto.Proto_common {
		common, err := protoutil.ConvertCommon(bizProto.GetBody())
		if err != nil {
			return err
		}
		s.rb <- common
	}
	return nil
}

func NewServer(c *xgrpc.ServerConfig, d *device.Device) *Server {
	//serve gRPC
	srv := &Server{
		UnimplementedDeviceServer: pb.UnimplementedDeviceServer{},
		device:                    d,
	}
	connectionTimeout := grpc.ConnectionTimeout(time.Second * 7)
	ws := xgrpc.NewServer(
		c,
		connectionTimeout,
	)
	pb.RegisterDeviceServer(ws.Server(), srv)
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	srv.Server = ws
	return srv
}
