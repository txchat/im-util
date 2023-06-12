package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/txchat/dtalk/api/proto/chat"
	xerror "github.com/txchat/dtalk/pkg/error"
)

var (
	ErrSendMessageReplyInnerError             = errors.New(chat.SendMessageReply_InnerError.String())
	ErrSendMessageReplyUnSupportedMessageType = errors.New(chat.SendMessageReply_UnSupportedMessageType.String())
	ErrSendMessageReplyInsufficientPermission = errors.New(chat.SendMessageReply_InsufficientPermission.String())
	ErrSendMessageReplyIllegalFormat          = errors.New(chat.SendMessageReply_IllegalFormat.String())
	ErrSendMessageReplyOutdatedFormat         = errors.New(chat.SendMessageReply_OutdatedFormat.String())
)

type Response struct {
	Result  int64           `json:"result"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type SendResp struct {
	Code     int32  `json:"code"`
	Mid      string `json:"mid"`
	Datetime int64  `json:"datetime"`
	Repeat   bool   `json:"repeat"`
}

type AuthenticationMetadata struct {
	Signature  string
	UUID       string
	Device     string
	DeviceName string
	Version    string
}

type ChatHTTPAPIClient struct {
	serverAddr                string
	timeout                   time.Duration
	getAuthenticationMetadata func() *AuthenticationMetadata
}

func NewChatHTTPAPIClient(serverAddr string, timeout time.Duration, authMetadata func() *AuthenticationMetadata) *ChatHTTPAPIClient {
	return &ChatHTTPAPIClient{
		serverAddr:                serverAddr,
		timeout:                   timeout,
		getAuthenticationMetadata: authMetadata,
	}
}

func (c *ChatHTTPAPIClient) SendChatMessage(ctx context.Context, data []byte) (*SendResp, error) {
	u, err := url.Parse(c.serverAddr + "/app/record/send")
	if err != nil {
		return nil, err
	}

	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	if c.getAuthenticationMetadata != nil {
		authData := c.getAuthenticationMetadata()
		req.Header.Set("FZM-SIGNATURE", authData.Signature)
		req.Header.Set("FZM-UUID", authData.UUID)
		req.Header.Set("FZM-DEVICE", authData.DeviceName)
		req.Header.Set("FZM-DEVICE-NAME", authData.DeviceName)
		req.Header.Set("FZM-VERSION", authData.Version)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with response code: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var dtalkResp Response
	err = json.Unmarshal(respBody, &dtalkResp)
	if err != nil {
		return nil, err
	}
	if dtalkResp.Result != xerror.CodeOK {
		return nil, fmt.Errorf("request failed with custom error %d:%s", dtalkResp.Result, dtalkResp.Message)
	}

	var finalResp SendResp
	err = json.Unmarshal(dtalkResp.Data, &finalResp)
	if err != nil {
		return nil, err
	}
	return &finalResp, nil
}

func GetMid(finalResp *SendResp) (string, error) {
	failedType := chat.SendMessageReply_FailedType(finalResp.Code)
	if failedType != chat.SendMessageReply_IsOK {
		return "", fmt.Errorf("send message failed reson: %s", failedType.String())
	}
	return finalResp.Mid, nil
}
