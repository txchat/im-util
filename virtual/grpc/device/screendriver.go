package device

import (
	"context"
	"github.com/txchat/im-util/pkg/net"
	"github.com/txchat/im-util/pressure/pkg/msggenerator"
	comet "github.com/txchat/im/api/comet/grpc"
	xproto "github.com/txchat/imparse/proto"
	"io"
	"time"
	xgrpc "github.com/txchat/dtalk/pkg/net/grpc"
	pb "github.com/txchat/im-util/pressure/pkg/device/api"
	"google.golang.org/grpc"
	"github.com/txchat/im-util/internel/device"
)

type Server struct {
	pb.UnimplementedDriverServer

	*xgrpc.Server
	device *device.Device
	rb chan []byte
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

func (s *Server) Input(req pb.Driver_InputServer) error {
	for {
		data, err := req.Recv()
		if err == io.EOF {
			return req.SendAndClose(&pb.InputReply{})
		}
		if err != nil {
			return err
		}
		s.device.SendMsg(, , data.GetText())
	}
}

func (s *Server) Output(req *pb.OutputReq, reply pb.Driver_OutputServer) error {
	for {
		msg := <- s.rb
		err := reply.Send(&pb.OutputReply{
			Msg: string(msg),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) OnReceive(c *net.IMConn, proto *comet.Proto) error{
	bizProto, err := msggenerator.ConvertBizProto(proto.GetBody())
	if err != nil {
		return err
	}
	if bizProto.GetEventType() == xproto.Proto_common {
		common, err := msggenerator.ConvertCommon(bizProto.GetBody())
		if err != nil {
			return err
		}
		s.rb <- common.Msg
	}
	return nil
}

func NewServer(c *xgrpc.ServerConfig, d *device.Device) *Server {
	//serve gRPC
	srv := &Server{
		UnimplementedDriverServer: pb.UnimplementedDriverServer{},
		device:                    d,
	}
	connectionTimeout := grpc.ConnectionTimeout(time.Second * 7)
	ws := xgrpc.NewServer(
		c,
		connectionTimeout,
	)
	pb.RegisterDriverServer(ws.Server(), srv)
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	srv.Server = ws
	return srv
}
