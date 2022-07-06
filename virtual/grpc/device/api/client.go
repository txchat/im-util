package device

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

type Client struct {
	conn   *grpc.ClientConn
	client DeviceClient
}

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

func New(addr string) *Client {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		panic(err)
	}
	return &Client{
		conn:   conn,
		client: NewDeviceClient(conn),
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Connect(ctx context.Context, in *ConnectReq, opts ...grpc.CallOption) (*ConnectReply, error) {
	return c.client.Connect(ctx, in, opts...)
}

func (c *Client) ReConnect(ctx context.Context, in *ReConnectReq, opts ...grpc.CallOption) (*ReConnectReply, error) {
	return c.client.ReConnect(ctx, in, opts...)
}

func (c *Client) DisConnect(ctx context.Context, in *DisConnectReq, opts ...grpc.CallOption) (*DisConnectReply, error) {
	return c.client.DisConnect(ctx, in, opts...)
}

func (c *Client) Focus(ctx context.Context, in *FocusReq, opts ...grpc.CallOption) (*FocusReply, error) {
	return c.client.Focus(ctx, in, opts...)
}

func (c *Client) Input(ctx context.Context, opts ...grpc.CallOption) (Device_InputClient, error) {
	return c.client.Input(ctx, opts...)
}

func (c *Client) Output(ctx context.Context, in *OutputReq, opts ...grpc.CallOption) (Device_OutputClient, error) {
	return c.client.Output(ctx, in, opts...)
}
