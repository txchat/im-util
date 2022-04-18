package device

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

type Client struct {
	conn   *grpc.ClientConn
	client DriverClient
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
		client: NewDriverClient(conn),
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

func (c *Client) ChangeCurrentPage(ctx context.Context, in *ChangeCurrentPageReq, opts ...grpc.CallOption) (*ChangeCurrentPageReply, error) {
	return c.client.ChangeCurrentPage(ctx, in, opts...)
}

func (c *Client) Input(ctx context.Context, opts ...grpc.CallOption) (Driver_InputClient, error) {
	return c.client.Input(ctx, opts...)
}

func (c *Client) Output(ctx context.Context, in *OutputReq, opts ...grpc.CallOption) (Driver_OutputClient, error) {
	return c.client.Output(ctx, in, opts...)
}
